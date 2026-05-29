package api

import (
	"github.com/geerew/friendle/app"
	"github.com/geerew/friendle/version"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initVersionRoutes initializes the version routes
func (r *Router) initVersionRoutes() {
	g := r.apiGroup("version")

	g.Get("", r.getVersion)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getVersion returns the application version
func (r *Router) getVersion(c *fiber.Ctx) error {
	// If running in dev mode, always return "dev" regardless of build version
	var currentVersion string
	if r.app.Config.AppMode == app.AppModeDev {
		currentVersion = "dev"
	} else {
		currentVersion = version.GetVersion()
	}

	return c.Status(fiber.StatusOK).JSON(&versionResponse{Version: currentVersion})
}
