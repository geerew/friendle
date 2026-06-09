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
	groupRoutes.Get("/self", r.requireAccess(accessSiteUser), r.listSelfGroups)
	groupRoutes.Get("/", r.requireAccess(accessSiteUser), r.listGroups)
	groupRoutes.Get("/:id", r.requireAccess(accessSiteUser), r.getGroup)
	// TODO: support group admin
	groupRoutes.Delete("/:id", r.requireAccess(accessSiteAdmin), r.deleteGroup)
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroup returns a group by ID
func (r *Router) getGroup(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	group, err := r.appSvc.Groups.Get(ctx, c.Params("id"))
	if err != nil {
		return serviceError(c, err)
	}

	return c.JSON(group)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// listGroups returns paginated groups, optionally filtered by name
func (r *Router) listGroups(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	page := paginationFromCtx(c)

	var (
		groups any
		err    error
	)

	if c.Context().QueryArgs().Has("name") {
		groups, err = r.appSvc.Groups.Search(ctx, page, c.Query("name", ""))
	} else {
		groups, err = r.appSvc.Groups.List(ctx, page)
	}

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

// listSelfGroups returns paginated groups the authenticated user belongs to
func (r *Router) listSelfGroups(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	page := paginationFromCtx(c)
	groups, err := r.appSvc.Groups.ListSelf(ctx, page)
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

// deleteGroup deletes a group
func (r *Router) deleteGroup(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	if err := r.appSvc.Groups.Delete(ctx, c.Params("id")); err != nil {
		return serviceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
