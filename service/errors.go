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
	ErrGroupAlreadyMember     = errors.New("already a group member")
	ErrGroupNotMember         = errors.New("not a group member")
	ErrGroupNotAdmin          = errors.New("not a group admin")
	ErrGroupJoinRequestPending  = errors.New("join request already pending")
	ErrGroupJoinRequestRejected = errors.New("join request was rejected")
	ErrGroupJoinRequestNotFound = errors.New("join request not found")
	ErrGroupMemberNotFound      = errors.New("group member not found")
	ErrGroupMemberSelf          = errors.New("cannot modify your own group membership")
	ErrGroupLastAdmin           = errors.New("unable to remove the last group admin")
)
