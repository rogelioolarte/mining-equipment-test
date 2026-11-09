package database

import (
	"context"
	"fmt"
	"minning/internal/domain/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionManager interface {
	DB(ctx context.Context) ports.Executor
	Begin(ctx context.Context) (context.Context, error)
	BeginTx(ctx context.Context, opts pgx.TxOptions) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type contextKey string

const txKey contextKey = "pg_tx_key"

type pgTxManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) (TransactionManager, error) {
	if pool == nil {
		return nil, fmt.Errorf("pool required")
	}
	return &pgTxManager{pool: pool}, nil
}

func (m *pgTxManager) DB(ctx context.Context) ports.Executor {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}
	return m.pool
}

func (m *pgTxManager) Begin(ctx context.Context) (context.Context, error) {
	return m.BeginTx(ctx, pgx.TxOptions{})
}

func (m *pgTxManager) BeginTx(ctx context.Context, opts pgx.TxOptions) (context.Context, error) {
	if _, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return ctx, fmt.Errorf("the transaction already exists in the context.")
	}

	tx, err := m.pool.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to start the transaction: %w", err)
	}

	return context.WithValue(ctx, txKey, tx), nil
}

func (m *pgTxManager) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("there is no active transaction in the context to commit.")
	}
	return tx.Commit(ctx)
}

func (m *pgTxManager) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("there is no active transaction in the context to roll back")
	}
	return tx.Rollback(ctx)
}

func (m *pgTxManager) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	txCtx, err := m.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting the transaction: %w", err)
	}

	committed := false
	defer func() {
		if p := recover(); p != nil {
			_ = m.Rollback(txCtx)
			panic(p)
		} else if !committed {
			_ = m.Rollback(txCtx)
		}
	}()

	if err := fn(txCtx); err != nil {
		return err
	}

	if err := m.Commit(txCtx); err != nil {
		return fmt.Errorf("error confirming the transaction: %w", err)
	}

	committed = true
	return nil
}
