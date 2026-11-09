package ports

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Executor is the minimal pgx contract required by repositories.
// Satisfied by *pgxpool.Pool, pgx.Tx, and sessionExecutor — making
// repository code agnostic to whether it runs inside a transaction or not.
type Executor interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
}

// Resolver returns the active executor for the current request scope.
type Resolver interface {
	DB(ctx context.Context) Executor
}

// TransactionManager extends Resolver with transaction lifecycle management.
type TransactionManager interface {
	Resolver
	Begin(ctx context.Context) (context.Context, error)
	BeginTx(ctx context.Context, opts pgx.TxOptions) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
