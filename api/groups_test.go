package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Test successfully listing the caller's groups
func TestListMyGroups(t *testing.T) {
	router, _, _ := setup(t, "user", types.SiteRoleUser)

	for _, path := range []string{"/api/groups", "/api/groups/"} {
		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.NoError(t, err)
		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		t.Logf("path=%s status=%d body=%s", path, status, string(body))
		require.Equal(t, http.StatusOK, status)
		require.Equal(t, "[]", string(body))
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Test search results are ordered by relevance
func TestSearchGroupsOrder(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)

	for _, name := range []string{"1", "10", "11", "2"} {
		group := &models.Group{Name: name, CreatedBy: "user", IntervalHours: 24, Timezone: "UTC"}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		time.Sleep(time.Millisecond)
	}

	params := url.Values{"q": {"1"}}
	status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups/search?"+params.Encode(), nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	var result struct {
		Items []struct {
			Name string `json:"name"`
		} `json:"items"`
		TotalItems int `json:"totalItems"`
	}
	require.NoError(t, json.Unmarshal(body, &result))
	require.Equal(t, 3, result.TotalItems)
	require.Equal(t, []string{"1", "10", "11"}, []string{
		result.Items[0].Name,
		result.Items[1].Name,
		result.Items[2].Name,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Test successfully cancelling a pending join request
func TestCancelJoinRequest(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)

	group := &models.Group{Name: "Join Me", CreatedBy: "user", IntervalHours: 24, Timezone: "UTC"}
	require.NoError(t, router.appDao.CreateGroup(ctx, group))

	createReq, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests", nil)
	require.NoError(t, err)
	status, _, err := requestHelper(t, router, createReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)

	cancelReq, err := http.NewRequest(http.MethodDelete, "/api/groups/"+group.ID+"/join-requests/me", nil)
	require.NoError(t, err)
	status, _, err = requestHelper(t, router, cancelReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	cancelAgain, err := http.NewRequest(http.MethodDelete, "/api/groups/"+group.ID+"/join-requests/me", nil)
	require.NoError(t, err)
	status, body, err := requestHelper(t, router, cancelAgain)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, status)
	require.Contains(t, string(body), "Pending request not found")
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Test join request authorization rules
func TestCreateJoinRequestAuthorization(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)

	other := &models.User{
		Username:     "other",
		DisplayName:  "Other",
		PasswordHash: "password",
		SiteRole:     types.SiteRoleUser,
	}
	require.NoError(t, router.appDao.CreateUser(ctx, other))

	group := &models.Group{Name: "Shared", CreatedBy: other.ID, IntervalHours: 24, Timezone: "UTC"}
	require.NoError(t, router.appDao.CreateGroup(ctx, group))

	// Test error due to a regular user requesting join for another user
	t.Run("forbidden for other user", func(t *testing.T) {
		body := bytes.NewBufferString(`{"userId":"` + other.ID + `"}`)
		req, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests", body)
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		status, respBody, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(respBody), "Forbidden")
	})

	// Test successfully requesting join for yourself
	t.Run("self", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests", nil)
		require.NoError(t, err)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)

		jr, err := router.appDao.GetJoinRequestByUser(ctx, group.ID, "user")
		require.NoError(t, err)
		require.NotNil(t, jr)
		require.Equal(t, types.JoinPending, jr.Status)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Test site admin direct member add bypasses join requests
func TestAdminAddGroupMember(t *testing.T) {
	router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

	member := &models.User{
		Username:     "member",
		DisplayName:  "Member",
		PasswordHash: "password",
		SiteRole:     types.SiteRoleUser,
	}
	require.NoError(t, router.appDao.CreateUser(ctx, member))

	group := &models.Group{Name: "Admin Add", CreatedBy: "admin", IntervalHours: 24, Timezone: "UTC"}
	require.NoError(t, router.appDao.CreateGroup(ctx, group))

	// Test error due to site admin attempting a join request for another user
	joinBody := bytes.NewBufferString(`{"userId":"` + member.ID + `"}`)
	joinReq, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests", joinBody)
	require.NoError(t, err)
	joinReq.Header.Set("Content-Type", "application/json")

	status, _, err := requestHelper(t, router, joinReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusForbidden, status)

	// Test successfully adding a member directly as site admin
	addBody := bytes.NewBufferString(`{"userId":"` + member.ID + `"}`)
	addReq, err := http.NewRequest(http.MethodPost, "/api/admin/groups/"+group.ID+"/members", addBody)
	require.NoError(t, err)
	addReq.Header.Set("Content-Type", "application/json")

	status, respBody, err := requestHelper(t, router, addReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)

	var resp map[string]string
	require.NoError(t, json.Unmarshal(respBody, &resp))
	require.Equal(t, member.ID, resp["userId"])

	m, err := router.appDao.GetGroupMember(ctx, group.ID, member.ID)
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, types.GroupRoleUser, m.GroupRole)

	jr, err := router.appDao.GetJoinRequestByUser(ctx, group.ID, member.ID)
	require.NoError(t, err)
	require.Nil(t, jr)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestCreateGroup exercises group creation
func TestCreateGroup(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)

	body := bytes.NewBufferString(`{"name":"My Group"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/groups/", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	status, respBody, err := requestHelper(t, router, req)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(respBody, &resp))
	require.Equal(t, "My Group", resp["name"])

	m, err := router.appDao.GetGroupMember(ctx, resp["id"].(string), "user")
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, types.GroupRoleAdmin, m.GroupRole)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGetGroup exercises fetching a group the caller belongs to
func TestGetGroup(t *testing.T) {
	// Test successfully fetching a group as a member
	t.Run("member", func(t *testing.T) {
		router, ctx, _ := setup(t, "user", types.SiteRoleUser)

		group := createTestGroupWithMember(t, router, ctx, "user", types.GroupRoleAdmin, "Detail Group")

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID, nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(body, &resp))
		require.Equal(t, group.ID, resp["id"])
		require.Equal(t, "Detail Group", resp["name"])
		require.NotNil(t, resp["round"])
	})

	// Test successfully fetching a group as site admin without membership
	t.Run("site admin", func(t *testing.T) {
		router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

		member := &models.User{Username: "member", DisplayName: "Member", SiteRole: types.SiteRoleUser}
		createTestUser(t, router, ctx, member)

		group := createTestGroupWithMember(t, router, ctx, member.ID, types.GroupRoleAdmin, "Admin View Group")

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID, nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(body, &resp))
		require.Equal(t, group.ID, resp["id"])
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestUpdateGroup exercises group updates by a group admin
func TestUpdateGroup(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)

	group := createTestGroupWithMember(t, router, ctx, "user", types.GroupRoleAdmin, "Old Name")

	body := bytes.NewBufferString(`{"name":"New Name"}`)
	req, err := http.NewRequest(http.MethodPatch, "/api/groups/"+group.ID, body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	status, respBody, err := requestHelper(t, router, req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(respBody, &resp))
	require.Equal(t, "New Name", resp["name"])

	updated, err := router.appDao.GetGroup(ctx, group.ID)
	require.NoError(t, err)
	require.Equal(t, "New Name", updated.Name)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestUpdateGroupJoinRequestApprove exercises approving a join request
func TestUpdateGroupJoinRequestApprove(t *testing.T) {
	router, ctx, principal := setup(t, "user", types.SiteRoleUser)

	gadmin := &models.User{Username: "gadmin", DisplayName: "Group Admin", SiteRole: types.SiteRoleUser}
	createTestUser(t, router, ctx, gadmin)

	group := createTestGroupWithMember(t, router, ctx, gadmin.ID, types.GroupRoleAdmin, "Approve Group")

	joinReq, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests", nil)
	require.NoError(t, err)
	status, _, err := requestHelper(t, router, joinReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)

	jr, err := router.appDao.GetJoinRequestByUser(ctx, group.ID, "user")
	require.NoError(t, err)
	require.NotNil(t, jr)

	principal.userID = gadmin.ID
	principal.role = types.SiteRoleUser

	approveReq, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests/"+jr.ID+"/approve", nil)
	require.NoError(t, err)
	status, _, err = requestHelper(t, router, approveReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	m, err := router.appDao.GetGroupMember(ctx, group.ID, "user")
	require.NoError(t, err)
	require.NotNil(t, m)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestUpdateGroupJoinRequestReject exercises rejecting a join request
func TestUpdateGroupJoinRequestReject(t *testing.T) {
	router, ctx, principal := setup(t, "user", types.SiteRoleUser)

	gadmin := &models.User{Username: "greject", DisplayName: "Group Admin", SiteRole: types.SiteRoleUser}
	createTestUser(t, router, ctx, gadmin)

	group := createTestGroupWithMember(t, router, ctx, gadmin.ID, types.GroupRoleAdmin, "Reject Group")

	joinReq, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests", nil)
	require.NoError(t, err)
	status, _, err := requestHelper(t, router, joinReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)

	jr, err := router.appDao.GetJoinRequestByUser(ctx, group.ID, "user")
	require.NoError(t, err)
	require.NotNil(t, jr)

	principal.userID = gadmin.ID
	principal.role = types.SiteRoleUser

	rejectReq, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests/"+jr.ID+"/reject", nil)
	require.NoError(t, err)
	status, _, err = requestHelper(t, router, rejectReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	m, err := router.appDao.GetGroupMember(ctx, group.ID, "user")
	require.NoError(t, err)
	require.Nil(t, m)

	updated, err := router.appDao.GetJoinRequest(ctx, jr.ID)
	require.NoError(t, err)
	require.Equal(t, types.JoinRejected, updated.Status)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestDeleteGroupMember exercises removing a member from a group
func TestDeleteGroupMember(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)

	group := createTestGroupWithMember(t, router, ctx, "user", types.GroupRoleAdmin, "Members Group")

	member := &models.User{Username: "member", DisplayName: "Member", SiteRole: types.SiteRoleUser}
	createTestUser(t, router, ctx, member)
	require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
		GroupID: group.ID, UserID: member.ID, GroupRole: types.GroupRoleUser,
	}))

	delReq, err := http.NewRequest(http.MethodDelete, "/api/groups/"+group.ID+"/members/"+member.ID, nil)
	require.NoError(t, err)
	status, _, err := requestHelper(t, router, delReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	removed, err := router.appDao.GetGroupMember(ctx, group.ID, member.ID)
	require.NoError(t, err)
	require.Nil(t, removed)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGetGroupLeaderboard exercises the group leaderboard endpoint
func TestGetGroupLeaderboard(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)

	group := createTestGroupWithMember(t, router, ctx, "user", types.GroupRoleAdmin, "Leaderboard Group")

	status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/leaderboard", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	var rows []map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &rows))
	require.Len(t, rows, 1)
	require.Equal(t, "user", rows[0]["userId"])
}
