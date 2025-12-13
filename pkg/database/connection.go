package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/TinyMurky/tinyurl/pkg/logging"
)

// DB is the pool of database Conn
type DB struct {
	Pool *sql.DB
}

// NewFromEnv sets up the database connections using the configuration in the
// process's environment variables. This should be called just once per server
// instance.
func NewFromEnv(ctx context.Context, cfg *Config) (*DB, error) {
	logger := logging.FromContext(ctx)

	dsn := cfg.ToDSN()

	pool, err := sql.Open("sqlite", dsn)

	if err != nil {
		return nil, fmt.Errorf("fail to open sqlite connection pool: %w", err)
	}

	newDB := DB{
		Pool: pool,
	}

	logger.Infof("Open connection pool with dsn: %s", dsn)
	return &newDB, nil
}

// Close releases database connections
func (db *DB) Close(ctx context.Context) {
	logger := logging.FromContext(ctx)
	logger.Info("Closing connection pool...")
	if err := db.Pool.Close(); err != nil {
		logger.Errorf("Closing connection pool error: %s", err.Error())
	}

	logger.Info("Closing connection pool successfully")
}
