package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/types"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAdminRecovery exercises admin password recovery via token file
func TestAdminRecovery(t *testing.T) {
	// Test successfully resetting an admin password
	t.Run("200", func(t *testing.T) {
		router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

		recoveryToken, err := auth.GenerateRecoveryToken(router.app.FS, "admin", "newpass1234", router.app.Config.DataDir)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/recovery", strings.NewReader(`{"token":"`+recoveryToken.Token+`"}`))
		req.Header.Set("Content-Type", "application/json")

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		record, err := router.appDao.GetUser(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: "admin"}))
		require.NoError(t, err)
		require.Equal(t, recoveryToken.PasswordHash, record.PasswordHash)

		_, err = auth.ValidateRecoveryToken(router.app.FS, recoveryToken.Token, router.app.Config.DataDir)
		require.Error(t, err)
	})

	// Test error due to an invalid recovery token
	t.Run("401", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		req := httptest.NewRequest(http.MethodPost, "/api/admin/recovery", strings.NewReader(`{"token":"invalid"}`))
		req.Header.Set("Content-Type", "application/json")

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, status)
		require.Contains(t, string(body), "Invalid or expired recovery token")
	})
}
