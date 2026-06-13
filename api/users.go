package api

import (
	"github.com/geerew/friendle/service"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initAdminUserRoutes registers site-admin user routes on the admin API group
func (r *Router) initAdminUserRoutes(adminGroup fiber.Router) {
	userRoutes := adminGroup.Group("/users")

	userRoutes.Post("/", r.requireAccess(accessSiteAdmin), r.createUser)
	userRoutes.Get("/", r.requireAccess(accessSiteAdmin), r.listUsers)
	userRoutes.Put("/:id", r.requireAccess(accessSiteAdmin), r.updateUser)
	userRoutes.Delete("/:id", r.requireAccess(accessSiteAdmin), r.deleteUser)
	userRoutes.Delete("/:id/sessions", r.requireAccess(accessSiteAdmin), r.deleteUserSessions)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createUser creates a user from the site admin API
func (r *Router) createUser(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	req := &service.UserCreateRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	err := r.appSvc.Users.Create(ctx, *req)
	if err != nil {
		return serviceError(c, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// listUsers returns a paginated site-wide user list
func (r *Router) listUsers(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	page := paginationFromCtx(c)
	users, err := r.appSvc.Users.List(ctx, page)
	if err != nil {
		return serviceError(c, err)
	}

	pResult, err := page.BuildResult(users)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// updateUser updates a user and refreshes sessions when the site role changes
func (r *Router) updateUser(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	req := &service.UserUpdateRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	view, roleChanged, err := r.appSvc.Users.Update(ctx, c.Params("id"), *req)
	if err != nil {
		return serviceError(c, err)
	}

	if roleChanged {
		if err := r.sessionManager.UpdateSessionRoleForUser(view.ID, view.SiteRole); err != nil {
			return errorResponse(c, fiber.StatusInternalServerError, "Error updating user sessions", err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(view)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteUser deletes a user and their sessions
func (r *Router) deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	_, ctx := principalAndCtx(c)

	if err := r.appSvc.Users.Delete(ctx, id); err != nil {
		return serviceError(c, err)
	}

	if err := r.sessionManager.DeleteUserSessions(id); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user sessions", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteUserSessions revokes all sessions for a user
func (r *Router) deleteUserSessions(c *fiber.Ctx) error {
	id := c.Params("id")

	err := r.sessionManager.DeleteUserSessions(id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user sessions", err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
