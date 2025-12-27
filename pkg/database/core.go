// Package database provides PostgreSQL database connection and query functionality.
package database

import (
	"context"
	"errors"

	"backend.atomicledger.com/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/oops"
)

// Ensure DatabaseInterface is implemented by Database
var _ DatabaseInterface = (*Database)(nil)

const (
	sqlPreviewMaxLen = 100
)

// Database wraps a PostgreSQL connection pool with logging functionality.
type Database struct {
	Pool   PoolInterface
	logger logger.Logger
}

// NewConnection creates a new database connection pool with proper error handling and logging.
func NewConnection(ctx context.Context, connString string, log logger.Logger) (*Database, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, oops.
			Code("db_pool_creation_failed").
			With("error_type", "connection_pool").
			Wrapf(err, "failed to create database connection pool")
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close() // Clean up the pool before returning error
		return nil, oops.
			Code("db_ping_failed").
			With("error_type", "connection_test").
			Wrapf(err, "failed to ping database after pool creation")
	}

	db := &Database{
		Pool:   pool,
		logger: log.With("component", "database"),
	}

	// Log successful connection with pool configuration
	poolConfig := pool.Config()
	db.logger.Info("database connection established",
		"max_connections", poolConfig.MaxConns,
		"min_connections", poolConfig.MinConns,
	)

	return db, nil
}

// Close closes the database connection pool.
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		db.logger.Info("database connection pool closed")
	}
}

// GetPool returns the underlying PostgreSQL connection pool.
func (db *Database) GetPool() PoolInterface {
	return db.Pool
}

// Query executes a query that returns multiple rows.
func (db *Database) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	rows, err := db.Pool.Query(ctx, sql, args...)
	if err != nil {
		db.logger.Error("query execution failed",
			"error", err,
			"sql_preview", truncateSQL(sql),
		)
		return nil, oops.
			Code("db_query_failed").
			With("operation", "query").
			With("sql_preview", truncateSQL(sql)).
			Wrapf(err, "database query failed")
	}
	return rows, nil
}

// QueryRow executes a query that returns a single row.
func (db *Database) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}

// QueryRowScan executes a query that returns a single row and scans it using the provided function.
func (db *Database) QueryRowScan(ctx context.Context, scanFunc func(row pgx.Row) error, sql string, args ...any) error {
	row := db.Pool.QueryRow(ctx, sql, args...)
	if err := scanFunc(row); err != nil {
		db.logger.Error("query row scan failed",
			"error", err,
			"sql_preview", truncateSQL(sql),
		)

		return oops.
			Code("db_query_row_scan_failed").
			With("operation", "query_row_scan").
			With("sql_preview", truncateSQL(sql)).
			Wrapf(err, "database query row scan failed")
	}

	return nil
}

// Exec executes a query that doesn't return rows (INSERT, UPDATE, DELETE).
func (db *Database) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	tag, err := db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		db.logger.Error("exec operation failed",
			"error", err,
			"sql_preview", truncateSQL(sql),
		)
		return tag, oops.
			Code("db_exec_failed").
			With("operation", "exec").
			With("sql_preview", truncateSQL(sql)).
			Wrapf(err, "database exec failed")
	}

	db.logger.Debug("exec operation completed",
		"rows_affected", tag.RowsAffected(),
	)

	return tag, nil
}

// HealthCheck performs a health check on the database connection.
func (db *Database) HealthCheck(ctx context.Context) error {
	if db.Pool == nil {
		return oops.
			Code("db_health_check_failed").
			With("error_type", "nil_pool").
			Wrapf(errors.New("database pool is nil"), "database health check failed")
	}
	if err := db.Pool.Ping(ctx); err != nil {
		return oops.
			Code("db_health_check_failed").
			Wrapf(err, "database health check failed")
	}
	return nil
}

// Shutdown gracefully closes the database connection pool.
func (db *Database) Shutdown(_ context.Context) error {
	db.logger.Info("shutting down database connection")
	if db.Pool != nil {
		db.Pool.Close()
	}
	return nil
}

// truncateSQL truncates SQL string for safe logging.
func truncateSQL(sql string) string {
	if len(sql) <= sqlPreviewMaxLen {
		return sql
	}
	return sql[:sqlPreviewMaxLen] + "..."
}
