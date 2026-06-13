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

func Test_GetGroup(t *testing.T) {
	// Test successfully retrieving a group record
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		dbOpts := NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: group.ID})
		require.NoError(t, dao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    userID,
			GroupRole: types.GroupRoleAdmin,
		}))

		record, err := dao.GetGroup(ctx, dbOpts)
		require.NoError(t, err)
		require.Equal(t, group.ID, record.ID)
		require.Equal(t, "Friends", record.Name)
		require.Equal(t, 1, record.MemberCount)
	})

	// Test member count is zero when a group has no members
	t.Run("no members", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Lonely", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		dbOpts := NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: group.ID})
		record, err := dao.GetGroup(ctx, dbOpts)
		require.NoError(t, err)
		require.Equal(t, 0, record.MemberCount)
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

		names := make([]string, len(groups))
		for i, record := range groups {
			names[i] = record.Name
			require.Equal(t, 1, record.MemberCount)
		}

		require.ElementsMatch(t, []string{"Friends", "Work"}, names)
	})

	// Test successfully filtering groups by a WHERE clause
	t.Run("where filter", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		friends := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, friends))

		work := &models.Group{Name: "Work", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, work))

		groups, err := dao.ListGroups(ctx, NewOptions().WithWhere(squirrel.Like{"LOWER(" + models.GROUP_TABLE_NAME + ")": "%friend%"}))
		require.NoError(t, err)
		require.Len(t, groups, 1)
		require.Equal(t, "Friends", groups[0].Name)
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

