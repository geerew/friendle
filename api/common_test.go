package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/geerew/friendle/app"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// testPrincipal holds mutable auth state for test middleware
type testPrincipal struct {
	userID string
	role   types.UserRole
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// setup creates a test router
func setup(t *testing.T, id string, role types.UserRole) (*Router, context.Context, *testPrincipal) {
	t.Helper()

	appConfig := &app.Config{
		HttpAddr:     "127.0.0.1:9081",
		DataDir:      t.TempDir(),
		AppMode:      app.AppModeTest,
		EnableSignup: true,
	}
	application, err := app.New(context.Background(), appConfig)
	require.NoError(t, err)

	principal := &testPrincipal{userID: id, role: role}

	stack := testMiddleware(principal)
	if id == "" {
		stack = func(r *Router) []fiber.Handler {
			return []fiber.Handler{corsMiddleWare()}
		}
	}

	router := New(application, stack)

	if id != "" {
		user := models.User{
			Base: models.Base{
				ID: id,
			},
			Username:     id,
			SiteRole:     role,
			PasswordHash: "password",
			DisplayName:  "Test User",
		}

		require.NoError(t, router.appDao.CreateUser(context.Background(), &user))
	}

	if id != "" {
		router.app.SetBootstrapped()
	}

	ctx := context.Background()
	if id != "" {
		ctx = context.WithValue(ctx, types.PrincipalContextKey, types.Principal{
			UserID:   id,
			SiteRole: role,
		})
	}

	return router, ctx, principal
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// testMiddleware returns a minimal stack that injects the test principal
func testMiddleware(principal *testPrincipal) MiddlewareStack {
	return func(r *Router) []fiber.Handler {
		return []fiber.Handler{
			corsMiddleWare(),
			requestPathMiddleware(r),
			bootstrapMiddleware(r),
			func(c *fiber.Ctx) error {
				if principal.userID != "" {
					c.Locals(types.PrincipalContextKey, types.Principal{
						UserID:   principal.userID,
						SiteRole: principal.role,
					})
				}

				return c.Next()
			},
		}
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// requestHelper sends a request to the router and returns the status code and body
func requestHelper(t *testing.T, router *Router, req *http.Request) (int, []byte, error) {
	t.Helper()

	resp, err := router.Test(req)
	if err != nil {
		return -1, nil, err
	}

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body, err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// unmarshalHelper unmarshals a body into a pagination result and a slice of items
func unmarshalHelper[T any](t *testing.T, body []byte) (pagination.PaginationResult, []T) {
	t.Helper()

	var respData pagination.PaginationResult
	err := json.Unmarshal(body, &respData)
	require.NoError(t, err)

	var resp []T
	for _, item := range respData.Items {
		var r T
		require.Nil(t, json.Unmarshal(item, &r))
		resp = append(resp, r)
	}

	return respData, resp
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createTestUser inserts a user, using a placeholder password hash when none is set
func createTestUser(t *testing.T, router *Router, ctx context.Context, user *models.User) {
	t.Helper()

	if user.PasswordHash == "" {
		user.PasswordHash = "test-password-hash"
	}

	require.NoError(t, router.appDao.CreateUser(ctx, user))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createTestUserWithPassword inserts a user with a hashed password
func createTestUserWithPassword(t *testing.T, router *Router, ctx context.Context, user *models.User, password string) {
	t.Helper()

	passwordHash, err := auth.GeneratePassword(password)
	require.NoError(t, err)

	user.PasswordHash = passwordHash
	require.NoError(t, router.appDao.CreateUser(ctx, user))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createTestGroupWithMember creates a group and adds the user as a member
func createTestGroupWithMember(t *testing.T, router *Router, ctx context.Context, userID string, groupRole types.GroupRole, name string) *models.Group {
	t.Helper()

	group := &models.Group{Name: name, CreatedBy: userID, IntervalHours: 24, Timezone: "UTC"}
	require.NoError(t, router.appDao.CreateGroup(ctx, group))

	groupMember := &models.GroupMember{
		GroupID:   group.ID,
		UserID:    userID,
		GroupRole: groupRole,
	}
	require.NoError(t, router.appDao.CreateGroupMember(ctx, groupMember))

	return group
}
