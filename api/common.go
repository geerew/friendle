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

	if writeErr := c.Status(status).JSON(resp); writeErr != nil {
		return writeErr
	}

	if status >= 400 {
		return fiber.NewError(status, message)
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

// principalAndCtx returns the session principal and a context carrying it. The handler must be
// registered behind require middleware or global auth that has already validated the session
func principalAndCtx(c *fiber.Ctx) (types.Principal, context.Context) {
	principal, ctx, _ := principalCtx(c)

	return principal, ctx
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// paginationFromCtx builds pagination options from query parameters
func paginationFromCtx(c *fiber.Ctx) *pagination.Pagination {
	return pagination.New(
		pagination.ParsePage(c.Query(pagination.PageQueryParam)),
		pagination.ParsePerPage(c.Query(pagination.PerPageQueryParam)),
	)
}
