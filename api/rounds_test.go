package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/service"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_GetRoundToday exercises today's round state for group members
func TestGroups_GetRoundToday(t *testing.T) {
	// Test successfully returning today's round for a group member
	t.Run("200 (round)", func(t *testing.T) {
		router, ctx, principal := setup(t, "alice", types.SiteRoleUser)

		group := &models.Group{Name: "Friends", CreatedBy: principal.userID}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleAdmin,
		}))

		bob := &models.User{
			Base:        models.Base{ID: "bob"},
			Username:    "bob",
			DisplayName: "Bob",
			SiteRole:    types.SiteRoleUser,
		}
		createTestUser(t, router, ctx, bob)
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    bob.ID,
			GroupRole: types.GroupRoleUser,
		}))

		today := utils.DateString(time.Now())
		round, err := router.appSvc.Rounds.CreateDailyRound(ctx, group.ID, today)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/round/today", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		var resp service.RoundTodayResponse
		require.NoError(t, json.Unmarshal(body, &resp))
		require.True(t, resp.MemberThresholdMet)
		require.NotNil(t, resp.Round)
		require.Equal(t, today, resp.Round.RoundDate)
		require.Equal(t, types.RoundAwaitingWord, resp.Round.Status)
		require.Equal(t, round.PickerUserID == principal.userID, resp.Round.IsPicker)

		var raw map[string]any
		require.NoError(t, json.Unmarshal(body, &raw))
		roundBody, ok := raw["round"].(map[string]any)
		require.True(t, ok)
		_, hasPickerID := roundBody["pickerUserId"]
		_, hasPickerName := roundBody["pickerDisplayName"]
		require.False(t, hasPickerID)
		require.False(t, hasPickerName)
	})

	// Test successfully avoiding back-to-back picker assignment
	t.Run("picker rotation", func(t *testing.T) {
		router, ctx, principal := setup(t, "alice", types.SiteRoleUser)

		group := &models.Group{Name: "Friends", CreatedBy: principal.userID}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleAdmin,
		}))

		bob := &models.User{
			Base:        models.Base{ID: "bob"},
			Username:    "bob",
			DisplayName: "Bob",
			SiteRole:    types.SiteRoleUser,
		}
		createTestUser(t, router, ctx, bob)
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    bob.ID,
			GroupRole: types.GroupRoleUser,
		}))

		now := time.Now()
		yesterday := utils.PreviousDateString(now)
		today := utils.DateString(now)

		first, err := router.appSvc.Rounds.CreateDailyRound(ctx, group.ID, yesterday)
		require.NoError(t, err)

		second, err := router.appSvc.Rounds.CreateDailyRound(ctx, group.ID, today)
		require.NoError(t, err)
		require.NotEqual(t, first.PickerUserID, second.PickerUserID)
	})

	// Test error due to a non-member requesting today's round
	t.Run("403 (non-member)", func(t *testing.T) {
		router, ctx, _ := setup(t, "bob", types.SiteRoleUser)

		alice := &models.User{
			Base:     models.Base{ID: "alice"},
			Username: "alice",
			SiteRole: types.SiteRoleUser,
		}
		createTestUser(t, router, ctx, alice)

		group := &models.Group{Name: "Friends", CreatedBy: alice.ID}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    alice.ID,
			GroupRole: types.GroupRoleAdmin,
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/round/today", nil)

		status, respBody, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(respBody), "Forbidden")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestRounds_CloseStaleRounds exercises closing open rounds from before today
func TestRounds_CloseStaleRounds(t *testing.T) {
	// Test successfully completing yesterday's awaiting-word round without creating today's
	t.Run("awaiting word completed", func(t *testing.T) {
		router, ctx, principal := setup(t, "alice", types.SiteRoleUser)

		group := &models.Group{Name: "Friends", CreatedBy: principal.userID}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleAdmin,
		}))

		bob := &models.User{
			Base:        models.Base{ID: "bob"},
			Username:    "bob",
			DisplayName: "Bob",
			SiteRole:    types.SiteRoleUser,
		}
		createTestUser(t, router, ctx, bob)
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    bob.ID,
			GroupRole: types.GroupRoleUser,
		}))

		now := time.Now()
		yesterday := utils.PreviousDateString(now)
		today := utils.DateString(now)

		yesterdayRound, err := router.appSvc.Rounds.CreateDailyRound(ctx, group.ID, yesterday)
		require.NoError(t, err)

		require.NoError(t, router.appSvc.Rounds.CloseStaleRounds(ctx))

		stored, err := router.appDao.GetRound(ctx, dao.NewOptions().WithWhere(squirrel.And{
			squirrel.Eq{models.ROUND_GROUP_ID: group.ID},
			squirrel.Eq{models.ROUND_ROUND_DATE: yesterday},
		}))
		require.NoError(t, err)
		require.NotNil(t, stored)
		require.Equal(t, types.RoundCompleted, stored.Status)

		picker, err := router.appDao.GetGroupMember(ctx, dao.NewOptions().WithWhere(squirrel.And{
			squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: group.ID},
			squirrel.Eq{models.GROUP_MEMBER_USER_ID: yesterdayRound.PickerUserID},
		}))
		require.NoError(t, err)
		require.NotNil(t, picker)
		require.Equal(t, 1, picker.PickerSkips)

		todayRound, err := router.appDao.GetRound(ctx, dao.NewOptions().WithWhere(squirrel.And{
			squirrel.Eq{models.ROUND_GROUP_ID: group.ID},
			squirrel.Eq{models.ROUND_ROUND_DATE: today},
		}))
		require.NoError(t, err)
		require.Nil(t, todayRound)
	})

	// Test successfully completing yesterday's active round
	t.Run("active completed", func(t *testing.T) {
		router, ctx, principal := setup(t, "alice", types.SiteRoleUser)

		group := &models.Group{Name: "Friends", CreatedBy: principal.userID}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleAdmin,
		}))

		bob := &models.User{
			Base:        models.Base{ID: "bob"},
			Username:    "bob",
			DisplayName: "Bob",
			SiteRole:    types.SiteRoleUser,
		}
		createTestUser(t, router, ctx, bob)
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    bob.ID,
			GroupRole: types.GroupRoleUser,
		}))

		yesterday := utils.PreviousDateString(time.Now())
		round, err := router.appSvc.Rounds.CreateDailyRound(ctx, group.ID, yesterday)
		require.NoError(t, err)

		round.Status = types.RoundActive
		require.NoError(t, router.appDao.UpdateRound(ctx, round))

		require.NoError(t, router.appSvc.Rounds.CloseStaleRounds(ctx))

		stored, err := router.appDao.GetRound(ctx, dao.NewOptions().WithWhere(squirrel.And{
			squirrel.Eq{models.ROUND_GROUP_ID: group.ID},
			squirrel.Eq{models.ROUND_ROUND_DATE: yesterday},
		}))
		require.NoError(t, err)
		require.NotNil(t, stored)
		require.Equal(t, types.RoundCompleted, stored.Status)
	})
}
