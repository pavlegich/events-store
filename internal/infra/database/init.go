// Package database contains methods for database initialization.
package database

import (
	"context"
	"database/sql"
	"fmt"
)

// go:embed migrations/*.sql
// var embedMigrations embed.FS

// Init initializes database and creates the tables
// from the specified migrations.
func Init(ctx context.Context, path string) (*sql.DB, error) {
	db, err := sql.Open("clickhouse", path)
	if err != nil {
		return nil, fmt.Errorf("Init: couldn't open database %w", err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Init: connection with database is died %w", err)
	}

	// Migrations
	// goose.SetBaseFS(embedMigrations)
	// err = goose.SetDialect("clickhouse")
	// if err != nil {
	// 	return nil, fmt.Errorf("Init: goose set dialect failed %w", err)
	// }
	// err = goose.Up(db, "migrations")
	// if err != nil {
	// 	return nil, fmt.Errorf("Init: goose up failed %w", err)
	// }

	return db, nil
}
