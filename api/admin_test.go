package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/types"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAdminListUsers exercises admin user listing
func TestAdminListUsers(t *testing.T) {
	// Test successfully listing users for a site admin
	t.Run("success", func(t *testing.T) {
		router, _ := setupAdmin(t)

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/admin/users", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		result, users := unmarshalHelper[map[string]any](t, body)
		require.GreaterOrEqual(t, result.TotalItems, 1)
		require.Len(t, users, result.TotalItems)
		require.Equal(t, pagination.DefaultPerPage, result.PerPage)
	})

	// Test error due to a non-admin caller
	t.Run("403", func(t *testing.T) {
		router, _ := setupUser(t)

		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/admin/users", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(body), "Site admin required")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAdminListGroups exercises admin group listing and pagination
func TestAdminListGroups(t *testing.T) {
	router, ctx := setupAdmin(t)

	group := &models.Group{Name: "Admin Test Group", CreatedBy: "admin"}
	require.NoError(t, router.appDao.CreateGroup(ctx, group))

	// Test successfully listing groups for a site admin
	t.Run("success", func(t *testing.T) {
		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/admin/groups", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		result, groups := unmarshalHelper[map[string]any](t, body)
		require.Equal(t, 1, result.TotalItems)
		require.Len(t, groups, 1)
	})

	// Test successfully paginating admin groups
	t.Run("pagination", func(t *testing.T) {
		for i := range 4 {
			group := &models.Group{Name: "Paginated Group " + string(rune('A'+i)), CreatedBy: "admin"}
			require.NoError(t, router.appDao.CreateGroup(ctx, group))
		}

		params := url.Values{
			pagination.PageQueryParam:    {"1"},
			pagination.PerPageQueryParam: {"2"},
		}
		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/admin/groups?"+params.Encode(), nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		result, groups := unmarshalHelper[map[string]any](t, body)
		require.Equal(t, 5, result.TotalItems)
		require.Len(t, groups, 2)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestDeleteAdminUser exercises site admin user deletion
func TestDeleteAdminUser(t *testing.T) {
	router, ctx := setupAdmin(t)

	target := &models.User{Username: "delete-me", DisplayName: "Delete Me", SiteRole: types.UserRoleUser}
	createTestUser(t, router, ctx, target)

	status, _, err := requestHelper(t, router, httptest.NewRequest(http.MethodDelete, "/api/admin/users/"+target.ID, nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	deleted, err := router.appDao.GetUser(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: target.ID}))
	require.NoError(t, err)
	require.Nil(t, deleted)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestDeleteAdminGroup exercises site admin group deletion
func TestDeleteAdminGroup(t *testing.T) {
	router, ctx := setupAdmin(t)

	group := &models.Group{Name: "Delete Group", CreatedBy: "admin"}
	require.NoError(t, router.appDao.CreateGroup(ctx, group))

	status, _, err := requestHelper(t, router, httptest.NewRequest(http.MethodDelete, "/api/admin/groups/"+group.ID, nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, status)

	deleted, err := router.appDao.GetGroup(ctx, group.ID)
	require.NoError(t, err)
	require.Nil(t, deleted)
}
