package api

import (
	"github.com/geerew/friendle/service"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initGroupRoutes initializes the group routes
func (r *Router) initGroupRoutes() {
	groupRoutes := r.apiGroup("groups")

	groupRoutes.Post("/", r.requireAccess(accessSiteUser), r.createGroup)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroup creates a friend group for the authenticated user
func (r *Router) createGroup(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	req := &service.CreateGroupRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	group, err := r.appSvc.Groups.Create(ctx, *req)
	if err != nil {
		return serviceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}
