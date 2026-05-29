package api

import (
	"context"
	"fmt"

	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// errorResponse writes a JSON error response and stores details for request logging
func errorResponse(c *fiber.Ctx, status int, message string, err error) error {
	resp := fiber.Map{"message": message}
	if err != nil {
		resp["error"] = err.Error()
	}
	if status >= 400 {
		c.Locals("api_error_message", message)
		if err != nil {
			c.Locals("api_error_detail", err.Error())
		}
	}

	return c.Status(status).JSON(resp)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// validatePassword checks password length constraints
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 128 {
		return fmt.Errorf("password must be no more than 128 characters")
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// principalCtx returns the authenticated principal and a context carrying it
func principalCtx(c *fiber.Ctx) (types.Principal, context.Context, error) {
	principal, ok := c.Locals(types.PrincipalContextKey).(types.Principal)
	if !ok {
		return types.Principal{}, nil, fmt.Errorf("missing principal")
	}
	ctx := c.UserContext()
	ctx = context.WithValue(ctx, types.PrincipalContextKey, principal)

	return principal, ctx, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// paginationFromCtx builds pagination options from query parameters
func paginationFromCtx(c *fiber.Ctx) *pagination.Pagination {
	return pagination.New(
		pagination.ParsePage(c.Query(pagination.PageQueryParam)),
		pagination.ParsePerPage(c.Query(pagination.PerPageQueryParam)),
	)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// protectedRoute requires a site admin principal
func protectedRoute(c *fiber.Ctx) error {
	principal, _, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Missing principal", nil)
	}
	if principal.SiteRole != types.SiteRoleAdmin {
		return errorResponse(c, fiber.StatusForbidden, "User is not an admin", nil)
	}

	return c.Next()
}
