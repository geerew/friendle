package database

import (
	"context"
	"testing"

	"github.com/geerew/friendle/utils/filesystem"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func TestNewSQLite(t *testing.T) {
	// Test successfully opening the database manager
	t.Run("success", func(t *testing.T) {
		fs := filesystem.New(afero.NewMemMapFs())

		dbManager, err := NewSQLite(&SQLiteConfig{
			DataDir: "./oc_data",
			FS:      fs,
			Testing: true,
		})

		require.NoError(t, err)
		require.NotNil(t, dbManager)
	})

	// Test error due to being unable to create data.db
	t.Run("write error", func(t *testing.T) {
		fs := filesystem.New(afero.NewReadOnlyFs(afero.NewMemMapFs()))

		dbManager, err := NewSQLite(&SQLiteConfig{
			DataDir: "./oc_data",
			FS:      fs,
			Testing: true,
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to create write database")
		require.Contains(t, err.Error(), "operation not permitted")
		require.Nil(t, dbManager)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func TestOpenSQLite(t *testing.T) {
	// Test successfully opening a sqlite connection
	t.Run("success", func(t *testing.T) {
		fs := filesystem.New(afero.NewMemMapFs())

		db, err := openSQLite(&sqliteConfig{
			DataDir:    "./oc_data",
			DSN:        "data.db",
			MigrateDir: "data",
			FS:         fs,
			Testing:    true,
		})

		require.NoError(t, err)
		require.NotNil(t, db)
	})

	// Test error due to being unable to create data.db
	t.Run("write error", func(t *testing.T) {
		fs := filesystem.New(afero.NewReadOnlyFs(afero.NewMemMapFs()))

		db, err := openSQLite(&sqliteConfig{
			DataDir:    "./oc_data",
			DSN:        "data.db",
			MigrateDir: "data",
			FS:         fs,
			Testing:    true,
		})

		require.Error(t, err)
		require.EqualError(t, err, "operation not permitted")
		require.Nil(t, db)
	})

	// Test error due to invalid migration directory
	t.Run("invalid migration", func(t *testing.T) {
		fs := filesystem.New(afero.NewMemMapFs())

		db, err := openSQLite(&sqliteConfig{
			DataDir:    "./oc_data",
			DSN:        "data.db",
			MigrateDir: "test",
			FS:         fs,
			Testing:    true,
		})

		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to run migrations in test")
		require.Contains(t, err.Error(), "test directory does not exist")
		require.Nil(t, db)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func TestQueryContext(t *testing.T) {
	// Test successfully querying multiple rows
	t.Run("simple", func(t *testing.T) {
		dbManager := setupTestDB(t)
		ctx := context.Background()

		_, err := dbManager.DataDb.ExecContext(ctx, "INSERT INTO test (name) VALUES ('test')")
		require.NoError(t, err)

		rows, err := dbManager.DataDb.QueryContext(ctx, "SELECT * FROM test")
		require.NoError(t, err)
		defer rows.Close()

		for rows.Next() {
			var id int
			var name string
			err = rows.Scan(&id, &name)
			require.NoError(t, err)
			require.Equal(t, 1, id)
			require.Equal(t, "test", name)
		}

		require.NoError(t, rows.Err())
	})

	// Test successfully querying multiple rows in a transaction
	t.Run("transaction", func(t *testing.T) {
		dbManager := setupTestDB(t)
		ctx := context.Background()

		_, err := dbManager.DataDb.ExecContext(ctx, "INSERT INTO test (name) VALUES ('test')")
		require.NoError(t, err)

		var id int
		var name string

		err = dbManager.DataDb.RunInTransaction(ctx, func(txCtx context.Context) error {
			rows, err := dbManager.DataDb.QueryContext(txCtx, "SELECT * FROM test")
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				err = rows.Scan(&id, &name)
				require.NoError(t, err)
			}

			return rows.Err()
		})

		require.NoError(t, err)
		require.Equal(t, "test", name)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func TestQueryRowContext(t *testing.T) {
	// Test successfully querying a single row
	t.Run("simple", func(t *testing.T) {
		dbManager := setupTestDB(t)
		ctx := context.Background()

		_, err := dbManager.DataDb.ExecContext(ctx, "INSERT INTO test (name) VALUES ('test')")
		require.NoError(t, err)

		var id int
		var name string
		err = dbManager.DataDb.QueryRowContext(ctx, "SELECT * FROM test").Scan(&id, &name)

		require.NoError(t, err)
		require.Equal(t, "test", name)
	})

	// Test successfully querying a single row in a transaction
	t.Run("transaction", func(t *testing.T) {
		dbManager := setupTestDB(t)
		ctx := context.Background()

		_, err := dbManager.DataDb.ExecContext(ctx, "INSERT INTO test (name) VALUES ('test')")
		require.NoError(t, err)

		var id int
		var name string

		err = dbManager.DataDb.RunInTransaction(ctx, func(txCtx context.Context) error {
			return dbManager.DataDb.QueryRowContext(txCtx, "SELECT * FROM test").Scan(&id, &name)
		})

		require.NoError(t, err)
		require.Equal(t, "test", name)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func TestExecContext(t *testing.T) {
	// Test successfully executing a non-query SQL statement
	t.Run("simple", func(t *testing.T) {
		dbManager := setupTestDB(t)
		ctx := context.Background()

		result, err := dbManager.DataDb.ExecContext(ctx, "INSERT INTO test (name) VALUES ('test')")
		require.NoError(t, err)

		rowAffected, err := result.RowsAffected()
		require.NoError(t, err)
		require.Equal(t, int64(1), rowAffected)
	})

	// Test successfully executing a non-query SQL statement in a transaction
	t.Run("transaction", func(t *testing.T) {
		dbManager := setupTestDB(t)
		ctx := context.Background()

		err := dbManager.DataDb.RunInTransaction(ctx, func(txCtx context.Context) error {
			_, err := dbManager.DataDb.ExecContext(txCtx, "INSERT INTO test (name) VALUES ('test')")

			return err
		})
		require.NoError(t, err)

		var count int
		err = dbManager.DataDb.QueryRowContext(ctx, "SELECT COUNT(*) FROM test").Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func TestGetContext(t *testing.T) {
	// Test successfully getting a single record
	t.Run("simple", func(t *testing.T) {
		dbManager := setupTestDB(t)
		ctx := context.Background()

		_, err := dbManager.DataDb.ExecContext(ctx, "INSERT INTO test (name) VALUES ('test1'), ('test2')")
		require.NoError(t, err)

		record := &struct {
			Id   int    `db:"id"`
			Name string `db:"name"`
		}{}

		err = dbManager.DataDb.GetContext(ctx, record, "SELECT * FROM test WHERE name = ?", "test1")
		require.NoError(t, err)
		require.Equal(t, 1, record.Id)
		require.Equal(t, "test1", record.Name)
	})

	// Test successfully getting a single record in a transaction
	t.Run("transaction", func(t *testing.T) {
		dbManager := setupTestDB(t)
		ctx := context.Background()

		_, err := dbManager.DataDb.ExecContext(ctx, "INSERT INTO test (name) VALUES ('test1'), ('test2')")
		require.NoError(t, err)

		record := &struct {
			Id   int    `db:"id"`
			Name string `db:"name"`
		}{}

		err = dbManager.DataDb.RunInTransaction(ctx, func(txCtx context.Context) error {
			return dbManager.DataDb.GetContext(txCtx, record, "SELECT * FROM test WHERE name = ?", "test1")
		})

		require.NoError(t, err)
		require.Equal(t, 1, record.Id)
		require.Equal(t, "test1", record.Name)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func TestSelectContext(t *testing.T) {
	// Test successfully selecting multiple records
	t.Run("simple", func(t *testing.T) {
		dbManager := setupTestDB(t)
		ctx := context.Background()

		records := []*struct {
			Id   int    `db:"id"`
			Name string `db:"name"`
		}{}

		_, err := dbManager.DataDb.ExecContext(ctx, "INSERT INTO test (name) VALUES ('test'), ('test2')")
		require.NoError(t, err)

		err = dbManager.DataDb.SelectContext(ctx, &records, "SELECT * FROM test")
		require.NoError(t, err)
		require.Len(t, records, 2)
		require.Equal(t, 1, records[0].Id)
		require.Equal(t, "test", records[0].Name)
		require.Equal(t, 2, records[1].Id)
		require.Equal(t, "test2", records[1].Name)
	})

	// Test successfully selecting multiple records in a transaction
	t.Run("transaction", func(t *testing.T) {
		dbManager := setupTestDB(t)
		ctx := context.Background()

		_, err := dbManager.DataDb.ExecContext(ctx, "INSERT INTO test (name) VALUES ('test'), ('test2')")
		require.NoError(t, err)

		records := []*struct {
			Id   int    `db:"id"`
			Name string `db:"name"`
		}{}

		err = dbManager.DataDb.RunInTransaction(ctx, func(txCtx context.Context) error {
			return dbManager.DataDb.SelectContext(txCtx, &records, "SELECT * FROM test")
		})

		require.NoError(t, err)
		require.Len(t, records, 2)
		require.Equal(t, 1, records[0].Id)
		require.Equal(t, "test", records[0].Name)
		require.Equal(t, 2, records[1].Id)
		require.Equal(t, "test2", records[1].Name)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// setupTestDB opens an in-memory manager and creates a simple test table
func setupTestDB(t *testing.T) *DatabaseManager {
	t.Helper()

	fs := filesystem.New(afero.NewMemMapFs())

	dbManager, err := NewSQLite(&SQLiteConfig{
		DataDir: "./oc_data",
		FS:      fs,
		Testing: true,
	})
	require.NoError(t, err)
	require.NotNil(t, dbManager)

	_, err = dbManager.DataDb.ExecContext(context.Background(), "CREATE TABLE test (id INTEGER PRIMARY KEY, name TEXT)")
	require.NoError(t, err)

	return dbManager
}
