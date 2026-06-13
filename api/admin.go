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

	r.initAdminUserRoutes(a)

	a.Post("/recovery", r.adminCreateRecovery)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// adminCreateRecovery applies a password reset from a recovery token written by the CLI
func (r *Router) adminCreateRecovery(c *fiber.Ctx) error {
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
