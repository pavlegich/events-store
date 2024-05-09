// Package database contains methods for database initialization.
package database

import (
	"context"
	"database/sql"
	"fmt"
)

// Init initializes database and creates the tables
// from the specified migrations.
func Init(ctx context.Context, driver string, path string) (*sql.DB, error) {
	db, err := sql.Open(driver, path)
	if err != nil {
		return nil, fmt.Errorf("Init: couldn't open database %w", err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("Init: connection with database is died %w", err)
	}

	return db, nil
}
