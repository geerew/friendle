package dao

import (
	"context"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_CreateGroup(t *testing.T) {
	// Test successfully creating a group record
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))
		require.NotEmpty(t, group.ID)
	})

	// Test error due to duplicate name
	t.Run("duplicate", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		duplicate := &models.Group{Name: "Friends", CreatedBy: userID}
		require.ErrorContains(t, dao.CreateGroup(ctx, duplicate), "UNIQUE constraint failed: "+models.GROUP_TABLE+"."+models.GROUP_NAME)
	})

	// Test error due to nil pointer
	t.Run("nil pointer", func(t *testing.T) {
		dao, ctx := setup(t)
		require.ErrorIs(t, dao.CreateGroup(ctx, nil), utils.ErrNilPtr)
	})

	// Test error due to empty name
	t.Run("empty name", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "", CreatedBy: userID}
		require.ErrorIs(t, dao.CreateGroup(ctx, group), utils.ErrGroupName)
	})

	// Test error due to empty creator
	t.Run("empty creator", func(t *testing.T) {
		dao, ctx := setup(t)

		group := &models.Group{Name: "Friends", CreatedBy: ""}
		require.ErrorIs(t, dao.CreateGroup(ctx, group), utils.ErrUserId)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_CreateGroupMember(t *testing.T) {
	// Test successfully creating a group member record
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		member := &models.GroupMember{
			GroupID:   group.ID,
			UserID:    userID,
			GroupRole: types.GroupRoleAdmin,
		}
		require.NoError(t, dao.CreateGroupMember(ctx, member))
		require.NotEmpty(t, member.ID)
	})

	// Test error due to nil pointer
	t.Run("nil pointer", func(t *testing.T) {
		dao, ctx := setup(t)
		require.ErrorIs(t, dao.CreateGroupMember(ctx, nil), utils.ErrNilPtr)
	})

	// Test error due to empty group id
	t.Run("empty group id", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		member := &models.GroupMember{
			GroupID:   "",
			UserID:    userID,
			GroupRole: types.GroupRoleAdmin,
		}
		require.ErrorIs(t, dao.CreateGroupMember(ctx, member), utils.ErrGroupId)
	})

	// Test error due to empty user id
	t.Run("empty user id", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		member := &models.GroupMember{
			GroupID:   group.ID,
			UserID:    "",
			GroupRole: types.GroupRoleAdmin,
		}
		require.ErrorIs(t, dao.CreateGroupMember(ctx, member), utils.ErrUserId)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// testUserID returns the seeded test user's id
func testUserID(t *testing.T, dao *DAO, ctx context.Context) string {
	t.Helper()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: "test-user"})
	user, err := dao.GetUser(ctx, dbOpts)
	require.NoError(t, err)

	return user.ID
}
