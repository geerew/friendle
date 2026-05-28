package session

import (
	"context"
	"testing"

	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/utils/filesystem"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func setup(tb testing.TB) (database.Database, context.Context) {
	tb.Helper()

	dbManager, err := database.NewSQLiteManager(&database.DatabaseManagerConfig{
		DataDir: "./oc_data",
		FS:   filesystem.New(afero.NewMemMapFs()),
		Testing: true,
	})

	require.NoError(tb, err)
	require.NotNil(tb, dbManager)

	return dbManager.DataDb, context.Background()
}
