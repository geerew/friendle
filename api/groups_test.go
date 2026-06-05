package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
		require.NotEmpty(t, groups[0].ID)
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
