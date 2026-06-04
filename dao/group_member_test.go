package dao

import (
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_CreateGroupMember(t *testing.T) {
	// Test successfully inserting a group member
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Member Group", user.ID)

		member := &models.GroupMember{
			GroupID: group.ID,
			UserID:  user.ID,
			GroupRole: types.GroupRoleAdmin,
		}
		require.NoError(t, dao.CreateGroupMember(ctx, member))
		require.NotEmpty(t, member.ID)
	})

	// Test error due to duplicate membership
	t.Run("duplicate", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Dup Group", user.ID)
		createTestGroupMember(t, dao, ctx, group.ID, user.ID, types.GroupRoleAdmin)

		err := dao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    user.ID,
			GroupRole: types.GroupRoleUser,
		})
		require.ErrorContains(t, err, "UNIQUE constraint failed")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_GetGroupMember(t *testing.T) {
	// Test successfully retrieving a group member
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Get Member Group", user.ID)
		member := createTestGroupMember(t, dao, ctx, group.ID, user.ID, types.GroupRoleAdmin)

		record, err := dao.GetGroupMember(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: member.ID}))
		require.NoError(t, err)
		require.Equal(t, member.ID, record.ID)
		require.Equal(t, group.ID, record.GroupID)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_ListGroupMembers(t *testing.T) {
	// Test successfully listing group members for a group
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		other := createTestUser(t, dao, ctx, "other-user")
		group := createTestGroup(t, dao, ctx, "List Members", user.ID)
		createTestGroupMember(t, dao, ctx, group.ID, user.ID, types.GroupRoleAdmin)
		createTestGroupMember(t, dao, ctx, group.ID, other.ID, types.GroupRoleUser)

		members, err := dao.ListGroupMembers(ctx, NewOptions().WithWhere(squirrel.Eq{
			models.GROUP_MEMBER_GROUP_ID: group.ID,
		}))
		require.NoError(t, err)
		require.Len(t, members, 2)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_CountGroupMembers(t *testing.T) {
	// Test successfully counting group members
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Count Members", user.ID)
		createTestGroupMember(t, dao, ctx, group.ID, user.ID, types.GroupRoleAdmin)

		count, err := dao.CountGroupMembers(ctx, NewOptions().WithWhere(squirrel.Eq{
			models.GROUP_MEMBER_GROUP_ID: group.ID,
		}))
		require.NoError(t, err)
		require.Equal(t, 1, count)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_UpdateGroupMember(t *testing.T) {
	// Test successfully updating a group member
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Update Member", user.ID)
		member := createTestGroupMember(t, dao, ctx, group.ID, user.ID, types.GroupRoleUser)

		member.GroupRole = types.GroupRoleAdmin
		member.TimesPicked = 2
		require.NoError(t, dao.UpdateGroupMember(ctx, member))

		record, err := dao.GetGroupMember(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: member.ID}))
		require.NoError(t, err)
		require.Equal(t, types.GroupRoleAdmin, record.GroupRole)
		require.Equal(t, 2, record.TimesPicked)
	})

	// Test error due to missing ID
	t.Run("missing id", func(t *testing.T) {
		dao, ctx := setup(t)

		require.ErrorIs(t, dao.UpdateGroupMember(ctx, &models.GroupMember{}), utils.ErrId)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_DeleteGroupMembers(t *testing.T) {
	// Test successfully deleting a group member
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		group := createTestGroup(t, dao, ctx, "Remove Member", user.ID)
		member := createTestGroupMember(t, dao, ctx, group.ID, user.ID, types.GroupRoleAdmin)

		require.NoError(t, dao.DeleteGroupMembers(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: member.ID})))

		record, err := dao.GetGroupMember(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: member.ID}))
		require.NoError(t, err)
		require.Nil(t, record)
	})

	// Test error due to missing where clause
	t.Run("missing where", func(t *testing.T) {
		dao, ctx := setup(t)

		require.ErrorIs(t, dao.DeleteGroupMembers(ctx, nil), utils.ErrWhere)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_ListMemberGroupIDsForUser(t *testing.T) {
	// Test successfully listing member group IDs for a user
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		groupA := createTestGroup(t, dao, ctx, "A", user.ID)
		groupB := createTestGroup(t, dao, ctx, "B", user.ID)
		createTestGroup(t, dao, ctx, "C", user.ID)
		createTestGroupMember(t, dao, ctx, groupA.ID, user.ID, types.GroupRoleAdmin)
		createTestGroupMember(t, dao, ctx, groupB.ID, user.ID, types.GroupRoleUser)

		ids, err := dao.ListMemberGroupIDsForUser(ctx, user.ID, []string{groupA.ID, groupB.ID, "missing"})
		require.NoError(t, err)
		require.ElementsMatch(t, []string{groupA.ID, groupB.ID}, ids)
	})

	// Test empty input returns empty slice
	t.Run("empty group ids", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)

		ids, err := dao.ListMemberGroupIDsForUser(ctx, user.ID, nil)
		require.NoError(t, err)
		require.Empty(t, ids)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_MembersWithMinPicks(t *testing.T) {
	// Test successfully returning members tied for fewest picks
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		other := createTestUser(t, dao, ctx, "picker-b")
		group := createTestGroup(t, dao, ctx, "Picker Group", user.ID)

		memberA := createTestGroupMember(t, dao, ctx, group.ID, user.ID, types.GroupRoleUser)
		memberB := createTestGroupMember(t, dao, ctx, group.ID, other.ID, types.GroupRoleUser)

		memberA.TimesPicked = 1
		require.NoError(t, dao.UpdateGroupMember(ctx, memberA))

		members, err := dao.MembersWithMinPicks(ctx, group.ID)
		require.NoError(t, err)
		require.Len(t, members, 1)
		require.Equal(t, memberB.ID, members[0].ID)
	})
}
