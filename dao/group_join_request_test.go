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

func Test_CreateJoinRequest(t *testing.T) {
	// Test successfully inserting a join request
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		requester := createTestUser(t, dao, ctx, "joiner")
		group := createTestGroup(t, dao, ctx, "Join Request Group", user.ID)

		jr := &models.GroupJoinRequest{GroupID: group.ID, UserID: requester.ID}
		require.NoError(t, dao.CreateJoinRequest(ctx, jr))
		require.Equal(t, types.JoinPending, jr.Status)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_GetJoinRequest(t *testing.T) {
	// Test successfully retrieving a join request
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		requester := createTestUser(t, dao, ctx, "joiner-get")
		group := createTestGroup(t, dao, ctx, "Get Join Request", user.ID)

		jr := &models.GroupJoinRequest{GroupID: group.ID, UserID: requester.ID}
		require.NoError(t, dao.CreateJoinRequest(ctx, jr))

		record, err := dao.GetJoinRequest(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: jr.ID}))
		require.NoError(t, err)
		require.Equal(t, jr.ID, record.ID)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_ListJoinRequests(t *testing.T) {
	// Test successfully listing join requests across multiple groups
	t.Run("all", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		joinerA := createTestUser(t, dao, ctx, "joiner-a")
		joinerB := createTestUser(t, dao, ctx, "joiner-b")
		joinerC := createTestUser(t, dao, ctx, "joiner-c")
		groupA := createTestGroup(t, dao, ctx, "List A", user.ID)
		groupB := createTestGroup(t, dao, ctx, "List B", user.ID)

		require.NoError(t, dao.CreateJoinRequest(ctx, &models.GroupJoinRequest{GroupID: groupA.ID, UserID: joinerA.ID}))
		require.NoError(t, dao.CreateJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: groupA.ID,
			UserID:  joinerB.ID,
			Status:  types.JoinApproved,
		}))
		require.NoError(t, dao.CreateJoinRequest(ctx, &models.GroupJoinRequest{GroupID: groupB.ID, UserID: joinerB.ID}))
		require.NoError(t, dao.CreateJoinRequest(ctx, &models.GroupJoinRequest{GroupID: groupB.ID, UserID: joinerC.ID}))

		requests, err := dao.ListJoinRequests(ctx, NewOptions())
		require.NoError(t, err)
		require.Len(t, requests, 4)

		byGroup := map[string]int{}
		for _, jr := range requests {
			byGroup[jr.GroupID]++
		}

		require.Equal(t, 2, byGroup[groupA.ID])
		require.Equal(t, 2, byGroup[groupB.ID])
	})

	// Test successfully listing join requests for a single group
	t.Run("single group", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		joinerA := createTestUser(t, dao, ctx, "joiner-single-a")
		joinerB := createTestUser(t, dao, ctx, "joiner-single-b")
		joinerC := createTestUser(t, dao, ctx, "joiner-single-c")
		groupA := createTestGroup(t, dao, ctx, "Filter A", user.ID)
		groupB := createTestGroup(t, dao, ctx, "Filter B", user.ID)

		require.NoError(t, dao.CreateJoinRequest(ctx, &models.GroupJoinRequest{GroupID: groupA.ID, UserID: joinerA.ID}))
		require.NoError(t, dao.CreateJoinRequest(ctx, &models.GroupJoinRequest{GroupID: groupA.ID, UserID: joinerB.ID}))
		require.NoError(t, dao.CreateJoinRequest(ctx, &models.GroupJoinRequest{GroupID: groupB.ID, UserID: joinerC.ID}))

		requests, err := dao.ListJoinRequests(ctx, NewOptions().WithWhere(squirrel.Eq{
			models.JOIN_REQUEST_GROUP_ID: groupA.ID,
		}))
		require.NoError(t, err)
		require.Len(t, requests, 2)

		for _, jr := range requests {
			require.Equal(t, groupA.ID, jr.GroupID)
		}

		userIDs := []string{requests[0].UserID, requests[1].UserID}
		require.ElementsMatch(t, []string{joinerA.ID, joinerB.ID}, userIDs)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_UpdateJoinRequest(t *testing.T) {
	// Test successfully updating a join request status
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		requester := createTestUser(t, dao, ctx, "joiner-update")
		group := createTestGroup(t, dao, ctx, "Update Join Request", user.ID)

		jr := &models.GroupJoinRequest{GroupID: group.ID, UserID: requester.ID}
		require.NoError(t, dao.CreateJoinRequest(ctx, jr))

		jr.Status = types.JoinApproved
		require.NoError(t, dao.UpdateJoinRequest(ctx, jr))

		record, err := dao.GetJoinRequest(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: jr.ID}))
		require.NoError(t, err)
		require.Equal(t, types.JoinApproved, record.Status)
	})

	// Test error due to missing ID
	t.Run("missing id", func(t *testing.T) {
		dao, ctx := setup(t)

		require.ErrorIs(t, dao.UpdateJoinRequest(ctx, &models.GroupJoinRequest{}), utils.ErrId)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_DeleteJoinRequests(t *testing.T) {
	// Test successfully deleting a join request
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		requester := createTestUser(t, dao, ctx, "joiner-delete")
		group := createTestGroup(t, dao, ctx, "Delete Join Request", user.ID)

		jr := &models.GroupJoinRequest{GroupID: group.ID, UserID: requester.ID}
		require.NoError(t, dao.CreateJoinRequest(ctx, jr))

		require.NoError(t, dao.DeleteJoinRequests(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: jr.ID})))

		record, err := dao.GetJoinRequest(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: jr.ID}))
		require.NoError(t, err)
		require.Nil(t, record)
	})

	// Test error due to missing where clause
	t.Run("missing where", func(t *testing.T) {
		dao, ctx := setup(t)

		require.ErrorIs(t, dao.DeleteJoinRequests(ctx, nil), utils.ErrWhere)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func Test_ListPendingJoinGroupIDsForUser(t *testing.T) {
	// Test successfully listing pending join group IDs
	t.Run("success", func(t *testing.T) {
		dao, ctx := setup(t)
		user := principalUser(t, dao, ctx)
		requester := createTestUser(t, dao, ctx, "joiner-list")
		groupA := createTestGroup(t, dao, ctx, "Pending A", user.ID)
		groupB := createTestGroup(t, dao, ctx, "Pending B", user.ID)

		require.NoError(t, dao.CreateJoinRequest(ctx, &models.GroupJoinRequest{GroupID: groupA.ID, UserID: requester.ID}))
		require.NoError(t, dao.CreateJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: groupB.ID,
			UserID:  requester.ID,
			Status:  types.JoinApproved,
		}))

		ids, err := dao.ListPendingJoinGroupIDsForUser(ctx, requester.ID, []string{groupA.ID, groupB.ID})
		require.NoError(t, err)
		require.Equal(t, []string{groupA.ID}, ids)
	})
}
