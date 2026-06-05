package api

import (
	"github.com/geerew/friendle/service"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initGroupRoutes initializes the group routes
func (r *Router) initGroupRoutes() {
	groupRoutes := r.apiGroup("groups")

	groupRoutes.Get("/self", r.requireAccess(accessSiteUser), r.listSelfGroups)
	groupRoutes.Post("/", r.requireAccess(accessSiteUser), r.createGroup)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// listSelfGroups returns paginated groups the authenticated user belongs to
func (r *Router) listSelfGroups(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	page := paginationFromCtx(c)
	groups, err := r.appSvc.Groups.ListSelfGroups(ctx, page)
	if err != nil {
		return serviceError(c, err)
	}

	pResult, err := page.BuildResult(groups)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroup creates a friend group for the authenticated user
func (r *Router) createGroup(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	req := &service.CreateGroupRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	group, err := r.appSvc.Groups.CreateGroup(ctx, *req)
	if err != nil {
		return serviceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}
