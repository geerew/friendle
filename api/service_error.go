package api

import (
	"errors"

	"github.com/geerew/friendle/service"
	"github.com/geerew/friendle/utils"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type serviceErrorMapping struct {
	err     error
	status  int
	message string
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// serviceErrorMappings maps service-layer sentinel errors to HTTP responses. Add a row when
// introducing a new service.Err* value
var serviceErrorMappings = []serviceErrorMapping{
	{err: utils.ErrPrincipal, status: fiber.StatusUnauthorized, message: "Unauthorized"},

	{err: service.ErrNotSiteAdmin, status: fiber.StatusForbidden, message: "Forbidden"},
	{err: service.ErrUserNotFound, status: fiber.StatusNotFound, message: "User not found"},
	{err: service.ErrUsernameTaken, status: fiber.StatusBadRequest, message: "Username already exists"},
	{err: service.ErrNoUpdateData, status: fiber.StatusBadRequest, message: "No data to update"},
	{err: service.ErrCredentialsRequired, status: fiber.StatusBadRequest, message: "Username and/or password cannot be empty"},
	{err: service.ErrInvalidCredentials, status: fiber.StatusUnauthorized, message: "Invalid username and/or password"},
	{err: service.ErrInvalidCurrentPassword, status: fiber.StatusBadRequest, message: "Invalid current password"},
	{err: service.ErrInvalidPassword, status: fiber.StatusBadRequest, message: "Invalid password"},
	{err: service.ErrLastAdmin, status: fiber.StatusBadRequest, message: "Unable to delete the last admin user"},
	{err: service.ErrPasswordTooShort, status: fiber.StatusBadRequest, message: "password must be at least 8 characters"},
	{err: service.ErrPasswordTooLong, status: fiber.StatusBadRequest, message: "password must be no more than 128 characters"},

	{err: service.ErrGroupNameRequired, status: fiber.StatusBadRequest, message: "Group name is required"},
	{err: service.ErrGroupNameTaken, status: fiber.StatusBadRequest, message: "Group name already exists"},
	{err: service.ErrGroupNameTooLong, status: fiber.StatusBadRequest, message: "Group name is too long"},
	{err: service.ErrGroupNotFound, status: fiber.StatusNotFound, message: "Group not found"},
	{err: service.ErrGroupSearchQueryRequired, status: fiber.StatusBadRequest, message: "Search query is required"},

	{err: utils.ErrApiQueryParse, status: fiber.StatusBadRequest, message: "Invalid query"},
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// serviceError maps a service-layer error to an HTTP JSON error response
func serviceError(c *fiber.Ctx, err error) error {
	for _, mapping := range serviceErrorMappings {
		if errors.Is(err, mapping.err) {
			return errorResponse(c, mapping.status, mapping.message, nil)
		}
	}

	return errorResponse(c, fiber.StatusInternalServerError, "Request failed", err)
}