func Test_DeleteGroups(t *testing.T) {
	// Test successfully deleting a group record
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		opts := NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: group.ID})
		require.NoError(t, dao.DeleteGroups(ctx, opts))

		records, err := dao.ListGroups(ctx, opts)
		require.NoError(t, err)
		require.Empty(t, records)
	})

	// Test no error when deleting a non-existent group record
	t.Run("not found", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		opts := NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: "non-existent"})
		require.NoError(t, dao.DeleteGroups(ctx, opts))

		records, err := dao.ListGroups(ctx, nil)
		require.NoError(t, err)
		require.Len(t, records, 1)
		require.Equal(t, group.ID, records[0].ID)
	})

	// Test error due to missing where clause
	t.Run("missing where", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		require.ErrorIs(t, dao.DeleteGroups(ctx, nil), utils.ErrWhere)

		records, err := dao.ListGroups(ctx, nil)
		require.NoError(t, err)
		require.Len(t, records, 1)
		require.Equal(t, group.ID, records[0].ID)
	})

	// Test deleting a group cascades to related records
	t.Run("cascade", func(t *testing.T) {
		dao, ctx := setup(t)

		adminID := testUserID(t, dao, ctx)

		member := &models.User{
			Username:     "group-member",
			DisplayName:  "Group Member",
			PasswordHash: "test-password",
			SiteRole:     types.SiteRoleUser,
		}
		require.NoError(t, dao.CreateUser(ctx, member))

		pendingUser := &models.User{
			Username:     "pending-user",
			DisplayName:  "Pending User",
			PasswordHash: "test-password",
			SiteRole:     types.SiteRoleUser,
		}
		require.NoError(t, dao.CreateUser(ctx, pendingUser))

		rejectedUser := &models.User{
			Username:     "rejected-user",
			DisplayName:  "Rejected User",
			PasswordHash: "test-password",
			SiteRole:     types.SiteRoleUser,
		}
		require.NoError(t, dao.CreateUser(ctx, rejectedUser))

		group := &models.Group{Name: "Cascade", CreatedBy: adminID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		require.NoError(t, dao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    adminID,
			GroupRole: types.GroupRoleAdmin,
		}))
		require.NoError(t, dao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    member.ID,
			GroupRole: types.GroupRoleUser,
		}))
		require.NoError(t, dao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: group.ID,
			UserID:  pendingUser.ID,
			Status:  types.JoinPending,
		}))
		require.NoError(t, dao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: group.ID,
			UserID:  rejectedUser.ID,
			Status:  types.JoinRejected,
		}))

		round := &models.Round{
			GroupID:      group.ID,
			RoundDate:    "2026-05-28",
			PickerUserID: adminID,
			Status:       types.RoundActive,
		}
		require.NoError(t, dao.CreateRound(ctx, round))
		require.NoError(t, dao.CreateRoundMember(ctx, &models.RoundMember{
			RoundID: round.ID,
			UserID:  adminID,
		}))

		guess := &models.RoundMemberGuess{
			RoundID: round.ID,
			UserID:  adminID,
			Attempt: 1,
			Word:    "hello",
			Result:  types.TileStates{types.TileAbsent, types.TileAbsent, types.TileAbsent, types.TileAbsent, types.TileAbsent},
			Outcome: types.GuessOutcomeIncorrect,
		}
		guess.RefreshId()
		result, err := guess.Result.Value()
		require.NoError(t, err)
		_, err = dao.db.ExecContext(ctx, `
			INSERT INTO round_member_guesses (id, round_id, user_id, attempt, word, result, outcome)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			guess.ID, guess.RoundID, guess.UserID, guess.Attempt, guess.Word, result, guess.Outcome,
		)
		require.NoError(t, err)

		groupWhere := squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: group.ID}
		memberCount, err := dao.CountGroupMembers(ctx, NewOptions().WithWhere(groupWhere))
		require.NoError(t, err)
		require.Equal(t, 2, memberCount)

		joinCount, err := dao.CountGroupJoinRequests(ctx, NewOptions().WithWhere(squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: group.ID}))
		require.NoError(t, err)
		require.Equal(t, 2, joinCount)

		rounds, err := dao.ListRounds(ctx, NewOptions().WithWhere(squirrel.Eq{models.ROUND_GROUP_ID: group.ID}))
		require.NoError(t, err)
		require.Len(t, rounds, 1)

		roundMembers, err := dao.ListRoundMembers(ctx, NewOptions().WithWhere(squirrel.Eq{models.ROUND_MEMBER_ROUND_ID: round.ID}))
		require.NoError(t, err)
		require.Len(t, roundMembers, 1)

		var guessCount int
		require.NoError(t, dao.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM round_member_guesses WHERE round_id = ?`, round.ID).Scan(&guessCount))
		require.Equal(t, 1, guessCount)

		require.NoError(t, dao.DeleteGroups(ctx, NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: group.ID})))

		deletedGroup, err := dao.GetGroup(ctx, NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: group.ID}))
		require.NoError(t, err)
		require.Nil(t, deletedGroup)

		memberCount, err = dao.CountGroupMembers(ctx, NewOptions().WithWhere(groupWhere))
		require.NoError(t, err)
		require.Zero(t, memberCount)

		joinCount, err = dao.CountGroupJoinRequests(ctx, NewOptions().WithWhere(squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: group.ID}))
		require.NoError(t, err)
		require.Zero(t, joinCount)

		rounds, err = dao.ListRounds(ctx, NewOptions().WithWhere(squirrel.Eq{models.ROUND_GROUP_ID: group.ID}))
		require.NoError(t, err)
		require.Empty(t, rounds)

		roundMembers, err = dao.ListRoundMembers(ctx, NewOptions().WithWhere(squirrel.Eq{models.ROUND_MEMBER_ROUND_ID: round.ID}))
		require.NoError(t, err)
		require.Empty(t, roundMembers)

		require.NoError(t, dao.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM round_member_guesses WHERE round_id = ?`, round.ID).Scan(&guessCount))
		require.Zero(t, guessCount)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_UpdateGroup(t *testing.T) {
	// Test successfully updating a group name
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		group := &models.Group{Name: "Friends", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		group.Name = "Best Friends"
		require.NoError(t, dao.UpdateGroup(ctx, group))

		updated, err := dao.GetGroup(ctx, NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: group.ID}))
		require.NoError(t, err)
		require.Equal(t, "Best Friends", updated.Name)
	})

	// Test error due to duplicate group name
	t.Run("duplicate name", func(t *testing.T) {
		dao, ctx := setup(t)

		userID := testUserID(t, dao, ctx)
		require.NoError(t, dao.CreateGroup(ctx, &models.Group{Name: "Friends", CreatedBy: userID}))

		group := &models.Group{Name: "Work", CreatedBy: userID}
		require.NoError(t, dao.CreateGroup(ctx, group))

		group.Name = "Friends"
		err := dao.UpdateGroup(ctx, group)
		require.ErrorContains(t, err, "UNIQUE constraint failed")
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
