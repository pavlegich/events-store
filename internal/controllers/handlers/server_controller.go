// Package handlers contains server controller object and
// methods for building the server route, command functions
// for activating the command handler in controller
// and command handlers.
package handlers

import (
	"context"
	"net/http"

	"github.com/pavlegich/events-store/internal/controllers/middlewares"
	"github.com/pavlegich/events-store/internal/infra/config"
)

// Controller contains database and configuration
// for building the server router.
type Controller struct {
	cfg *config.Config
}

// NewController creates and returns new server controller.
func NewController(ctx context.Context, cfg *config.Config) *Controller {
	return &Controller{
		cfg: cfg,
	}
}

// BuildRoute creates new router and appends handlers and middlewares to it.
func (c *Controller) BuildRoute(ctx context.Context) (http.Handler, error) {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello!"))
	})

	handler := middlewares.Recovery(router)
	handler = middlewares.WithLogging(handler)

	return handler, nil
}
