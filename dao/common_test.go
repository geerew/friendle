package dao

import (
	"context"
	"testing"

	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/filesystem"
	"github.com/geerew/friendle/utils/types"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func setup(tb testing.TB) (*DAO, context.Context) {
	tb.Helper()

	dbManager, err := database.NewSQLite(&database.SQLiteConfig{
		DataDir: tb.TempDir(),
		FS:      filesystem.New(afero.NewMemMapFs()),
		Testing: true,
	})

	require.NoError(tb, err)
	require.NotNil(tb, dbManager)

	dao := &DAO{db: dbManager.DataDb}

	user := &models.User{
		Username:     "test-user",
		DisplayName:  "Test User",
		PasswordHash: "test-password",
		SiteRole:     types.SiteRoleAdmin,
	}
	require.NoError(tb, dao.CreateUser(context.Background(), user))

	ctx := context.Background()
	principal := types.Principal{
		UserID:   user.ID,
		SiteRole: user.SiteRole,
	}
	ctx = context.WithValue(ctx, types.PrincipalContextKey, principal)

	return dao, ctx
}
