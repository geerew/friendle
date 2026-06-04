package api

import (
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initAdminRoutes initializes admin routes
func (r *Router) initAdminRoutes() {
	a := r.apiGroup("admin")

	a.Post("/groups/:id/members", r.requireAccess(accessSiteAdmin), r.createGroupMember)

	// Recovery (unauthenticated; requires a valid .recovery-token on disk)
	a.Post("/recovery", r.createAdminRecovery)
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

	// Get the group
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: groupID})
	g, err := r.appDao.GetGroup(ctx, dbOpts)
	if err != nil || g == nil {
		return errorResponse(c, fiber.StatusNotFound, "Group not found", nil)
	}

	// Check if the user is already a member of the group
	dbOpts = dao.NewOptions().
		WithWhere(squirrel.Eq{
			models.GROUP_MEMBER_GROUP_ID: groupID,
			models.GROUP_MEMBER_USER_ID:  userID,
		})

	if m, _ := r.appDao.GetGroupMember(ctx, dbOpts); m != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Already a member", nil)
	}

	role := types.NewGroupRole(req.GroupRole)
	if !role.IsValid() {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid group role", nil)
	}

	// Delete any pending join requests for this user and group
	dbOpts = dao.NewOptions().WithWhere(squirrel.Eq{
		models.JOIN_REQUEST_GROUP_ID: groupID,
		models.JOIN_REQUEST_USER_ID:  userID,
		models.JOIN_REQUEST_STATUS:   types.JoinPending,
	})
	_ = r.appDao.DeleteJoinRequests(ctx, dbOpts)

	// Create the group member
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createAdminRecovery applies a password reset from a recovery token written by the CLI
func (r *Router) createAdminRecovery(c *fiber.Ctx) error {
	req := &adminRecoveryRequest{}
	if err := c.BodyParser(req); err != nil || strings.TrimSpace(req.Token) == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Token required", nil)
	}

	recoveryToken, err := auth.ValidateRecoveryToken(r.app.FS, strings.TrimSpace(req.Token), r.app.Config.DataDir)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Invalid or expired recovery token", err)
	}

	ctx := c.UserContext()
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: recoveryToken.Username})
	user, err := r.appDao.GetUser(ctx, dbOpts)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error looking up user", err)
	}

	if user == nil {
		return errorResponse(c, fiber.StatusNotFound, "User not found", nil)
	}

	if user.SiteRole != types.SiteRoleAdmin {
		return errorResponse(c, fiber.StatusForbidden, "User is not an admin", nil)
	}

	user.PasswordHash = recoveryToken.PasswordHash
	if err := r.appDao.UpdateUser(ctx, user); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error updating password", err)
	}

	if err := r.sessionManager.DeleteUserSessions(user.ID); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user sessions", err)
	}

	if err := auth.DeleteRecoveryToken(r.app.FS, r.app.Config.DataDir); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting recovery token", err)
	}

	return c.SendStatus(fiber.StatusOK)
}
