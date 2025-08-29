package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{pool: pool}
}

func GetConnection() (*pgxpool.Pool, error) {
	ctx := context.Background()
	dsn := "postgres://user:pass@localhost:5430/db?sslmode=disable" // TODO Добавить cfg
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}
	
	return pool, nil
}

func (db *DB) Close() {
	db.pool.Close()
}
