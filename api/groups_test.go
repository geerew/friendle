package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/service"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_Create exercises group creation
func TestGroups_Create(t *testing.T) {
	// Test successfully creating a group
	t.Run("201 (created)", func(t *testing.T) {
		router, _, _ := setup(t, "alice", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPost, "/api/groups/", strings.NewReader(`{"name":"Friends"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)

		var resp service.GroupResponse
		require.NoError(t, json.Unmarshal(body, &resp))
		require.Equal(t, "Friends", resp.Name)
		require.Equal(t, "alice", resp.CreatedBy)
		require.Equal(t, 1, resp.MemberCount)
		require.NotEmpty(t, resp.ID)
	})

	// Test error due to invalid JSON
	t.Run("400 (bind error)", func(t *testing.T) {
		router, _, _ := setup(t, "alice", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPost, "/api/groups/", strings.NewReader(`{`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Error parsing data")
	})

	// Test error due to missing name
	t.Run("400 (missing name)", func(t *testing.T) {
		router, _, _ := setup(t, "alice", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPost, "/api/groups/", strings.NewReader(`{}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Group name is required")
	})

	// Test error due to duplicate name
	t.Run("400 (duplicate name)", func(t *testing.T) {
		router, _, _ := setup(t, "alice", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPost, "/api/groups/", strings.NewReader(`{"name":"Friends"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)

		req = httptest.NewRequest(http.MethodPost, "/api/groups/", strings.NewReader(`{"name":"Friends"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Group name already exists")
	})

	// Test successfully listing the caller's groups
	t.Run("200 (list self)", func(t *testing.T) {
		router, _, _ := setup(t, "alice", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPost, "/api/groups/", strings.NewReader(`{"name":"Friends"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)

		req = httptest.NewRequest(http.MethodGet, "/api/groups/self", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		_, groups := unmarshalHelper[service.GroupResponse](t, body)
		require.Len(t, groups, 1)
		require.Equal(t, "Friends", groups[0].Name)
		require.Equal(t, 1, groups[0].MemberCount)
		require.NotEmpty(t, groups[0].ID)
	})

	// Test successfully listing all groups with pagination
	t.Run("200 (list)", func(t *testing.T) {
		router, _, _ := setup(t, "alice", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPost, "/api/groups/", strings.NewReader(`{"name":"Friends"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)

		req = httptest.NewRequest(http.MethodGet, "/api/groups/", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		pResult, groups := unmarshalHelper[service.GroupResponse](t, body)
		require.Equal(t, 1, pResult.TotalItems)
		require.Len(t, groups, 1)
		require.Equal(t, "Friends", groups[0].Name)
		require.Equal(t, 1, groups[0].MemberCount)
	})

	// Test list all includes groups the caller does not belong to
	t.Run("200 (list all non-member)", func(t *testing.T) {
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

		req := httptest.NewRequest(http.MethodGet, "/api/groups/self", nil)
		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		_, selfGroups := unmarshalHelper[models.Group](t, body)
		require.Empty(t, selfGroups)

		req = httptest.NewRequest(http.MethodGet, "/api/groups/", nil)
		status, body, err = requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		pResult, groups := unmarshalHelper[service.GroupResponse](t, body)
		require.Equal(t, 1, pResult.TotalItems)
		require.Len(t, groups, 1)
		require.Equal(t, "Friends", groups[0].Name)
	})

	// Test error due to missing authentication
	t.Run("401 (unauthorized)", func(t *testing.T) {
		router, _, _ := setup(t, "", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPost, "/api/groups/", strings.NewReader(`{"name":"Friends"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, status)
		require.Contains(t, string(body), "Unauthorized")

		req = httptest.NewRequest(http.MethodGet, "/api/groups/self", nil)

		status, body, err = requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, status)
		require.Contains(t, string(body), "Unauthorized")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_Get exercises fetching a single group
func TestGroups_Get(t *testing.T) {
	// Test successfully fetching a group the caller belongs to
	t.Run("200 (found)", func(t *testing.T) {
		router, _, _ := setup(t, "alice", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPost, "/api/groups/", strings.NewReader(`{"name":"Friends"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)

		var created service.GroupResponse
		require.NoError(t, json.Unmarshal(body, &created))

		req = httptest.NewRequest(http.MethodGet, "/api/groups/"+created.ID, nil)

		status, body, err = requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		var resp service.GroupResponse
		require.NoError(t, json.Unmarshal(body, &resp))
		require.Equal(t, created.ID, resp.ID)
		require.Equal(t, "Friends", resp.Name)
		require.Equal(t, 1, resp.MemberCount)
		require.Equal(t, "alice", resp.CreatedBy)
		require.Equal(t, types.GroupRoleAdmin, *resp.GroupRole)
		require.NotNil(t, resp.AdminSummary)
		require.Equal(t, 0, resp.AdminSummary.PendingJoinRequestCount)
		require.Equal(t, 0, resp.AdminSummary.RejectedJoinRequestCount)
	})

	// Test successfully fetching a group as a non-admin member
	t.Run("200 (member)", func(t *testing.T) {
		router, ctx, principal := setup(t, "bob", types.SiteRoleUser)

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
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleUser,
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID, nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		var resp service.GroupResponse
		require.NoError(t, json.Unmarshal(body, &resp))
		require.Equal(t, types.GroupRoleUser, *resp.GroupRole)
		require.Equal(t, 2, resp.MemberCount)
		require.Nil(t, resp.AdminSummary)
	})

	// Test successfully returning join request counts for a group admin
	t.Run("200 (admin summary)", func(t *testing.T) {
		router, ctx, principal := setup(t, "alice", types.SiteRoleUser)

		group := &models.Group{Name: "Friends", CreatedBy: principal.userID}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleAdmin,
		}))

		bob := &models.User{
			Base:     models.Base{ID: "bob"},
			Username: "bob",
			SiteRole: types.SiteRoleUser,
		}
		createTestUser(t, router, ctx, bob)

		carol := &models.User{
			Base:     models.Base{ID: "carol"},
			Username: "carol",
			SiteRole: types.SiteRoleUser,
		}
		createTestUser(t, router, ctx, carol)

		require.NoError(t, router.appDao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: group.ID,
			UserID:  bob.ID,
			Status:  types.JoinPending,
		}))
		require.NoError(t, router.appDao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: group.ID,
			UserID:  carol.ID,
			Status:  types.JoinRejected,
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID, nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		var resp service.GroupResponse
		require.NoError(t, json.Unmarshal(body, &resp))
		require.NotNil(t, resp.AdminSummary)
		require.Equal(t, 1, resp.AdminSummary.PendingJoinRequestCount)
		require.Equal(t, 1, resp.AdminSummary.RejectedJoinRequestCount)
	})

	// Test successfully fetching a group the caller does not belong to
	t.Run("200 (non-member)", func(t *testing.T) {
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

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID, nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		var resp service.GroupResponse
		require.NoError(t, json.Unmarshal(body, &resp))
		require.Equal(t, group.ID, resp.ID)
		require.Equal(t, "Friends", resp.Name)
		require.Nil(t, resp.AdminSummary)
	})

	// Test error due to a non-existent group
	t.Run("404 (not found)", func(t *testing.T) {
		router, _, _ := setup(t, "alice", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodGet, "/api/groups/missing-group-id", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, status)
		require.Contains(t, string(body), "Group not found")
	})

	// Test error due to missing authentication
	t.Run("401 (unauthorized)", func(t *testing.T) {
		router, _, _ := setup(t, "", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodGet, "/api/groups/some-id", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, status)
		require.Contains(t, string(body), "Unauthorized")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_Delete exercises site-admin group deletion
func TestGroups_Delete(t *testing.T) {
	// Test successfully deleting a group
	t.Run("204 (deleted)", func(t *testing.T) {
		router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

		group := &models.Group{Name: "Delete Me", CreatedBy: "admin"}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))

		status, _, err := requestHelper(t, router, httptest.NewRequest(http.MethodDelete, "/api/groups/"+group.ID, nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, status)

		deleted, err := router.appDao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: group.ID}))
		require.NoError(t, err)
		require.Nil(t, deleted)
	})

	// Test error due to a non-existent group
	t.Run("404 (not found)", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodDelete, "/api/groups/missing-group-id", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, status)
		require.Contains(t, string(body), "Group not found")
	})

	// Test error due to non-admin caller
	t.Run("403 (forbidden)", func(t *testing.T) {
		router, ctx, _ := setup(t, "alice", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPost, "/api/groups/", strings.NewReader(`{"name":"Friends"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)

		var created service.GroupResponse
		require.NoError(t, json.Unmarshal(body, &created))

		status, body, err = requestHelper(t, router, httptest.NewRequest(http.MethodDelete, "/api/groups/"+created.ID, nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(body), "Forbidden")

		stillThere, err := router.appDao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: created.ID}))
		require.NoError(t, err)
		require.NotNil(t, stillThere)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_Search exercises group name search
func TestGroups_Search(t *testing.T) {
	// Test successfully searching groups by name
	t.Run("200 (found)", func(t *testing.T) {
		router, ctx, _ := setup(t, "alice", types.SiteRoleUser)

		group := &models.Group{Name: "Friends", CreatedBy: "alice"}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))

		other := &models.Group{Name: "Work", CreatedBy: "alice"}
		require.NoError(t, router.appDao.CreateGroup(ctx, other))

		req := httptest.NewRequest(http.MethodGet, "/api/groups/?name=friend", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		pResult, groups := unmarshalHelper[service.GroupResponse](t, body)
		require.Equal(t, 1, pResult.TotalItems)
		require.Len(t, groups, 1)
		require.Equal(t, "Friends", groups[0].Name)
		require.Nil(t, groups[0].GroupRole)
		require.Nil(t, groups[0].JoinRequestStatus)
	})

	// Test prefix matches are ranked before substring matches
	t.Run("200 (prefix order)", func(t *testing.T) {
		router, ctx, _ := setup(t, "alice", types.SiteRoleUser)

		for _, name := range []string{"lightening", "ten", "tento", "toten"} {
			require.NoError(t, router.appDao.CreateGroup(ctx, &models.Group{Name: name, CreatedBy: "alice"}))
		}

		req := httptest.NewRequest(http.MethodGet, "/api/groups/?name=ten", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		pResult, groups := unmarshalHelper[service.GroupResponse](t, body)
		require.Equal(t, 4, pResult.TotalItems)
		require.Len(t, groups, 4)
		require.Equal(t, "ten", groups[0].Name)
		require.Equal(t, "tento", groups[1].Name)
		require.Equal(t, "lightening", groups[2].Name)
		require.Equal(t, "toten", groups[3].Name)
	})

	// Test successfully returning viewer status for each search result
	t.Run("200 (member status)", func(t *testing.T) {
		router, ctx, principal := setup(t, "alice", types.SiteRoleUser)

		memberGroup := &models.Group{Name: "Member Group", CreatedBy: "alice"}
		require.NoError(t, router.appDao.CreateGroup(ctx, memberGroup))
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   memberGroup.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleUser,
		}))

		adminGroup := &models.Group{Name: "Admin Group", CreatedBy: "alice"}
		require.NoError(t, router.appDao.CreateGroup(ctx, adminGroup))
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   adminGroup.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleAdmin,
		}))

		requestedGroup := &models.Group{Name: "Requested Group", CreatedBy: "alice"}
		require.NoError(t, router.appDao.CreateGroup(ctx, requestedGroup))
		require.NoError(t, router.appDao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: requestedGroup.ID,
			UserID:  principal.userID,
			Status:  types.JoinPending,
		}))

		rejectedGroup := &models.Group{Name: "Rejected Group", CreatedBy: "alice"}
		require.NoError(t, router.appDao.CreateGroup(ctx, rejectedGroup))
		require.NoError(t, router.appDao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: rejectedGroup.ID,
			UserID:  principal.userID,
			Status:  types.JoinRejected,
		}))

		openGroup := &models.Group{Name: "Open Group", CreatedBy: "alice"}
		require.NoError(t, router.appDao.CreateGroup(ctx, openGroup))

		req := httptest.NewRequest(http.MethodGet, "/api/groups/?name=group", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		pResult, groups := unmarshalHelper[service.GroupResponse](t, body)
		require.Equal(t, 5, pResult.TotalItems)
		require.Len(t, groups, 5)

		byName := make(map[string]service.GroupResponse, len(groups))
		for _, group := range groups {
			byName[group.Name] = group
		}

		require.Equal(t, types.GroupRoleAdmin, *byName["Admin Group"].GroupRole)
		require.Equal(t, types.GroupRoleUser, *byName["Member Group"].GroupRole)
		require.Equal(t, types.JoinPending, *byName["Requested Group"].JoinRequestStatus)
		require.Equal(t, types.JoinRejected, *byName["Rejected Group"].JoinRequestStatus)
		require.Nil(t, byName["Open Group"].GroupRole)
		require.Nil(t, byName["Open Group"].JoinRequestStatus)
	})

	// Test error due to missing search name
	t.Run("400 (missing name)", func(t *testing.T) {
		router, _, _ := setup(t, "alice", types.SiteRoleUser)

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/groups/?name=", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Name is required")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_RequestJoin exercises creating a group join request
func TestGroups_RequestJoin(t *testing.T) {
	// Test successfully creating a join request
	t.Run("201 (created)", func(t *testing.T) {
		router, ctx, principal := setup(t, "bob", types.SiteRoleUser)

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

		req := httptest.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)

		var resp service.GroupResponse
		require.NoError(t, json.Unmarshal(body, &resp))
		require.Equal(t, group.ID, resp.ID)
		require.Nil(t, resp.GroupRole)
		require.Equal(t, types.JoinPending, *resp.JoinRequestStatus)

		stored, err := router.appDao.GetGroupJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.And{
			squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: group.ID},
			squirrel.Eq{models.JOIN_REQUEST_USER_ID: principal.userID},
		}))
		require.NoError(t, err)
		require.NotNil(t, stored)
		require.Equal(t, types.JoinPending, stored.Status)
	})

	// Test error due to an existing rejected join request
	t.Run("400 (rejected)", func(t *testing.T) {
		router, ctx, principal := setup(t, "bob", types.SiteRoleUser)

		group := &models.Group{Name: "Closed", CreatedBy: principal.userID}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		require.NoError(t, router.appDao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: group.ID,
			UserID:  principal.userID,
			Status:  types.JoinRejected,
		}))

		req := httptest.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Join request was rejected")
	})

	// Test error due to an existing pending join request
	t.Run("400 (pending)", func(t *testing.T) {
		router, ctx, principal := setup(t, "bob", types.SiteRoleUser)

		group := &models.Group{Name: "Pending", CreatedBy: principal.userID}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		require.NoError(t, router.appDao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: group.ID,
			UserID:  principal.userID,
			Status:  types.JoinPending,
		}))

		req := httptest.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Join request already pending")
	})

	// Test error due to the caller already being a group member
	t.Run("400 (member)", func(t *testing.T) {
		router, ctx, principal := setup(t, "bob", types.SiteRoleUser)

		group := &models.Group{Name: "Member", CreatedBy: principal.userID}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleUser,
		}))

		req := httptest.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/join", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Already a group member")
	})

	// Test error due to a missing group
	t.Run("404 (not found)", func(t *testing.T) {
		router, _, _ := setup(t, "bob", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPost, "/api/groups/missing-group-id/join", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, status)
		require.Contains(t, string(body), "Group not found")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_ListMembers exercises listing group members
func TestGroups_ListMembers(t *testing.T) {
	// Test successfully listing members as a group admin
	t.Run("200 (admin)", func(t *testing.T) {
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

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/members", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		respData, members := unmarshalHelper[service.GroupMemberResponse](t, body)
		require.Equal(t, 2, respData.TotalItems)
		require.Len(t, members, 2)
		require.Equal(t, "Test User", members[0].DisplayName)
		require.Equal(t, types.GroupRoleAdmin, members[0].GroupRole)
		require.Equal(t, "Bob", members[1].DisplayName)
		require.Equal(t, types.GroupRoleUser, members[1].GroupRole)
	})

	// Test successfully listing members as a group user
	t.Run("200 (member)", func(t *testing.T) {
		router, ctx, principal := setup(t, "bob", types.SiteRoleUser)

		group := &models.Group{Name: "Friends", CreatedBy: principal.userID}
		require.NoError(t, router.appDao.CreateGroup(ctx, group))
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleUser,
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/members", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		respData, members := unmarshalHelper[service.GroupMemberResponse](t, body)
		require.Equal(t, 1, respData.TotalItems)
		require.Len(t, members, 1)
	})

	// Test successfully paginating group members
	t.Run("200 (pagination)", func(t *testing.T) {
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

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/members?page=1&perPage=1", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		respData, members := unmarshalHelper[service.GroupMemberResponse](t, body)
		require.Equal(t, 1, respData.Page)
		require.Equal(t, 1, respData.PerPage)
		require.Equal(t, 2, respData.TotalItems)
		require.Equal(t, 2, respData.TotalPages)
		require.Len(t, members, 1)
		require.Equal(t, "Test User", members[0].DisplayName)
		require.Equal(t, types.GroupRoleAdmin, members[0].GroupRole)
	})

	// Test error due to a non-member caller
	t.Run("403 (forbidden)", func(t *testing.T) {
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

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/members", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(body), "Forbidden")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_ListPendingJoinRequests exercises listing pending join requests
func TestGroups_ListPendingJoinRequests(t *testing.T) {
	// Test successfully listing pending join requests as a group admin
	t.Run("200 (admin)", func(t *testing.T) {
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
		require.NoError(t, router.appDao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: group.ID,
			UserID:  bob.ID,
			Status:  types.JoinPending,
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/pending", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		respData, requests := unmarshalHelper[service.GroupJoinRequestResponse](t, body)
		require.Equal(t, 1, respData.TotalItems)
		require.Len(t, requests, 1)
		require.Equal(t, "Bob", requests[0].DisplayName)
		require.Equal(t, bob.ID, requests[0].UserID)
	})

	// Test error due to a non-admin group member
	t.Run("403 (member)", func(t *testing.T) {
		router, ctx, principal := setup(t, "bob", types.SiteRoleUser)

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
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleUser,
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/pending", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(body), "Forbidden")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_ListRejectedJoinRequests exercises listing rejected join requests
func TestGroups_ListRejectedJoinRequests(t *testing.T) {
	// Test successfully listing rejected join requests as a group admin
	t.Run("200 (admin)", func(t *testing.T) {
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
		require.NoError(t, router.appDao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: group.ID,
			UserID:  bob.ID,
			Status:  types.JoinRejected,
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/rejected", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		respData, requests := unmarshalHelper[service.GroupJoinRequestResponse](t, body)
		require.Equal(t, 1, respData.TotalItems)
		require.Len(t, requests, 1)
		require.Equal(t, "Bob", requests[0].DisplayName)
		require.Equal(t, bob.ID, requests[0].UserID)
	})

	// Test error due to a non-admin group member
	t.Run("403 (member)", func(t *testing.T) {
		router, ctx, principal := setup(t, "bob", types.SiteRoleUser)

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
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleUser,
		}))

		req := httptest.NewRequest(http.MethodGet, "/api/groups/"+group.ID+"/rejected", nil)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(body), "Forbidden")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_ApproveJoinRequest exercises approving pending join requests
func TestGroups_ApproveJoinRequest(t *testing.T) {
	// Test successfully approving a pending join request
	t.Run("204 (approve)", func(t *testing.T) {
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
		require.NoError(t, router.appDao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: group.ID,
			UserID:  bob.ID,
			Status:  types.JoinPending,
		}))

		req := httptest.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/pending/"+bob.ID+"/approve", nil)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, status)

		member, err := router.appDao.ListGroupMembers(ctx, dao.NewOptions().WithWhere(squirrel.And{
			squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: group.ID},
			squirrel.Eq{models.GROUP_MEMBER_USER_ID: bob.ID},
		}))
		require.NoError(t, err)
		require.Len(t, member, 1)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_DeclineJoinRequest exercises declining pending join requests
func TestGroups_DeclineJoinRequest(t *testing.T) {
	// Test successfully declining a pending join request
	t.Run("204 (decline)", func(t *testing.T) {
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
		require.NoError(t, router.appDao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
			GroupID: group.ID,
			UserID:  bob.ID,
			Status:  types.JoinPending,
		}))

		req := httptest.NewRequest(http.MethodPost, "/api/groups/"+group.ID+"/pending/"+bob.ID+"/decline", nil)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, status)

		stored, err := router.appDao.GetGroupJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.And{
			squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: group.ID},
			squirrel.Eq{models.JOIN_REQUEST_USER_ID: bob.ID},
		}))
		require.NoError(t, err)
		require.NotNil(t, stored)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestGroups_UpdateMemberRole exercises updating group member roles
func TestGroups_UpdateMemberRole(t *testing.T) {
	// Test successfully promoting a member to group admin
	t.Run("200 (promote)", func(t *testing.T) {
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

		body := []byte(`{"groupRole":"group_admin"}`)
		req := httptest.NewRequest(http.MethodPatch, "/api/groups/"+group.ID+"/members/"+bob.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		status, respBody, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		member := service.GroupMemberResponse{}
		require.NoError(t, json.Unmarshal(respBody, &member))
		require.Equal(t, bob.ID, member.UserID)
		require.Equal(t, types.GroupRoleAdmin, member.GroupRole)
	})

	// Test error due to a non-admin group member
	t.Run("403 (member)", func(t *testing.T) {
		router, ctx, principal := setup(t, "bob", types.SiteRoleUser)

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
		require.NoError(t, router.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.userID,
			GroupRole: types.GroupRoleUser,
		}))

		body := []byte(`{"groupRole":"group_admin"}`)
		req := httptest.NewRequest(http.MethodPatch, "/api/groups/"+group.ID+"/members/"+alice.ID, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		status, respBody, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(respBody), "Forbidden")
	})
}
