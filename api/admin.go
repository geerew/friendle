package api

import (
	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) initAdminRoutes() {
	a := r.apiGroup("admin")
	a.Get("/users", r.requireAuth, r.requireSiteAdmin, r.adminListUsers)
	a.Get("/groups", r.requireAuth, r.requireSiteAdmin, r.adminListGroups)
	a.Delete("/users/:id", r.requireAuth, r.requireSiteAdmin, r.adminDeleteUser)
	a.Delete("/groups/:id", r.requireAuth, r.requireSiteAdmin, r.adminDeleteGroup)
}

func (r *Router) adminListUsers(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	dbOpts := dao.NewOptions().WithPagination(paginationFromCtx(c))
	users, err := r.appDao.ListAdminUsers(ctx, dbOpts)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "List failed", err)
	}

	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	groupRows, err := r.appDao.ListUserGroupSummariesForUserIDs(ctx, userIDs)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "List failed", err)
	}

	pResult, err := dbOpts.Pagination.BuildResult(
		adminUserResponseHelper(users, userGroupSummariesByUserID(groupRows)),
	)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

func (r *Router) adminListGroups(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	dbOpts := dao.NewOptions().WithPagination(paginationFromCtx(c))
	groups, err := r.appDao.ListAdminGroups(ctx, dbOpts)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "List failed", err)
	}

	pResult, err := dbOpts.Pagination.BuildResult(adminGroupResponseHelper(groups))
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

func (r *Router) adminDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: id})
	if err := r.appDao.DeleteUsers(ctx, dbOpts); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user", err)
	}

	if err := r.sessionManager.DeleteUserSessions(id); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user sessions", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (r *Router) adminDeleteGroup(c *fiber.Ctx) error {
	id := c.Params("id")

	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: id})
	if err := r.appDao.DeleteGroups(ctx, dbOpts); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting group", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
