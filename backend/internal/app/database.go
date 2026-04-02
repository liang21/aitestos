// Package app provides database connection management
package app

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/liang21/aitestos/internal/config"
)

// NewDB creates a new database connection pool
func NewDB(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate database config: %w", err)
	}

	db, err := sqlx.Connect("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

// DBCloser wraps sqlx.DB to implement Closer interface
type DBCloser struct {
	*sqlx.DB
	name string
}

// NewDBCloser creates a closer wrapper for database
func NewDBCloser(db *sqlx.DB, name string) *DBCloser {
	return &DBCloser{DB: db, name: name}
}

// Name returns the database name
func (d *DBCloser) Name() string {
	return d.name
}

// Close closes the database connection pool
func (d *DBCloser) Close(ctx context.Context) error {
	return d.DB.Close()
}
