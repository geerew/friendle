package api

// TODO When 2 windows are open with bootstrap, and the first one bootstraps, the second should
//   1) error on submit or
//   2) redirect to the login page on refresh

import (
	"strings"
	"time"

	"github.com/geerew/friendle/app"
	"github.com/geerew/friendle/utils/logger"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const requestPathLocalsKey = "request_path"

// requestPathInfo classifies the incoming path once per request
type requestPathInfo struct {
	uiAsset bool
	authUI  bool
	api     bool
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// corsMiddleWare creates a CORS middleware
func corsMiddleWare() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, HEAD, PATCH",
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// requestLoggingMiddleware creates a request logging middleware
func requestLoggingMiddleware(log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		// Only log API requests (paths starting with /api)
		path := c.Path()
		if !strings.HasPrefix(path, "/api/") {
			return err
		}

		status := c.Response().StatusCode()
		apiLogger := log.WithComponent(string(app.ComponentAPI))

		// Pull any error info stored by errorResponse
		errMsg, _ := c.Locals("api_error_message").(string)
		errDetail, _ := c.Locals("api_error_detail").(string)

		switch {
		case status >= 500:
			evt := apiLogger.Error().
				Str("method", c.Method()).
				Int("status", status).
				Dur("duration", duration).
				Str("ip", c.IP())
			if errMsg != "" {
				evt = evt.Str("error_message", errMsg)
			}

			if errDetail != "" {
				evt = evt.Str("error_detail", errDetail)
			}
			evt.Msg(c.Path())
		case status >= 400:
			evt := apiLogger.Warn().
				Str("method", c.Method()).
				Int("status", status).
				Dur("duration", duration).
				Str("ip", c.IP())
			if errMsg != "" {
				evt = evt.Str("error_message", errMsg)
			}

			if errDetail != "" {
				evt = evt.Str("error_detail", errDetail)
			}
			evt.Msg(c.Path())
		default:
			apiLogger.Debug().
				Str("method", c.Method()).
				Int("status", status).
				Dur("duration", duration).
				Str("ip", c.IP()).
				Msg(c.Path())
		}

		return err
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// requestPathMiddleware classifies the request path and stores it on the context
func requestPathMiddleware(r *Router) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		c.Locals(requestPathLocalsKey, requestPathInfo{
			uiAsset: r.isUIAsset(path),
			authUI:  strings.HasPrefix(path, "/auth/"),
			api:     strings.HasPrefix(path, "/api/"),
		})

		return c.Next()
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// bootstrapMiddleware checks if the app is bootstrapped. If not, it redirects
// to /auth/bootstrap
//
// Bootstrapping is the process of setting up the app for the first time. It involves
// the creation of 1 admin user, which the /auth/bootstrap endpoint handles
func bootstrapMiddleware(r *Router) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		pathInfo := requestPath(c)

		// If not bootstrapped, force everything through /auth/bootstrap and /api/auth/bootstrap
		if !r.app.IsBootstrapped() {
			if pathInfo.uiAsset {
				return c.Next()
			}

			// API check
			if pathInfo.api {
				if strings.HasPrefix(path, "/api/auth/bootstrap/") {
					c.Locals("bootstrapping", true)
					return c.Next()
				}

				return errorResponse(c, fiber.StatusForbidden, "app is not bootstrapped", nil)
			}

			// UI check
			if strings.HasPrefix(path, "/auth/bootstrap") {
				c.Locals("bootstrapping", true)
				return c.Next()
			}

			return c.Redirect("/auth/bootstrap")
		}

		// If bootstrapped and someone accesses bootstrap URL, redirect appropriately
		if strings.HasPrefix(path, "/auth/bootstrap/") {
			session, err := r.sessionManager.Get(c)

			if err != nil || session.Fresh() {
				return c.Redirect("/auth/login")
			}

			return c.Redirect("/")
		}

		return c.Next()
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// sessionMiddleware loads the session and sets the principal when the caller is logged in
func sessionMiddleware(r *Router) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if bootstrapping, _ := c.Locals("bootstrapping").(bool); bootstrapping {
			return c.Next()
		}

		pathInfo := requestPath(c)
		if pathInfo.uiAsset || strings.HasPrefix(c.Path(), "/api/auth/logout") {
			return c.Next()
		}

		session, err := r.sessionManager.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if session.Fresh() {
			return c.Next()
		}

		userID, ok1 := session.Get("id").(string)
		userRole, ok2 := session.Get("role").(string)
		if !ok1 || !ok2 || userID == "" || userRole == "" {
			return c.Next()
		}

		role := types.NewSiteRole(userRole)
		c.Locals(types.PrincipalContextKey, types.Principal{
			UserID:   userID,
			SiteRole: role,
			Role:     role,
		})

		return c.Next()
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// uiAuthMiddleware redirects browser requests based on login state
func uiAuthMiddleware(r *Router) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if bootstrapping, _ := c.Locals("bootstrapping").(bool); bootstrapping {
			return c.Next()
		}

		pathInfo := requestPath(c)
		if pathInfo.uiAsset || pathInfo.api || strings.HasPrefix(c.Path(), "/api/auth/logout") {
			return c.Next()
		}

		session, err := r.sessionManager.Get(c)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if session.Fresh() {
			if pathInfo.authUI {
				if !r.app.Config.EnableSignup && strings.HasPrefix(c.Path(), "/auth/register") {
					return c.Redirect("/auth/login")
				}

				return c.Next()
			}

			return c.Redirect("/auth/login")
		}

		if pathInfo.authUI {
			return c.Redirect("/")
		}

		if p, ok := c.Locals(types.PrincipalContextKey).(types.Principal); ok {
			if p.SiteRole != types.SiteRoleAdmin && r.isProtectedUIPage(c.Path()) {
				return c.Redirect("/")
			}
		}

		return c.Next()
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Route middleware
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// routeAccess is the access level required for a route
type routeAccess int

const (
	accessAuth routeAccess = iota
	accessSiteAdmin
	accessGroupMember
	accessGroupAdmin
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// require returns per-route API authorization middleware. It assumes sessionMiddleware has
// already attached a principal for logged-in callers; it does not perform UI redirects
func (r *Router) require(level routeAccess) fiber.Handler {
	return func(c *fiber.Ctx) error {
		switch level {
		case accessAuth:
			if _, _, err := principalCtx(c); err != nil {
				return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
			}

		case accessSiteAdmin:
			p, _, err := principalCtx(c)
			if err != nil {
				return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
			}

			if p.SiteRole != types.SiteRoleAdmin {
				return errorResponse(c, fiber.StatusForbidden, "Site admin required", nil)
			}

		case accessGroupMember:
			p, ctx, err := principalCtx(c)
			if err != nil {
				return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
			}

			if _, err := r.membership(ctx, c.Params("id"), p.UserID); err != nil {
				return errorResponse(c, fiber.StatusForbidden, "Not a member", nil)
			}

		case accessGroupAdmin:
			p, ctx, err := principalCtx(c)
			if err != nil {
				return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
			}

			m, err := r.appDao.GetGroupMember(ctx, c.Params("id"), p.UserID)
			if err != nil || m == nil || m.GroupRole != types.GroupRoleAdmin {
				return errorResponse(c, fiber.StatusForbidden, "Group admin required", nil)
			}

		default:
			return errorResponse(c, fiber.StatusInternalServerError, "Unknown access level", nil)
		}

		return c.Next()
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// requestPath returns path classification stored by requestPathMiddleware
func requestPath(c *fiber.Ctx) requestPathInfo {
	info, _ := c.Locals(requestPathLocalsKey).(requestPathInfo)

	return info
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// isUIAsset reports whether the path is a dev/prod UI bundle asset or static file
func (r *Router) isUIAsset(path string) bool {
	if r.app.Config.AppMode == app.AppModeDev &&
		(strings.HasPrefix(path, "/node_modules/") ||
			strings.HasPrefix(path, "/.svelte-kit/") ||
			strings.HasPrefix(path, "/src/") ||
			strings.HasPrefix(path, "/@")) {
		return true
	}

	if r.app.Config.AppMode != app.AppModeDev && strings.HasPrefix(path, "/_app/") {
		return true
	}

	if strings.HasPrefix(path, "/apple-touch-icon.png") ||
		strings.HasPrefix(path, "/favicon.") ||
		strings.HasPrefix(path, "/fonts/") ||
		strings.HasPrefix(path, "/web-app-manifest-") {
		return true
	}

	return false
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// isProtectedUIPage checks if the request is intended for a protected UI page
func (r *Router) isProtectedUIPage(path string) bool {
	return strings.HasPrefix(path, "/admin")
}
