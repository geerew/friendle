package api

import (
	"context"
	"fmt"

	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

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

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if len(password) > 128 {
		return fmt.Errorf("password must be no more than 128 characters")
	}
	return nil
}

func principalCtx(c *fiber.Ctx) (types.Principal, context.Context, error) {
	principal, ok := c.Locals(types.PrincipalContextKey).(types.Principal)
	if !ok {
		return types.Principal{}, nil, fmt.Errorf("missing principal")
	}
	ctx := c.UserContext()
	ctx = context.WithValue(ctx, types.PrincipalContextKey, principal)
	return principal, ctx, nil
}

func paginationFromCtx(c *fiber.Ctx) *pagination.Pagination {
	return pagination.New(
		pagination.ParsePage(c.Query(pagination.PageQueryParam)),
		pagination.ParsePerPage(c.Query(pagination.PerPageQueryParam)),
	)
}

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
