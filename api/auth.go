package api

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initAuthRoutes initializes the auth routes
func (r *Router) initAuthRoutes() {
	authGroup := r.apiGroup("auth")

	authGroup.Get("/signup-status", r.signupStatus)
	authGroup.Post("/bootstrap/:token", r.bootstrap)
	authGroup.Post("/register", r.register)
	authGroup.Post("/login", r.login)
	authGroup.Post("/logout", r.logout)

	authGroup.Get("/me", r.requireAccess(accessAuth), r.getMe)
	authGroup.Put("/me", r.requireAccess(accessAuth), r.updateMe)
	authGroup.Delete("/me", r.requireAccess(accessAuth), r.deleteMe)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// signupStatus returns whether self-service registration is enabled
func (r *Router) signupStatus(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(signupStatusResponse{
		Enabled: r.app.Config.EnableSignup,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// bootstrap creates the first admin user using a one-time bootstrap token
func (r *Router) bootstrap(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Bootstrap token is required", nil)
	}

	// Check if already bootstrapped first
	if r.app.IsBootstrapped() {
		return errorResponse(c, fiber.StatusForbidden, "Application is already bootstrapped", nil)
	}

	// Validate bootstrap token
	if err := auth.ValidateBootstrapToken(token, r.app.Config.DataDir, r.app.FS); err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Invalid or expired bootstrap token", nil)
	}

	// Create admin user using existing register logic
	err := r.register(c)
	if err == nil {
		r.app.SetBootstrapped()
		auth.DeleteBootstrapToken(r.app.Config.DataDir, r.app.FS)
	}

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// register creates a user account and starts a session
func (r *Router) register(c *fiber.Ctx) error {
	if r.app.IsBootstrapped() && !r.app.Config.EnableSignup {
		return errorResponse(c, fiber.StatusForbidden, "Sign-up is disabled", nil)
	}

	registerReq := &registerRequest{}

	if err := c.BodyParser(registerReq); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	if registerReq.Username == "" || registerReq.Password == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Username and/or password cannot be empty", nil)
	}

	if err := validatePassword(registerReq.Password); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	passwordHash, err := auth.GeneratePassword(registerReq.Password)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error hashing password", err)
	}

	user := &models.User{
		Username:     registerReq.Username,
		DisplayName:  registerReq.Username, // Set the display name to the username by default
		PasswordHash: passwordHash,
	}

	// The first user will always be an admin
	if !r.app.IsBootstrapped() {
		user.SiteRole = types.UserRoleAdmin
	} else {
		user.SiteRole = types.UserRoleUser
	}

	err = r.appDao.CreateUser(c.UserContext(), user)
	if err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return errorResponse(c, fiber.StatusBadRequest, "Username already exists", nil)
		}

		return errorResponse(c, fiber.StatusInternalServerError, "Error creating user", err)
	}

	err = r.sessionManager.SetSession(c, user.ID, user.SiteRole)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error setting session", err)
	}

	return c.Status(fiber.StatusCreated).JSON(&userResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		SiteRole:    user.SiteRole,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// login authenticates a user and starts a session
func (r *Router) login(c *fiber.Ctx) error {
	loginReq := &loginRequest{}

	if err := c.BodyParser(loginReq); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	if loginReq.Username == "" || loginReq.Password == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Username and/or password cannot be empty", nil)
	}

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: loginReq.Username})
	user, err := r.appDao.GetUser(c.UserContext(), dbOpts)
	if err != nil || user == nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Invalid username and/or password", nil)
	}

	if !auth.ComparePassword(user.PasswordHash, loginReq.Password) {
		return errorResponse(c, fiber.StatusUnauthorized, "Invalid username and/or password", nil)
	}

	err = r.sessionManager.SetSession(c, user.ID, user.SiteRole)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error setting session", err)
	}

	return c.Status(fiber.StatusOK).JSON(&userResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		SiteRole:    user.SiteRole,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// logout ends the current session
func (r *Router) logout(c *fiber.Ctx) error {
	err := r.sessionManager.DeleteSession(c)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting session", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getMe returns the authenticated user's profile
func (r *Router) getMe(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	user, err := r.getUserByPrincipal(ctx, principal)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error getting user information", err)
	}

	return c.Status(fiber.StatusOK).JSON(&userResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		SiteRole:    user.SiteRole,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// updateMe updates the authenticated user's profile or password
func (r *Router) updateMe(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	user, err := r.getUserByPrincipal(ctx, principal)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error getting user information", err)
	}

	updateReq := &selfUpdateRequest{}
	if err := c.BodyParser(updateReq); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	if updateReq.DisplayName == "" && updateReq.Password == "" {
		return errorResponse(c, fiber.StatusBadRequest, "No data to update", nil)
	}

	if updateReq.DisplayName != "" {
		user.DisplayName = updateReq.DisplayName
	}

	if updateReq.Password != "" {
		if !auth.ComparePassword(user.PasswordHash, updateReq.CurrentPassword) {
			return errorResponse(c, fiber.StatusBadRequest, "Invalid current password", nil)
		}

		passwordHash, err := auth.GeneratePassword(updateReq.Password)
		if err != nil {
			return errorResponse(c, fiber.StatusInternalServerError, "Error hashing password", err)
		}
		user.PasswordHash = passwordHash
	}

	err = r.appDao.UpdateUser(ctx, user)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error updating user", err)
	}

	return c.Status(fiber.StatusOK).JSON(&userResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		SiteRole:    user.SiteRole,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteMe deletes the authenticated user's account
func (r *Router) deleteMe(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	user, err := r.getUserByPrincipal(ctx, principal)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error getting user information", err)
	}

	deleteReq := &selfDeleteRequest{}
	if err := c.BodyParser(deleteReq); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	if !auth.ComparePassword(user.PasswordHash, deleteReq.CurrentPassword) {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid password", nil)
	}

	if user.SiteRole == types.UserRoleAdmin {
		// Count the number of admin users and fail if there is only one
		dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_SITE_ROLE: types.UserRoleAdmin})
		adminCount, err := r.appDao.CountUsers(ctx, dbOpts)
		if err != nil {
			return errorResponse(c, fiber.StatusInternalServerError, "Error counting admin users", err)
		}

		if adminCount == 1 {
			return errorResponse(c, fiber.StatusBadRequest, "Unable to delete the last admin user", nil)
		}
	}

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: principal.UserID})
	err = r.appDao.DeleteUsers(ctx, dbOpts)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user", err)
	}

	err = r.sessionManager.DeleteUserSessions(user.ID)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user sessions", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getUserByPrincipal retrieves a user by the principal's user ID
func (r *Router) getUserByPrincipal(ctx context.Context, principal types.Principal) (*models.User, error) {
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: principal.UserID})
	return r.appDao.GetUser(ctx, dbOpts)
}
