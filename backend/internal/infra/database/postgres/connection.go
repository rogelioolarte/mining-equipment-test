package postgres

import (
	"context"
	"fmt"

	"minning/internal/infra/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	Pool *pgxpool.Pool
}

func NewConnection(cfg *config.DB) (*Connection, error) {
	pool, err := pgxpool.New(context.Background(), cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	pool.Config().MaxConns = int32(cfg.MaxConns)
	pool.Config().MinConns = int32(cfg.MaxConns)
	pool.Config().MaxConnIdleTime = 5 * time.Minute

	return &Connection{Pool: pool}, nil

}

func NewTaskConnection(cfg *config.DB) (*Connection, error) {
	pool, err := pgxpool.New(context.Background(), cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	pool.Config().MaxConns = int32(10)
	pool.Config().MinConns = int32(2)
	pool.Config().MaxConnIdleTime = 5 * time.Minute

	return &Connection{Pool: pool}, nil

}

func (c *Connection) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
}

func (c *Connection) GetPool() *pgxpool.Pool {
	return c.Pool
}
