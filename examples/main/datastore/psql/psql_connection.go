package psql

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// pool represents a single instance (singleton) of a connection pool.
var pool *pgxpool.Pool

// Pool returns a single instance of the connection pool at any point.
func Pool() (*pgxpool.Pool, error) {
	if pool == nil {
		config, err := pgxpool.ParseConfig(os.Getenv("DB_CONNECTION_STRING"))
		if err != nil {
			return nil, fmt.Errorf("could not create database connection config: %w", err)
		}

		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			return nil, fmt.Errorf("could not create database connection pool: %w", err)
		}
	}

	return pool, nil
}
