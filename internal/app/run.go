// Package app contains the main methods for running the server.
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	_ "github.com/golang/mock/mockgen/model"
	"github.com/pavlegich/events-store/internal/controllers/handlers"
	"github.com/pavlegich/events-store/internal/infra/config"
	"github.com/pavlegich/events-store/internal/infra/database"
	"github.com/pavlegich/events-store/internal/infra/logger"
	"github.com/pavlegich/events-store/internal/repository"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

// Run initializes the main app components and runs the server.
func Run() error {
	// Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Logger
	err := logger.Init(ctx, "Info")
	if err != nil {
		return fmt.Errorf("Run: logger initialization failed %w", err)
	}
	defer logger.Log.Sync()

	// Configuration
	cfg := config.NewConfig(ctx)
	err = cfg.ParseFlags(ctx)
	if err != nil {
		return fmt.Errorf("Run: parse flags failed %w", err)
	}

	// Database
	db, err := database.Init(ctx, cfg.DSN)
	if err != nil {
		return fmt.Errorf("Run: initialize database failed %w", err)
	}
	defer db.Close()

	// Router
	ctrl := handlers.NewController(ctx, cfg)
	repo := repository.NewEventRepository(ctx, db)
	router, err := ctrl.BuildRoute(ctx, repo)
	if err != nil {
		return fmt.Errorf("Run: build server route failed %w", err)
	}

	// Server
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	// Server graceful shutdown
	go func() {
		<-ctx.Done()
		if ctx.Err() != nil {
			ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 5*time.Second)
			defer cancelShutdown()

			logger.Log.Info("shutting down gracefully...",
				zap.Error(ctx.Err()))

			err := srv.Shutdown(ctxShutdown)
			if err != nil {
				logger.Log.Error("server shutdown failed",
					zap.Error(err))
			}
		}
	}()

	logger.Log.Info("running server", zap.String("addr", srv.Addr))

	return srv.ListenAndServe()
}
