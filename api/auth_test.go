package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAuth_Register exercises user registration
func TestAuth_Register(t *testing.T) {

	// Test successfully registering a new user
	t.Run("201 (created)", func(t *testing.T) {
		router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

		router.app.SetBootstrapped()

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"username": "test", "password": "abcd1234" }`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)

		dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: "test"})
		record, err := router.appDao.GetUser(ctx, dbOpts)
		require.NoError(t, err)
		require.NotEqual(t, "password", record.PasswordHash)
		require.Equal(t, types.SiteRoleUser, record.SiteRole)
	})

	// Test error due to invalid data
	t.Run("400 (bind error)", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Error parsing data")
	})

	// Test error due to missing but valid data data
	t.Run("400 (invalid data)", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		// Missing username and password
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Username and/or password cannot be empty")

		// Missing password
		req = httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"username": "test"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err = requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Username and/or password cannot be empty")

		// Missing username
		req = httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"password": "password123"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err = requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Username and/or password cannot be empty")

		// Empty values
		req = httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"username": "", "password": ""}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err = requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Username and/or password cannot be empty")
	})

	// Test error due to user already existing (case-insensitive)
	t.Run("400 (existing user)", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"username": "test", "password": "abcd1234" }`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, status)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Username already exists")

		req = httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"username": "TEST", "password": "abcd1234" }`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err = requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Username already exists")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAuth_Bootstrap exercises first-run bootstrap
func TestAuth_Bootstrap(t *testing.T) {
	// Test successfully bootstrapping the application
	t.Run("201 (created)", func(t *testing.T) {
		router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

		// Clear the admin user to make it unbootstrapped
		dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: "admin"})
		err := router.appDao.DeleteUsers(ctx, dbOpts)
		require.NoError(t, err)
		require.NoError(t, router.app.RefreshBootstrapped())

		// Generate a bootstrap token using the app's data directory and filesystem
		bootstrapToken, err := auth.GenerateBootstrapToken(router.app.Config.DataDir, router.app.FS)
		require.NoError(t, err)

		// Create user with token
		req := httptest.NewRequest(http.MethodPost, "/api/auth/bootstrap/"+bootstrapToken.Token, strings.NewReader(`{"username": "test", "password": "abcd1234" }`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, _, err2 := requestHelper(t, router, req)
		require.NoError(t, err2)
		require.Equal(t, http.StatusCreated, status)

		dbOpts2 := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: "test"})
		record, err3 := router.appDao.GetUser(ctx, dbOpts2)
		require.NoError(t, err3)
		require.NotEqual(t, "password", record.PasswordHash)
		require.Equal(t, types.SiteRoleAdmin, record.SiteRole)
		require.True(t, router.app.IsBootstrapped())
	})

	t.Run("401 (invalid token)", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		// Clear the admin user to make it unbootstrapped
		dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: "admin"})
		err := router.appDao.DeleteUsers(context.Background(), dbOpts)
		require.NoError(t, err)
		require.NoError(t, router.app.RefreshBootstrapped())

		req := httptest.NewRequest(http.MethodPost, "/api/auth/bootstrap/invalid-token", strings.NewReader(`{"username": "test", "password": "abcd1234" }`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, _, err2 := requestHelper(t, router, req)
		require.NoError(t, err2)
		require.Equal(t, http.StatusUnauthorized, status)
	})

	t.Run("403 (already bootstrapped)", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		// Generate a bootstrap token using the app's data directory and filesystem
		bootstrapToken, err := auth.GenerateBootstrapToken(router.app.Config.DataDir, router.app.FS)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/bootstrap/"+bootstrapToken.Token, strings.NewReader(`{"username": "test", "password": "abcd1234" }`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAuth_Login exercises user login
func TestAuth_Login(t *testing.T) {
	t.Run("200 (success)", func(t *testing.T) {
		router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

		passwordHash, err := auth.GeneratePassword("abcd1234")
		require.NoError(t, err)

		user := &models.User{
			Username:     "test",
			DisplayName:  "Test",
			PasswordHash: passwordHash,
			SiteRole:         types.SiteRoleAdmin,
		}
		require.NoError(t, router.appDao.CreateUser(ctx, user))

		dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: user.ID})
		_, err = router.appDao.GetUser(ctx, dbOpts)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username": "test", "password": "abcd1234" }`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)
	})

	t.Run("400 (bind error)", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Error parsing data")
	})

	t.Run("400 (invalid data)", func(t *testing.T) {
		router, _, _ := setup(t, "admin", types.SiteRoleAdmin)

		// Missing both
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Username and/or password cannot be empty")

		// Missing password
		req = httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username": "test"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err = requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Username and/or password cannot be empty")

		// Missing username
		req = httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"password": "password123"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err = requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Username and/or password cannot be empty")

		// Both empty
		req = httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username": "", "password": ""}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err = requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, status)
		require.Contains(t, string(body), "Username and/or password cannot be empty")
	})

	t.Run("401 (invalid user)", func(t *testing.T) {
		router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

		passwordHash, err := auth.GeneratePassword("abcd1234")
		require.NoError(t, err)

		user := &models.User{
			Username:     "test",
			DisplayName:  "Test",
			PasswordHash: passwordHash,
			SiteRole:         types.SiteRoleAdmin,
		}
		require.NoError(t, router.appDao.CreateUser(ctx, user))

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username": "invalid", "password": "abcd1234" }`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, status)
		require.Contains(t, string(body), "Invalid username and/or password")
	})

	t.Run("401 (invalid password)", func(t *testing.T) {
		router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

		passwordHash, err := auth.GeneratePassword("abcd1234")
		require.NoError(t, err)

		user := &models.User{
			Username:     "test",
			DisplayName:  "Test",
			PasswordHash: passwordHash,
			SiteRole:         types.SiteRoleAdmin,
		}
		require.NoError(t, router.appDao.CreateUser(ctx, user))

		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username": "test", "password": "wrongpass123" }`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusUnauthorized, status)
		require.Contains(t, string(body), "Invalid username and/or password")
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAuth_SignupStatus exercises the signup status endpoint
func TestAuth_SignupStatus(t *testing.T) {
	router, _, _ := setup(t, "", types.SiteRoleUser)

	status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/auth/signup-status", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	var resp signupStatusResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	require.True(t, resp.Enabled)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAuth_GetMe exercises the current user profile endpoint
func TestAuth_GetMe(t *testing.T) {
	router, _, _ := setup(t, "user", types.SiteRoleUser)

	status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/api/auth/me", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	var resp userResponse
	require.NoError(t, json.Unmarshal(body, &resp))
	require.Equal(t, "user", resp.Username)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAuth_UpdateMe exercises self-service profile updates
func TestAuth_UpdateMe(t *testing.T) {
	// Test successfully updating display name
	t.Run("display name", func(t *testing.T) {
		router, ctx, _ := setup(t, "user", types.SiteRoleUser)

		req := httptest.NewRequest(http.MethodPut, "/api/auth/me", strings.NewReader(`{"displayName":"Updated Name"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, status)

		var resp userResponse
		require.NoError(t, json.Unmarshal(body, &resp))
		require.Equal(t, "Updated Name", resp.DisplayName)

		user, err := router.appDao.GetUser(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: "user"}))
		require.NoError(t, err)
		require.Equal(t, "Updated Name", user.DisplayName)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAuth_DeleteMe exercises self-service account deletion
func TestAuth_DeleteMe(t *testing.T) {
	// Test successfully deleting your own account
	t.Run("success", func(t *testing.T) {
		router, ctx, _ := setup(t, "user", types.SiteRoleUser)

		user, err := router.appDao.GetUser(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: "user"}))
		require.NoError(t, err)
		passwordHash, err := auth.GeneratePassword("abcd1234")
		require.NoError(t, err)
		user.PasswordHash = passwordHash
		require.NoError(t, router.appDao.UpdateUser(ctx, user))

		req := httptest.NewRequest(http.MethodDelete, "/api/auth/me", strings.NewReader(`{"currentPassword":"abcd1234"}`))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		status, _, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, status)

		deleted, err := router.appDao.GetUser(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: "user"}))
		require.NoError(t, err)
		require.Nil(t, deleted)
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestAuth_Logout exercises session logout
func TestAuth_Logout(t *testing.T) {
	// Test successfully logging out after login
	t.Run("204", func(t *testing.T) {
		router, ctx, _ := setup(t, "admin", types.SiteRoleAdmin)

		createTestUserWithPassword(t, router, ctx, &models.User{
			Username:    "logout-user",
			DisplayName: "Logout",
			SiteRole:    types.SiteRoleUser,
		}, "abcd1234")

		loginReq := httptest.NewRequest(
			http.MethodPost,
			"/api/auth/login",
			strings.NewReader(`{"username":"logout-user","password":"abcd1234"}`),
		)
		loginReq.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		loginResp, err := router.Test(loginReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, loginResp.StatusCode)

		logoutReq := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		for _, cookie := range loginResp.Cookies() {
			logoutReq.AddCookie(cookie)
		}

		status, _, err := requestHelper(t, router, logoutReq)
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, status)
	})
}
