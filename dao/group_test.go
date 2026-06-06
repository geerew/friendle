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

func Test_GetGroup(t *testing.T) {
	// Test successfully retrieving a group record
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		dbOpts := NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: group.ID})
		record, err := dao.GetGroup(ctx, dbOpts)
		require.NoError(t, err)
		require.Equal(t, group.ID, record.ID)
		require.Equal(t, "Friends", record.Name)
	})

	// Test no error when retrieving a non-existent group record
	t.Run("not found", func(t *testing.T) {
		dao, ctx := setup(t)

		record, err := dao.GetGroup(ctx, NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: "missing"}))
		require.NoError(t, err)
		require.Nil(t, record)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_ListGroups(t *testing.T) {
	// Test successfully listing groups for a member
	t.Run("member groups", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))
		require.NoError(t, dao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    userID,
			GroupRole: types.GroupRoleAdmin,
		}))

		other := &models.Group{Name: "Work", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, other))
		require.NoError(t, dao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   other.ID,
			UserID:    userID,
			GroupRole: types.GroupRoleUser,
		}))

		where, err := MemberGroupsWhere(userID)
		require.NoError(t, err)

		groups, err := dao.ListGroups(ctx, NewOptions().WithWhere(where))
		require.NoError(t, err)
		require.Len(t, groups, 2)
		require.Equal(t, "Friends", groups[0].Name)
		require.Equal(t, "Work", groups[1].Name)
	})

	// Test empty list when the user has no memberships
	t.Run("no memberships", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		where, err := MemberGroupsWhere(userID)
		require.NoError(t, err)

		groups, err := dao.ListGroups(ctx, NewOptions().WithWhere(where))
		require.NoError(t, err)
		require.Empty(t, groups)
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
