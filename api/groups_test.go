package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/types"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Test successfully listing the caller's groups
func TestListMyGroups(t *testing.T) {
	router, _, _ := setup(t, "user", types.SiteRoleUser)

	for _, path := range []string{"/api/groups/mine", "/api/groups/mine/"} {
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

// TestListGroups exercises the site-wide group list
func TestListGroups(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)

	group := &models.Group{Name: "List Test Group", CreatedBy: "user"}
	require.NoError(t, router.appDao.CreateGroup(ctx, group))

	// Test successfully listing groups without site admin access
	t.Run("success", func(t *testing.T) {
		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		result, groups := unmarshalHelper[map[string]any](t, body)
		require.Equal(t, 1, result.TotalItems)
		require.Len(t, groups, 1)
	})

	// Test successfully paginating groups
	t.Run("pagination", func(t *testing.T) {
		for i := range 4 {
			group := &models.Group{Name: "Paginated Group " + string(rune('A'+i)), CreatedBy: "user"}
			require.NoError(t, router.appDao.CreateGroup(ctx, group))
		}

		params := url.Values{
			pagination.PageQueryParam:    {"1"},
			pagination.PerPageQueryParam: {"2"},
		}
		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups?"+params.Encode(), nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		result, groups := unmarshalHelper[map[string]any](t, body)
		require.Equal(t, 5, result.TotalItems)
		require.Len(t, groups, 2)
	})

	// Test error due to an unauthenticated caller
	t.Run("401", func(t *testing.T) {
		router, _, _ := setup(t, "", types.SiteRoleUser)
		router.app.SetBootstrapped()

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, status)
		require.Contains(t, string(body), "Unauthorized")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestDeleteGroup exercises group deletion authorization
func TestDeleteGroup(t *testing.T) {
	// Test successfully deleting a group as a site admin
	t.Run("site admin", func(t *testing.T) {
		router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

		group := &models.Group{Name: "Delete Group", CreatedBy: "admin"}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))

		status, _, err := requestHelper(t, router, httptest.NewRequest(http.MethodDelete, "/api/groups/"+group.ID, nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, status)

		deleted, err := router.appDao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: group.ID}))
		require.NoError(t, err)
		require.Nil(t, deleted)
	})

	// Test successfully deleting a group as a group admin
	t.Run("group admin", func(t *testing.T) {
		router, ctx, _ := setup(t, "user", types.SiteRoleUser)

		group := createTestGroupWithMember(t, router, ctx, "user", types.GroupRoleAdmin, "Group Admin Delete")

		status, _, err := requestHelper(t, router, httptest.NewRequest(http.MethodDelete, "/api/groups/"+group.ID, nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, status)

		deleted, err := router.appDao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: group.ID}))
		require.NoError(t, err)
		require.Nil(t, deleted)
	})

	// Test error due to a non-admin group member
	t.Run("403 member", func(t *testing.T) {
		router, ctx, _ := setup(t, "user", types.SiteRoleUser)

		group := createTestGroupWithMember(t, router, ctx, "user", types.GroupRoleUser, "Member Delete")

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodDelete, "/api/groups/"+group.ID, nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(body), "Forbidden")

		existing, err := router.appDao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: group.ID}))
		require.NoError(t, err)
		require.NotNil(t, existing)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Test search results are ordered by relevance
func TestSearchGroupsOrder(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)

	for _, name := range []string{"1", "10", "11", "2"} {
		group := &models.Group{Name: name, CreatedBy: "user"}
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

	group := &models.Group{Name: "Join Me", CreatedBy: "user"}
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

	group := &models.Group{Name: "Shared", CreatedBy: other.ID}
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

		jr, err := router.appDao.GetJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
			"group_id": group.ID,
			"user_id":  "user",
		}))
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

	group := &models.Group{Name: "Admin Add", CreatedBy: "admin"}
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

	m, err := router.appDao.GetGroupMember(ctx, dao.NewOptions().WithWhere(squirrel.Eq{"group_id": group.ID, "user_id": member.ID}))
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, types.GroupRoleUser, m.GroupRole)

	jr, err := router.appDao.GetJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
			"group_id": group.ID,
			"user_id":  member.ID,
		}))
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

	m, err := router.appDao.GetGroupMember(ctx, dao.NewOptions().WithWhere(squirrel.Eq{"group_id": resp["id"].(string), "user_id": "user"}))
	require.NoError(t, err)
	require.NotNil(t, m)
	require.Equal(t, types.GroupRoleAdmin, m.GroupRole)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestCreateGroupDuplicateName exercises case-insensitive group name uniqueness
func TestCreateGroupDuplicateName(t *testing.T) {
	router, ctx, _ := setup(t, "user", types.SiteRoleUser)

	group := &models.Group{Name: "Unique Group", CreatedBy: "user"}
	require.NoError(t, router.appDao.CreateGroup(ctx, group))

	body := bytes.NewBufferString(`{"name":"unique group"}`)
	req, err := http.NewRequest(http.MethodPost, "/api/groups/", body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	status, respBody, err := requestHelper(t, router, req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, status)
	require.Contains(t, string(respBody), "Group name already exists")
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
		require.NotNil(t, resp["members"])
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

	updated, err := router.appDao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: group.ID}))
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

	jr, err := router.appDao.GetJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
			"group_id": group.ID,
			"user_id":  "user",
		}))
	require.NoError(t, err)
	require.NotNil(t, jr)

	principal.userID = gadmin.ID
	principal.role = types.SiteRoleUser

	approveReq, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests/"+jr.ID+"/approve", nil)
	require.NoError(t, err)
	status, _, err = requestHelper(t, router, approveReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	m, err := router.appDao.GetGroupMember(ctx, dao.NewOptions().WithWhere(squirrel.Eq{"group_id": group.ID, "user_id": "user"}))
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

	jr, err := router.appDao.GetJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
			"group_id": group.ID,
			"user_id":  "user",
		}))
	require.NoError(t, err)
	require.NotNil(t, jr)

	principal.userID = gadmin.ID
	principal.role = types.SiteRoleUser

	rejectReq, err := http.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join-requests/"+jr.ID+"/reject", nil)
	require.NoError(t, err)
	status, _, err = requestHelper(t, router, rejectReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	m, err := router.appDao.GetGroupMember(ctx, dao.NewOptions().WithWhere(squirrel.Eq{"group_id": group.ID, "user_id": "user"}))
	require.NoError(t, err)
	require.Nil(t, m)

	updated, err := router.appDao.GetJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: jr.ID}))
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

	removed, err := router.appDao.GetGroupMember(ctx, dao.NewOptions().WithWhere(squirrel.Eq{"group_id": group.ID, "user_id": member.ID}))
	require.NoError(t, err)
	require.Nil(t, removed)
}
