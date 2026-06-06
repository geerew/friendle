package service

import "errors"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var (
	ErrNotSiteAdmin           = errors.New("not a site admin")
	ErrUserNotFound           = errors.New("user not found")
	ErrUsernameTaken          = errors.New("username already exists")
	ErrNoUpdateData           = errors.New("no data to update")
	ErrCredentialsRequired    = errors.New("username and password required")
	ErrPasswordTooShort       = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong        = errors.New("password must be no more than 128 characters")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrInvalidCurrentPassword = errors.New("invalid current password")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrLastAdmin              = errors.New("unable to delete the last admin user")
	ErrGroupNameRequired      = errors.New("group name is required")
	ErrGroupNameTaken         = errors.New("group name already exists")
	ErrGroupNameTooLong       = errors.New("group name is too long")
	ErrGroupNotFound          = errors.New("group not found")
	ErrGroupSearchQueryRequired = errors.New("group search query is required")
)
