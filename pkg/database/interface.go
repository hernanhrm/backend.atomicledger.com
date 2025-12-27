package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// PoolInterface defines the interface for database pool operations
// This allows for mocking in tests.
type PoolInterface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Close()
}

// DatabaseInterface defines the interface for database operations
type DatabaseInterface interface {
	GetPool() PoolInterface
	Close()
	Shutdown(ctx context.Context) error
}
