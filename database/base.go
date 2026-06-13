package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Database is the driver-agnostic data access contract used by the DAO layer
type Database interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	RunInTransaction(ctx context.Context, fn func(context.Context) error) error
	DB() *sqlx.DB
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DatabaseManager holds the application data store
type DatabaseManager struct {
	DataDb Database
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewManager wraps a Database implementation for use by the application
func NewManager(dataDb Database) *DatabaseManager {
	return &DatabaseManager{DataDb: dataDb}
}
