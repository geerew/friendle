package api

import (
	"github.com/geerew/friendle/app"
	"github.com/geerew/friendle/version"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type versionAPI struct {
	r *Router
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initVersionRoutes initializes the version routes
func (r *Router) initVersionRoutes() {
	versionAPI := versionAPI{
		r: r,
	}

	g := r.apiGroup("version")

	g.Get("", versionAPI.getVersion)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getVersion returns the application version
func (api *versionAPI) getVersion(c *fiber.Ctx) error {
	// If running in dev mode, always return "dev" regardless of build version
	var currentVersion string
	if api.r.app.Config.AppMode == app.AppModeDev {
		currentVersion = "dev"
	} else {
		currentVersion = version.GetVersion()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"version": currentVersion,
	})
}
