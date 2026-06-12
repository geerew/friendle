package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geerew/friendle/app"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TestBootstrapMiddleware exercises unbootstrapped route guards
func TestBootstrapMiddleware(t *testing.T) {
	router := setupUnbootstrappedRouter(t)
	bootstrapPath := router.app.BootstrapPath()
	require.NotEmpty(t, bootstrapPath)

	// Test error due to a request for the site root while unbootstrapped
	t.Run("403 root", func(t *testing.T) {
		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, status)
		require.Contains(t, string(body), "app is not bootstrapped")
	})

	// Test error due to a bare bootstrap path without the token
	t.Run("404 bootstrap prefix", func(t *testing.T) {
		status, body, err := requestHelper(t, router, httptest.NewRequest(http.MethodGet, "/auth/bootstrap/", nil))
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, status)
		require.Contains(t, string(body), "Not found")
	})

	// Test error due to an incorrect bootstrap token in the path
	t.Run("404 wrong token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/auth/bootstrap/wrong-token/", nil)
		status, body, err := requestHelper(t, router, req)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, status)
		require.Contains(t, string(body), "Not found")
		require.NotContains(t, string(body), bootstrapPath)
	})
}

// setupUnbootstrappedRouter creates a router awaiting first-run bootstrap
func setupUnbootstrappedRouter(t *testing.T) *Router {
	t.Helper()

	appConfig := &app.Config{
		HttpAddr:     "127.0.0.1:9081",
		DataDir:      t.TempDir(),
		AppMode:      app.AppModeTest,
		EnableSignup: true,
	}
	application, err := app.New(context.Background(), appConfig)
	require.NoError(t, err)
	require.False(t, application.IsBootstrapped())

	stack := func(r *Router) []fiber.Handler {
		return []fiber.Handler{
			corsMiddleWare(),
			requestPathMiddleware(r),
			bootstrapMiddleware(r),
		}
	}

	return New(application, stack)
}
