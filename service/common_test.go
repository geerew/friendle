package service

import (
	"context"
	"testing"
	"time"

	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/filesystem"
	"github.com/geerew/friendle/utils/types"
	"github.com/geerew/friendle/utils/words"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// roundTestEnv holds a service test fixture with database access for seeding
type roundTestEnv struct {
	svc *Service
	dao *dao.DAO
	ctx context.Context
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// setupRoundTest creates an in-memory service backed by SQLite
func setupRoundTest(tb testing.TB) *roundTestEnv {
	tb.Helper()

	dbManager, err := database.NewSQLite(&database.SQLiteConfig{
		DataDir: tb.TempDir(),
		FS:      filesystem.New(afero.NewMemMapFs()),
		Testing: true,
	})
	require.NoError(tb, err)

	dict, err := words.New()
	require.NoError(tb, err)

	appDao := dao.New(dbManager.DataDb)

	return &roundTestEnv{
		svc: New(dbManager.DataDb, dict),
		dao: appDao,
		ctx: context.Background(),
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ctxWithPrincipal returns a context carrying the given caller identity
func ctxWithPrincipal(ctx context.Context, userID string, siteRole types.SiteRole) context.Context {
	return context.WithValue(ctx, types.PrincipalContextKey, types.Principal{
		UserID:   userID,
		SiteRole: siteRole,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createTestUser inserts a user with the given username
func createTestUser(tb testing.TB, env *roundTestEnv, username string, siteRole types.SiteRole) *models.User {
	tb.Helper()

	user := &models.User{
		Username:     username,
		DisplayName:  username,
		PasswordHash: "test-password",
		SiteRole:     siteRole,
	}
	require.NoError(tb, env.dao.CreateUser(env.ctx, user))

	return user
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createTestGroup inserts a group owned by createdBy
func createTestGroup(tb testing.TB, env *roundTestEnv, name, createdBy string) *models.Group {
	tb.Helper()

	group := &models.Group{Name: name, CreatedBy: createdBy}
	require.NoError(tb, env.dao.CreateGroup(env.ctx, group))

	return group
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createTestGroupMember inserts a group membership row
func createTestGroupMember(tb testing.TB, env *roundTestEnv, groupID, userID string, role types.GroupRole) *models.GroupMember {
	tb.Helper()

	member := &models.GroupMember{
		GroupID:   groupID,
		UserID:    userID,
		GroupRole: role,
	}
	require.NoError(tb, env.dao.CreateGroupMember(env.ctx, member))

	return member
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createTestRound inserts a round for today with the given status
func createTestRound(tb testing.TB, env *roundTestEnv, groupID, pickerID string, status types.RoundStatus, word *string) *models.Round {
	tb.Helper()

	round := &models.Round{
		GroupID:      groupID,
		RoundDate:    time.Now().Format("2006-01-02"),
		PickerUserID: pickerID,
		Status:       status,
		WordPlain:    word,
	}
	require.NoError(tb, env.dao.CreateRound(env.ctx, round))

	return round
}
