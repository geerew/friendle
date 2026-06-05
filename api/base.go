package api

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/geerew/friendle/app"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/service"
	"github.com/geerew/friendle/utils/logger"
	"github.com/geerew/friendle/utils/session"
	"github.com/gofiber/fiber/v2"
	fibersession "github.com/gofiber/fiber/v2/middleware/session"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MiddlewareStack builds the global middleware chain for a router
type MiddlewareStack func(r *Router) []fiber.Handler

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Router defines a router
type Router struct {
	fiberApp       *fiber.App
	app            *app.App
	appDao         *dao.DAO
	appSvc         *service.Service
	sessionManager *session.SessionManager
	logger         *logger.Logger
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// New creates a router, applies middleware, and registers routes
func New(application *app.App, stack MiddlewareStack) *Router {
	if stack == nil {
		stack = DefaultMiddleware
	}

	log := application.Logger.WithComponent(string(app.ComponentAPI))
	r := &Router{
		app:    application,
		appDao: dao.New(application.DbManager.DataDb),
		appSvc: service.New(application.DbManager.DataDb),
		logger: log,
	}

	r.createSessionStore()

	r.fiberApp = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	for _, handler := range stack(r) {
		r.fiberApp.Use(handler)
	}

	r.bindUi()
	r.initAuthRoutes()
	r.initGroupRoutes()
	r.initAdminRoutes()
	r.initVersionRoutes()

	return r
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Serve serves the API and UI
func (r *Router) Serve() error {
	ln, err := net.Listen("tcp", r.app.Config.HttpAddr)
	if err != nil {
		return err
	}

	r.logger.Info().
		Str("url", fmt.Sprintf("http://%s", r.app.Config.HttpAddr)).
		Msg("Server started")

	return r.fiberApp.Listener(ln)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Test sends a request through the router for use in tests
func (r *Router) Test(req *http.Request, msTimeout ...int) (*http.Response, error) {
	return r.fiberApp.Test(req, msTimeout...)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createSessionStore creates the session store
func (r *Router) createSessionStore() {
	config := fibersession.Config{
		KeyLookup:      "cookie:session",
		Expiration:     7 * (24 * time.Hour),
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	}

	sqliteStorage := session.NewSqliteStorage(r.app.DbManager.DataDb, 10*time.Second)

	r.sessionManager = session.New(r.app.DbManager.DataDb, config, sqliteStorage)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// apiGroup returns a new API router group
func (r *Router) apiGroup(groupPath string) fiber.Router {
	return r.fiberApp.Group("/api/" + groupPath)
}
