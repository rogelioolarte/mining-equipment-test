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

	// Transaction encapsula la lógica completa: Begin, defer recover/rollback, Commit y Rollback por errores.
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type contextKey string

const txKey contextKey = "pg_tx_key"

type pgTxManager struct {
	pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) (TransactionManager, error) {
	if pool == nil {
		return nil, fmt.Errorf("Pool Required")
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
		return ctx, fmt.Errorf("la transacción ya existe en el contexto")
	}

	tx, err := m.pool.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("fallo al iniciar la transacción: %w", err)
	}

	return context.WithValue(ctx, txKey, tx), nil
}

func (m *pgTxManager) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("no hay transacción activa en el contexto para confirmar")
	}
	return tx.Commit(ctx)
}

func (m *pgTxManager) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if !ok {
		return fmt.Errorf("no hay transacción activa en el contexto para revertir")
	}
	return tx.Rollback(ctx)
}

// Transaction implementa el bloque con recover() y rollback automático integrado.
func (m *pgTxManager) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	txCtx, err := m.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error al iniciar la transacción: %w", err)
	}

	committed := false
	defer func() {
		if p := recover(); p != nil {
			_ = m.Rollback(txCtx)
			panic(p) // Propaga el panic después de hacer el rollback seguro
		} else if !committed {
			_ = m.Rollback(txCtx) // Rollback si ocurre un error y no se hizo commit
		}
	}()

	if err := fn(txCtx); err != nil {
		return err // El defer detectará que 'committed' es falso y hará rollback
	}

	if err := m.Commit(txCtx); err != nil {
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}

	committed = true
	return nil
}
