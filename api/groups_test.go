package api

import (
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
