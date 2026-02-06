package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"assetManager/internal/config"
)

// DB wraps sqlx.DB for the application
type DB struct {
	*sqlx.DB
}

// New creates a new database connection
func New(cfg *config.DatabaseConfig) (*DB, error) {
	db, err := sqlx.Connect("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
