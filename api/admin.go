package api

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initAdminRoutes initializes admin routes
func (r *Router) initAdminRoutes() {
	a := r.apiGroup("admin")

	// Users
	a.Get("/users", r.requireAuth, r.requireSiteAdmin, r.getAdminUsers)
	a.Delete("/users/:id", r.requireAuth, r.requireSiteAdmin, r.deleteAdminUser)

	// Groups
	a.Get("/groups", r.requireAuth, r.requireSiteAdmin, r.getAdminGroups)
	a.Delete("/groups/:id", r.requireAuth, r.requireSiteAdmin, r.deleteAdminGroup)
	a.Post("/groups/:id/members", r.requireAuth, r.requireSiteAdmin, r.createAdminGroupMember)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getAdminUsers returns paginated users for site admins
func (r *Router) getAdminUsers(c *fiber.Ctx) error {
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getAdminGroups returns paginated groups for site admins
func (r *Router) getAdminGroups(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	dbOpts := dao.NewOptions().
		WithOrderBy(utils.StringSplit(c.Query("orderBy", ""), ",")...).
		WithPagination(paginationFromCtx(c))
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteAdminUser deletes a user and their sessions
func (r *Router) deleteAdminUser(c *fiber.Ctx) error {
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteAdminGroup deletes a group
func (r *Router) deleteAdminGroup(c *fiber.Ctx) error {
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createAdminGroupMember adds a user to a group as site admin
func (r *Router) createAdminGroupMember(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	req := &adminAddGroupMemberRequest{}
	if err := c.BodyParser(req); err != nil || strings.TrimSpace(req.UserID) == "" {
		return errorResponse(c, fiber.StatusBadRequest, "userId required", nil)
	}

	groupID := c.Params("id")
	userID := strings.TrimSpace(req.UserID)

	g, err := r.appDao.GetGroup(ctx, groupID)
	if err != nil || g == nil {
		return errorResponse(c, fiber.StatusNotFound, "Group not found", nil)
	}

	if m, _ := r.appDao.GetGroupMember(ctx, groupID, userID); m != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Already a member", nil)
	}

	role := types.NewGroupRole(req.GroupRole)
	if !role.IsValid() {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid group role", nil)
	}

	_ = r.appDao.DeletePendingJoinRequest(ctx, groupID, userID)

	member := &models.GroupMember{GroupID: groupID, UserID: userID, GroupRole: role}
	if err := r.appDao.CreateGroupMember(ctx, member); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to add member", err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":        member.ID,
		"userId":    userID,
		"groupId":   groupID,
		"groupRole": role,
	})
}
