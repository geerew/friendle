package api

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initAdminRoutes initializes admin routes
func (r *Router) initAdminRoutes() {
	a := r.apiGroup("admin")

	// Users
	a.Get("/users", r.require(accessSiteAdmin), r.getUsers)
	a.Post("/users", r.require(accessSiteAdmin), r.createUser)
	a.Put("/users/:id", r.require(accessSiteAdmin), r.updateUser)
	a.Delete("/users/:id", r.require(accessSiteAdmin), r.deleteUser)
	a.Delete("/users/:id/sessions", r.require(accessSiteAdmin), r.deleteUserSessions)

	// Groups
	a.Get("/groups", r.require(accessSiteAdmin), r.getGroups)
	a.Delete("/groups/:id", r.require(accessSiteAdmin), r.deleteGroup)
	a.Post("/groups/:id/members", r.require(accessSiteAdmin), r.createGroupMember)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getUsers returns a paginated site-wide user list
func (r *Router) getUsers(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

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

	pResult, err := dbOpts.Pagination.BuildResult(adminUserResponseHelper(users, userGroupSummariesByUserID(groupRows)))
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createUser creates a user from the site admin API
func (r *Router) createUser(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	userReq := &userRequest{}

	if err := c.BodyParser(userReq); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	if userReq.Username == "" || userReq.Password == "" {
		return errorResponse(c, fiber.StatusBadRequest, "A username and password are required", nil)
	}

	if err := validatePassword(userReq.Password); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	if userReq.Role == "" {
		userReq.Role = types.UserRoleUser.String()
	}

	passwordHash, err := auth.GeneratePassword(userReq.Password)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error hashing password", err)
	}

	user := &models.User{
		Username:     userReq.Username,
		DisplayName:  userReq.Username,
		PasswordHash: passwordHash,
		SiteRole:     types.NewUserRole(userReq.Role),
	}

	if userReq.DisplayName != "" {
		user.DisplayName = userReq.DisplayName
	}

	if err := r.appDao.CreateUser(ctx, user); err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return errorResponse(c, fiber.StatusBadRequest, "Username already exists", nil)
		}

		return errorResponse(c, fiber.StatusInternalServerError, "Error creating user", err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// updateUser updates a user and refreshes sessions when the site role changes
func (r *Router) updateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	userReq := &userRequest{}
	if err := c.BodyParser(userReq); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	if userReq.DisplayName == "" && userReq.Password == "" && userReq.Role == "" {
		return errorResponse(c, fiber.StatusBadRequest, "No data to update", nil)
	}

	_, ctx := principalAndCtx(c)

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: id})
	user, err := r.appDao.GetUser(ctx, dbOpts)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error looking up user", err)
	}

	if user == nil {
		return errorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	if userReq.DisplayName != "" {
		user.DisplayName = userReq.DisplayName
	}

	if userReq.Password != "" {
		if err := validatePassword(userReq.Password); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		passwordHash, err := auth.GeneratePassword(userReq.Password)
		if err != nil {
			return errorResponse(c, fiber.StatusInternalServerError, "Error hashing password", err)
		}
		user.PasswordHash = passwordHash
	}

	if userReq.Role != "" {
		if user.SiteRole.String() == userReq.Role {
			userReq.Role = ""
		} else {
			user.SiteRole = types.NewUserRole(userReq.Role)
		}
	}

	err = r.appDao.UpdateUser(ctx, user)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error updating user", err)
	}

	if userReq.Role != "" {
		if err := r.sessionManager.UpdateSessionRoleForUser(id, user.SiteRole); err != nil {
			return errorResponse(c, fiber.StatusInternalServerError, "Error updating user sessions", err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(&userResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		SiteRole:    user.SiteRole,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroups returns a paginated site-wide group list
func (r *Router) getGroups(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

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

// deleteUser deletes a user and their sessions
func (r *Router) deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	_, ctx := principalAndCtx(c)

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

// deleteUserSessions revokes all sessions for a user
func (r *Router) deleteUserSessions(c *fiber.Ctx) error {
	id := c.Params("id")

	err := r.sessionManager.DeleteUserSessions(id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user sessions", err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteGroup deletes a group
func (r *Router) deleteGroup(c *fiber.Ctx) error {
	id := c.Params("id")

	_, ctx := principalAndCtx(c)

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: id})
	if err := r.appDao.DeleteGroups(ctx, dbOpts); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting group", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroupMember adds a user to a group from the site admin API
func (r *Router) createGroupMember(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

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

	return c.Status(fiber.StatusCreated).JSON(&adminGroupMemberResponse{
		ID:        member.ID,
		UserID:    userID,
		GroupID:   groupID,
		GroupRole: role,
	})
}
