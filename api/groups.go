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
	groupRoutes.Get("/", r.requireAccess(accessSiteUser), r.listGroups)
	groupRoutes.Get("/self", r.requireAccess(accessSiteUser), r.listSelfGroups)
	groupRoutes.Get("/search", r.requireAccess(accessSiteUser), r.searchGroups)
	groupRoutes.Get("/:id", r.requireAccess(accessSiteUser), r.getGroup)
	groupRoutes.Delete("/:id", r.requireAccess(accessSiteAdmin), r.deleteGroup)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// listGroups returns paginated groups
func (r *Router) listGroups(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	page := paginationFromCtx(c)
	groups, err := r.appSvc.Groups.ListGroups(ctx, page)
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

// searchGroups returns paginated groups matching a name query
func (r *Router) searchGroups(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	page := paginationFromCtx(c)
	groups, err := r.appSvc.Groups.SearchGroups(ctx, page, c.Query("q", ""))
	if err != nil {
		return serviceError(c, err)
	}

	pResult, err := page.BuildResult(groups)
	if err != nil {
		return serviceError(c, err)
	}

	return c.JSON(pResult)
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteGroup deletes a group
func (r *Router) deleteGroup(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	if err := r.appSvc.Groups.DeleteGroup(ctx, c.Params("id")); err != nil {
		return serviceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroup returns a group the authenticated user belongs to
func (r *Router) getGroup(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	group, err := r.appSvc.Groups.GetGroup(ctx, c.Params("id"))
	if err != nil {
		return serviceError(c, err)
	}

	return c.JSON(group)
}
