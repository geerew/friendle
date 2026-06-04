package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/types"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestListUsers exercises site-admin user listing
func TestListUsers(t *testing.T) {
	// Test successfully listing users for a site admin
	t.Run("success", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/admin/users", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		result, users := unmarshalHelper[map[string]any](t, body)
		require.GreaterOrEqual(t, result.TotalItems, 1)
		require.Len(t, users, result.TotalItems)
		require.Equal(t, pagination.DefaultPerPage, result.PerPage)
	})

	// Test error due to a non-admin caller
	t.Run("forbidden", func(t *testing.T) {
		router, _, _ := setup(t, "user", types.SiteRoleUser)

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/admin/users", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(body), "Forbidden")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestCreateUser exercises site-admin user creation
func TestCreateUser(t *testing.T) {
	// Test successfully creating a user
	t.Run("created", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/users", strings.NewReader(`{"username": "testuser", "password": "password123"}`))
		req.Header.Set("Content-Type", "application/json")

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)
	})

	// Test error due to a non-admin caller
	t.Run("forbidden", func(t *testing.T) {
		router, _, _ := setup(t, "user", types.SiteRoleUser)

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodPost, "/api/admin/users", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(body), "Forbidden")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestUpdateUser exercises site-admin user updates
func TestUpdateUser(t *testing.T) {
	router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

	user := &models.User{
		Username:     "test",
		DisplayName:  "Test",
		PasswordHash: "test-password-hash",
		SiteRole:     types.SiteRoleUser,
	}
	require.NoError(t, router.appDao.CreateUser(ctx, user))

	req := httptest.NewRequest(http.MethodPut, "/api/admin/users/"+user.ID, strings.NewReader(`{"displayName": "Bob"}`))
	req.Header.Set("Content-Type", "application/json")

	status, _, err := requestHelper(t, router, req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	record, err := router.appDao.GetUser(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: user.ID}))
	require.NoError(t, err)
	require.Equal(t, "Bob", record.DisplayName)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestDeleteUser exercises site-admin user deletion
func TestDeleteUser(t *testing.T) {
	router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

	target := &models.User{Username: "delete-me", DisplayName: "Delete Me", SiteRole: types.SiteRoleUser}
	createTestUser(t, router, ctx, target)

	status, _, err := requestHelper(t, router, httptest.NewRequest(http.MethodDelete, "/api/admin/users/"+target.ID, nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	deleted, err := router.appDao.GetUser(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: target.ID}))
	require.NoError(t, err)
	require.Nil(t, deleted)
}
