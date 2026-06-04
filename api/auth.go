package api

import (
	"github.com/geerew/friendle/service"
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

	authGroup.Get("/me", r.requireAccess(accessSiteUser), r.getMe)
	authGroup.Put("/me", r.requireAccess(accessSiteUser), r.updateMe)
	authGroup.Delete("/me", r.requireAccess(accessSiteUser), r.deleteMe)
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

	if r.app.IsBootstrapped() {
		return errorResponse(c, fiber.StatusForbidden, "Application is already bootstrapped", nil)
	}

	if err := auth.ValidateBootstrapToken(token, r.app.Config.DataDir, r.app.FS); err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Invalid or expired bootstrap token", nil)
	}

	req := &service.RegisterRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	ctx := c.UserContext()
	user, err := r.appSvc.Auth.Register(ctx, *req, types.SiteRoleAdmin)
	if err != nil {
		return serviceError(c, err)
	}

	if err := r.sessionManager.SetSession(c, user.ID, user.SiteRole); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error setting session", err)
	}

	r.app.SetBootstrapped()
	auth.DeleteBootstrapToken(r.app.Config.DataDir, r.app.FS)

	return c.Status(fiber.StatusCreated).JSON(user)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// register creates a user account and starts a session
func (r *Router) register(c *fiber.Ctx) error {
	if r.app.IsBootstrapped() && !r.app.Config.EnableSignup {
		return errorResponse(c, fiber.StatusForbidden, "Sign-up is disabled", nil)
	}

	req := &service.RegisterRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	ctx := c.UserContext()
	user, err := r.appSvc.Auth.Register(ctx, *req, types.SiteRoleUser)
	if err != nil {
		return serviceError(c, err)
	}

	if err := r.sessionManager.SetSession(c, user.ID, user.SiteRole); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error setting session", err)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// login authenticates a user and starts a session
func (r *Router) login(c *fiber.Ctx) error {
	req := &service.LoginRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	ctx := c.UserContext()
	user, err := r.appSvc.Auth.Login(ctx, *req)
	if err != nil {
		return serviceError(c, err)
	}

	if err := r.sessionManager.SetSession(c, user.ID, user.SiteRole); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error setting session", err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
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
	_, ctx := principalAndCtx(c)

	user, err := r.appSvc.Auth.GetMe(ctx)
	if err != nil {
		return serviceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// updateMe updates the authenticated user's profile or password
func (r *Router) updateMe(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	req := &service.UpdateMeRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	user, err := r.appSvc.Auth.UpdateMe(ctx, *req)
	if err != nil {
		return serviceError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteMe deletes the authenticated user's account
func (r *Router) deleteMe(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	req := &service.DeleteMeRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	if err := r.appSvc.Auth.DeleteMe(ctx, *req); err != nil {
		return serviceError(c, err)
	}

	if err := r.sessionManager.DeleteUserSessions(principal.UserID); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting user sessions", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
