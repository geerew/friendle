package utils

import "errors"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var (
	// Generic
	ErrNilPtr = errors.New("nil pointer")

	// API
	ErrApiQueryParse = errors.New("list query parse error")

	// DB
	ErrWhere           = errors.New("where clause cannot be empty")
	ErrPrincipal       = errors.New("principal not found in context")
	ErrGroupMembership = errors.New("group membership not found in context")

	// Model
	ErrId           = errors.New("id cannot be empty")
	ErrUsername     = errors.New("username cannot be empty")
	ErrUserPassword = errors.New("user password cannot be empty")
	ErrUserId       = errors.New("user id cannot be empty")
)
