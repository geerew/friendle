package api

import (
	"errors"
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

// initUserRoutes initializes the user routes
func (r *Router) initUserRoutes() {
	g := r.apiGroup("users")

	g.Get("", protectedRoute, r.getUsers)
	g.Post("", protectedRoute, r.createUser)
	g.Put("/:id", protectedRoute, r.updateUser)
	g.Delete("/:id", protectedRoute, r.deleteUser)

	g.Delete("/:id/sessions", protectedRoute, r.deleteUserSession)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getUsers returns a paginated list of users
func (r *Router) getUsers(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Missing principal", nil)
	}

	dbOpts := dao.NewOptions().
		WithOrderBy(utils.StringSplit(c.Query("orderBy", ""), ",")...).
		WithApiQuery(c.Query("q", "")).
		WithPagination(paginationFromCtx(c))

	users, err := r.appDao.ListUsers(ctx, dbOpts)
	if err != nil {
		if errors.Is(err, utils.ErrApiQueryParse) {
			return errorResponse(c, fiber.StatusBadRequest, "Error parsing query", err)
		}

		return errorResponse(c, fiber.StatusInternalServerError, "Error looking up users", err)
	}

	pResult, err := dbOpts.Pagination.BuildResult(userResponseHelper(users))
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.Status(fiber.StatusOK).JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createUser creates a user as site admin
func (r *Router) createUser(c *fiber.Ctx) error {
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

	// Default the role to a user when not provided
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

	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Missing principal", nil)
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

// updateUser updates a user and revokes sessions when the role changes
func (r *Router) updateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	userReq := &userRequest{}
	if err := c.BodyParser(userReq); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	if userReq.DisplayName == "" && userReq.Password == "" && userReq.Role == "" {
		return errorResponse(c, fiber.StatusBadRequest, "No data to update", nil)
	}

	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Missing principal", nil)
	}

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
		// Do nothing when the role is the same
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

	// Update all the sessions for the user with the new role
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

// deleteUser deletes a user and their sessions
func (r *Router) deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Missing principal", nil)
	}

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: id})
	err = r.appDao.DeleteUsers(ctx, dbOpts)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user", err)
	}

	err = r.sessionManager.DeleteUserSessions(id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user sessions", err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteUserSession revokes all sessions for a user
func (r *Router) deleteUserSession(c *fiber.Ctx) error {
	id := c.Params("id")

	err := r.sessionManager.DeleteUserSessions(id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user sessions", err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
