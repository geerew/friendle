package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGetGroupRound exercises fetching the current round state
func TestGetGroupRound(t *testing.T) {
	// Test successfully returning none when no round exists
	t.Run("none", func(t *testing.T) {
		router, ctx, _ := setup(t, "user", types.SiteRoleUser)
		group := createTestGroupWithMember(t, router, ctx, "user", types.GroupRoleAdmin, "Round Group")

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/rounds/current", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(body, &resp))
		require.Equal(t, "none", resp["status"])
	})

	// Test successfully returning an awaiting-word round
	t.Run("awaiting word", func(t *testing.T) {
		router, ctx, _ := setup(t, "user", types.SiteRoleUser)
		group := createTestGroupWithMember(t, router, ctx, "user", types.GroupRoleAdmin, "Active Round")

		roundDate := time.Now().Format("2006-01-02")
		round := &models.Round{
			GroupID: group.ID, RoundDate: roundDate, PickerUserID: "user",
			Status: types.RoundAwaitingWord,
		}
		require.NoError(t, router.appDao.CreateRound(ctx, round))

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/rounds/current", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(body, &resp))
		require.Equal(t, string(types.RoundAwaitingWord), resp["status"])
		require.Equal(t, "picker", resp["yourRole"])
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestCreateGroupRoundWord exercises submitting the picker's word
func TestCreateGroupRoundWord(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)
	group := createTestGroupWithMember(t, router, ctx, "user", types.GroupRoleAdmin, "Word Group")

	roundDate := time.Now().Format("2006-01-02")
	round := &models.Round{
		GroupID: group.ID, RoundDate: roundDate, PickerUserID: "user",
		Status: types.RoundAwaitingWord,
	}
	require.NoError(t, router.appDao.CreateRound(ctx, round))

	body := bytes.NewBufferString(`{"word":"about"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/rounds/current/word", body)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	status, _, err := requestHelper(t, router, req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	updated, err := router.appDao.GetRound(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
		"group_id": group.ID,
		"round_date": roundDate,
	}))
	require.NoError(t, err)
	require.Equal(t, types.RoundActive, updated.Status)
	require.NotNil(t, updated.WordPlain)
	require.Equal(t, "ABOUT", *updated.WordPlain)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestCreateGroupRoundGuess exercises submitting a guess
func TestCreateGroupRoundGuess(t *testing.T) {
	router, ctx, principal := setup(t, "user", types.SiteRoleUser)

	picker := &models.User{Username: "picker", DisplayName: "Picker", SiteRole: types.SiteRoleUser}
	guesser := &models.User{Username: "guesser", DisplayName: "Guesser", SiteRole: types.SiteRoleUser}
	createTestUser(t, router, ctx, picker)
	createTestUser(t, router, ctx, guesser)

	group := &models.Group{Name: "Guess Group", CreatedBy: picker.ID}
	require.NoError(t, router.appDao.CreateGroup(ctx, group))
	require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
		GroupID: group.ID, UserID: picker.ID, GroupRole: types.GroupRoleAdmin,
	}))
	require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
		GroupID: group.ID, UserID: guesser.ID, GroupRole: types.GroupRoleUser,
	}))

	word := "ABOUT"
	roundDate := time.Now().Format("2006-01-02")
	round := &models.Round{
		GroupID: group.ID, RoundDate: roundDate, PickerUserID: picker.ID,
		Status: types.RoundActive, WordPlain: &word,
	}
	require.NoError(t, router.appDao.CreateRound(ctx, round))

	principal.userID = guesser.ID
	principal.role = types.SiteRoleUser

	body := bytes.NewBufferString(`{"word":"about"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/rounds/current/guesses", body)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	status, respBody, err := requestHelper(t, router, req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(respBody, &resp))
	require.Equal(t, true, resp["won"])
	require.Equal(t, true, resp["finished"])
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGetGroupRoundReveal exercises the round reveal endpoint
func TestGetGroupRoundReveal(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)
	group := createTestGroupWithMember(t, router, ctx, "user", types.GroupRoleAdmin, "Reveal Group")

	word := "ABOUT"
	roundDate := time.Now().Format("2006-01-02")
	round := &models.Round{
		GroupID: group.ID, RoundDate: roundDate, PickerUserID: "user",
		Status: types.RoundCompleted, WordPlain: &word,
	}
	require.NoError(t, router.appDao.CreateRound(ctx, round))

	status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/rounds/current/reveal", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &resp))
	require.Equal(t, string(types.RoundCompleted), resp["status"])
	require.Equal(t, "ABOUT", resp["word"])
}
