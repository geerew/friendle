package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
