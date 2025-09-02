package postgres

import (
	"context"
	"fmt"

	"github.com/ent1k1377/subscriptions/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}

func GetConnection(cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	ctx := context.Background()
	dsn := cfg.DSN()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return pool, nil
}

func (db *DB) Close() {
	db.pool.Close()
}
