package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/geerew/friendle/migrations"
	"github.com/geerew/friendle/utils/filesystem"
	"github.com/geerew/friendle/utils/security"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	migrateDirData        = "data"
	modeReadWrite         = "rwc"
	modeReadOnly          = "ro"
	dsnData               = "data.db"
	defaultMaxLockRetries = 5
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var (
	gooseOnce             sync.Once
	defaultRetryIntervals = []time.Duration{
		50 * time.Millisecond,
		100 * time.Millisecond,
		150 * time.Millisecond,
		200 * time.Millisecond,
		300 * time.Millisecond,
		400 * time.Millisecond,
		500 * time.Millisecond,
		700 * time.Millisecond,
		1000 * time.Millisecond,
	}
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SQLiteConfig holds settings for opening the SQLite data database
type SQLiteConfig struct {
	DataDir string
	FS      *filesystem.FS
	Testing bool
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewSQLite opens the data database with separate read and write pools
func NewSQLite(config *SQLiteConfig) (*DatabaseManager, error) {
	dsnName := dsnName(dsnData, config.Testing)

	writeDb, err := openSQLite(&sqliteConfig{
		DataDir:    config.DataDir,
		DSN:        dsnName,
		MigrateDir: migrateDirData,
		FS:         config.FS,
		Testing:    config.Testing,
		Mode:       modeReadWrite,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create write database: %w", err)
	}

	configurePool(writeDb, 1, 1)

	readDb, err := openSQLite(&sqliteConfig{
		DataDir: config.DataDir,
		DSN:     dsnName,
		FS:      config.FS,
		Testing: config.Testing,
		Mode:    modeReadOnly,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create read database: %w", err)
	}

	configurePool(readDb, 10, 5)

	return NewManager(&sqliteDB{
		read:  readDb,
		write: writeDb,
	}), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// sqliteDB wraps separate SQLite read and write pools
type sqliteDB struct {
	read  *sqlx.DB
	write *sqlx.DB
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ExecContext executes a non-query SQL statement on the write pool with lock retries
func (db *sqliteDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx := txFromContext(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}

	var (
		res sql.Result
		err error
	)

	for attempt := 0; attempt <= defaultMaxLockRetries; attempt++ {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}

		res, err = db.write.ExecContext(ctx, query, args...)
		if err == nil {
			return res, nil
		}

		if !isLockError(err) {
			return res, err
		}

		if attempt == defaultMaxLockRetries {
			break
		}

		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case <-time.After(retryInterval(attempt)):
		}
	}

	return res, fmt.Errorf("%w after %d retries", err, defaultMaxLockRetries)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// QueryContext executes a query on the read pool, or the active transaction when present
func (db *sqliteDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if tx := txFromContext(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}

	return db.read.QueryContext(ctx, query, args...)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// QueryRowContext executes a single-row query on the read pool, or the active transaction
func (db *sqliteDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if tx := txFromContext(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}

	return db.read.QueryRowContext(ctx, query, args...)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetContext loads one row into dest from the read pool, or the active transaction
func (db *sqliteDB) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	if tx := txFromContext(ctx); tx != nil {
		return tx.GetContext(ctx, dest, query, args...)
	}

	return db.read.GetContext(ctx, dest, query, args...)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SelectContext loads rows into dest from the read pool, or the active transaction
func (db *sqliteDB) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	if tx := txFromContext(ctx); tx != nil {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return db.read.SelectContext(ctx, dest, query, args...)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RunInTransaction runs fn inside a write transaction with lock retries
func (db *sqliteDB) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	if txFromContext(ctx) != nil {
		return fn(ctx)
	}

	var lastErr error

	for attempt := 0; attempt <= defaultMaxLockRetries; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		sqlxTx, err := db.write.BeginTxx(ctx, nil)
		if err != nil {
			if !isLockError(err) {
				return err
			}

			lastErr = err

			if attempt < defaultMaxLockRetries {
				if err := sleepWithContext(ctx, retryInterval(attempt)); err != nil {
					return err
				}

				continue
			}

			break
		}

		txCtx := context.WithValue(ctx, txKey{}, sqlxTx)
		err = fn(txCtx)

		if err == nil {
			if commitErr := sqlxTx.Commit(); commitErr != nil {
				sqlxTx.Rollback()

				if isLockError(commitErr) && attempt < defaultMaxLockRetries {
					lastErr = commitErr

					if err := sleepWithContext(ctx, retryInterval(attempt)); err != nil {
						return err
					}

					continue
				}

				return commitErr
			}

			return nil
		}

		sqlxTx.Rollback()
		lastErr = err

		if !isLockError(err) || attempt >= defaultMaxLockRetries {
			return err
		}

		if err := sleepWithContext(ctx, retryInterval(attempt)); err != nil {
			return err
		}
	}

	return fmt.Errorf("%w after %d retries", lastErr, defaultMaxLockRetries)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DB returns the underlying write pool
func (db *sqliteDB) DB() *sqlx.DB {
	return db.write
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// txKey is the context key used to carry an active transaction
type txKey struct{}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// sqliteConfig holds settings for opening one SQLite pool
type sqliteConfig struct {
	DataDir    string
	DSN        string
	MigrateDir string
	FS         *filesystem.FS
	Mode       string
	Testing    bool
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// txFromContext returns the active transaction stored in ctx, or nil
func txFromContext(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// openSQLite opens a SQLite pool, runs migrations when configured, and verifies connectivity
func openSQLite(config *sqliteConfig) (*sqlx.DB, error) {
	if err := config.FS.MkdirAll(config.DataDir, os.ModePerm); err != nil {
		return nil, err
	}

	pragmaParts := []string{
		"cache=shared",
		"_busy_timeout=10000",
		"_journal_mode=WAL",
		"_journal_size_limit=200000000",
		"_synchronous=NORMAL",
		"_foreign_keys=1",
		"_cache_size=-16000",
	}

	if config.Mode != "" {
		pragmaParts = append([]string{fmt.Sprintf("mode=%s", config.Mode)}, pragmaParts...)
	}

	pragma := strings.Join(pragmaParts, "&")
	dsn := fmt.Sprintf("file:%s?%s", filepath.Join(config.DataDir, config.DSN), pragma)

	if config.Testing {
		dsn += "&mode=memory"
	}

	conn, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(1)
	conn.SetMaxOpenConns(1)

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if config.MigrateDir != "" {
		if err := migrate(conn, config.MigrateDir); err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// migrate runs goose migrations for the given directory
func migrate(db *sqlx.DB, migrateDir string) error {
	gooseOnce.Do(func() {
		goose.SetLogger(goose.NopLogger())
		goose.SetBaseFS(migrations.EmbedMigrations)
		if err := goose.SetDialect("sqlite3"); err != nil {
			panic(fmt.Errorf("failed to set goose dialect: %w", err))
		}
	})

	if err := goose.Up(db.DB, migrateDir); err != nil {
		return fmt.Errorf("failed to run migrations in %s: %w", migrateDir, err)
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// isLockError reports whether err is a SQLite lock contention error
func isLockError(err error) bool {
	if err == nil {
		return false
	}

	s := err.Error()

	return strings.Contains(s, "database is locked") || strings.Contains(s, "table is locked")
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// sleepWithContext waits for duration or until ctx is cancelled
func sleepWithContext(ctx context.Context, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// retryInterval returns the backoff delay for the given attempt
func retryInterval(attempt int) time.Duration {
	if attempt < 0 || attempt >= len(defaultRetryIntervals) {
		return defaultRetryIntervals[len(defaultRetryIntervals)-1]
	}

	return defaultRetryIntervals[attempt]
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// dsnName returns the database file name, randomised in tests to avoid collisions
func dsnName(baseName string, testing bool) string {
	if testing {
		return fmt.Sprintf("%s_memdb_%s", strings.TrimSuffix(baseName, ".db"), security.PseudorandomString(8))
	}

	return baseName
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// configurePool sets connection pool limits on db
func configurePool(db *sqlx.DB, maxOpen, maxIdle int) {
	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(time.Hour)
	db.SetConnMaxIdleTime(10 * time.Minute)
}
