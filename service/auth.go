package service

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateMeRequest represents an authenticated user's profile update request
type UpdateMeRequest struct {
	DisplayName     string `json:"displayName"`
	CurrentPassword string `json:"currentPassword"`
	Password        string `json:"password"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteMeRequest represents an authenticated user's account deletion request
type DeleteMeRequest struct {
	CurrentPassword string `json:"currentPassword"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Auth orchestrates registration, login, and self-service profile operations
type Auth struct {
	deps
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// newAuth creates an Auth service
func newAuth(d deps) *Auth {
	return &Auth{deps: d}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Register creates a user account with the given site role
func (a *Auth) Register(ctx context.Context, req RegisterRequest, siteRole types.SiteRole) (*UserResponse, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)
	displayName := strings.TrimSpace(req.DisplayName)

	if username == "" || password == "" {
		return nil, ErrCredentialsRequired
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	if displayName == "" {
		displayName = username
	}

	passwordHash, err := auth.GeneratePassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: passwordHash,
		SiteRole:     siteRole,
	}

	if err := a.dao.CreateUser(ctx, user); err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return nil, ErrUsernameTaken
		}

		return nil, err
	}

	return userResponseBuilder(user), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Login authenticates a user by username and password
func (a *Auth) Login(ctx context.Context, req LoginRequest) (*UserResponse, error) {
	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)

	if username == "" || password == "" {
		return nil, ErrCredentialsRequired
	}

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: username})
	user, err := a.dao.GetUser(ctx, dbOpts)
	if err != nil {
		return nil, err
	}

	if user == nil || !auth.ComparePassword(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}

	return userResponseBuilder(user), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetMe returns the authenticated user's profile
func (a *Auth) GetMe(ctx context.Context, userID string) (*UserResponse, error) {
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: userID})
	user, err := a.dao.GetUser(ctx, dbOpts)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return userResponseBuilder(user), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateMe updates the authenticated user's profile or password
func (a *Auth) UpdateMe(ctx context.Context, userID string, req UpdateMeRequest) (*UserResponse, error) {
	displayName := strings.TrimSpace(req.DisplayName)
	password := strings.TrimSpace(req.Password)
	currentPassword := strings.TrimSpace(req.CurrentPassword)

	if displayName == "" && password == "" {
		return nil, ErrNoUpdateData
	}

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: userID})
	user, err := a.dao.GetUser(ctx, dbOpts)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if displayName != "" {
		user.DisplayName = displayName
	}

	if password != "" {
		if !auth.ComparePassword(user.PasswordHash, currentPassword) {
			return nil, ErrInvalidCurrentPassword
		}

		if err := validatePassword(password); err != nil {
			return nil, err
		}

		passwordHash, err := auth.GeneratePassword(password)
		if err != nil {
			return nil, err
		}

		user.PasswordHash = passwordHash
	}

	if err := a.dao.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return userResponseBuilder(user), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteMe deletes the authenticated user's account
func (a *Auth) DeleteMe(ctx context.Context, userID string, req DeleteMeRequest) error {
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: userID})
	user, err := a.dao.GetUser(ctx, dbOpts)
	if err != nil {
		return err
	}

	if user == nil {
		return ErrUserNotFound
	}

	if !auth.ComparePassword(user.PasswordHash, strings.TrimSpace(req.CurrentPassword)) {
		return ErrInvalidPassword
	}

	// Ensure there is at least one admin user remaining
	if user.SiteRole == types.SiteRoleAdmin {
		dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_SITE_ROLE: types.SiteRoleAdmin})
		adminCount, err := a.dao.CountUsers(ctx, dbOpts)
		if err != nil {
			return err
		}

		if adminCount == 1 {
			return ErrLastAdmin
		}
	}

	dbOpts = dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: userID})
	return a.dao.DeleteUsers(ctx, dbOpts)
}
