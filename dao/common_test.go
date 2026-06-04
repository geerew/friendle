package dao

import (
	"context"
	"testing"

	"github.com/Masterminds/squirrel"
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// principalUser returns the default user created by setup
func principalUser(tb testing.TB, dao *DAO, ctx context.Context) *models.User {
	tb.Helper()

	user, err := dao.GetUser(ctx, NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: "test-user"}))
	require.NoError(tb, err)
	require.NotNil(tb, user)

	return user
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createTestUser inserts a user with the given username
func createTestUser(tb testing.TB, dao *DAO, ctx context.Context, username string) *models.User {
	tb.Helper()

	user := &models.User{
		Username:     username,
		DisplayName:  username,
		PasswordHash: "test-password",
		SiteRole:     types.SiteRoleUser,
	}
	require.NoError(tb, dao.CreateUser(ctx, user))

	return user
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createTestGroup inserts a group owned by createdBy
func createTestGroup(tb testing.TB, dao *DAO, ctx context.Context, name, createdBy string) *models.Group {
	tb.Helper()

	group := &models.Group{Name: name, CreatedBy: createdBy}
	require.NoError(tb, dao.CreateGroup(ctx, group))

	return group
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createTestGroupMember inserts a group membership row
func createTestGroupMember(tb testing.TB, dao *DAO, ctx context.Context, groupID, userID string, role types.GroupRole) *models.GroupMember {
	tb.Helper()

	member := &models.GroupMember{
		GroupID:   groupID,
		UserID:    userID,
		GroupRole: role,
	}
	require.NoError(tb, dao.CreateGroupMember(ctx, member))

	return member
}
