package service

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var (
	defaultUsersListOrderBy = []string{models.USER_TABLE_CREATED_AT + " desc"}
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UserResponse represents a user response
type UserResponse struct {
	ID          string         `json:"id"`
	Username    string         `json:"username"`
	DisplayName string         `json:"displayName"`
	SiteRole    types.SiteRole `json:"siteRole"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UserCreateRequest represents a user create request
type UserCreateRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	SiteRole    string `json:"siteRole"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UserUpdateRequest represents a user update request
type UserUpdateRequest struct {
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	SiteRole    string `json:"siteRole"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Users represents the user service
type Users struct {
	deps
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// newUsers creates a user service
func newUsers(d deps) *Users {
	return &Users{deps: d}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// List returns paginated slice of users
func (u *Users) List(ctx context.Context, page *pagination.Pagination) ([]*UserResponse, error) {
	users, err := u.dao.ListUsers(ctx, dao.NewOptions().
		WithPagination(page).
		WithOrderBy(defaultUsersListOrderBy...))
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return []*UserResponse{}, nil
	}

	return usersResponseBuilder(users), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Create creates a user
func (u *Users) Create(ctx context.Context, req UserCreateRequest) error {
	username := strings.TrimSpace(req.Username)
	displayName := strings.TrimSpace(req.DisplayName)
	password := strings.TrimSpace(req.Password)
	role := strings.TrimSpace(req.SiteRole)

	if displayName == "" {
		displayName = username
	}

	if username == "" || password == "" {
		return ErrCredentialsRequired
	}

	if err := validatePassword(password); err != nil {
		return err
	}

	passwordHash, err := auth.GeneratePassword(password)
	if err != nil {
		return err
	}

	if role == "" {
		role = types.SiteRoleUser.String()
	}

	user := &models.User{
		Username:     username,
		DisplayName:  displayName,
		PasswordHash: passwordHash,
		SiteRole:     types.NewSiteRole(role),
	}

	// Create the user
	if err := u.dao.CreateUser(ctx, user); err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return ErrUsernameTaken
		}

		return err
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Update updates a user and reports whether the site role changed
func (u *Users) Update(ctx context.Context, userID string, req UserUpdateRequest) (*UserResponse, bool, error) {
	displayName := strings.TrimSpace(req.DisplayName)
	password := strings.TrimSpace(req.Password)
	role := strings.TrimSpace(req.SiteRole)

	if displayName == "" && password == "" && role == "" {
		return nil, false, ErrNoUpdateData
	}

	// Get the existing user
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: userID})
	user, err := u.dao.GetUser(ctx, dbOpts)
	if err != nil {
		return nil, false, err
	}

	if user == nil {
		return nil, false, ErrUserNotFound
	}

	if displayName != "" {
		user.DisplayName = displayName
	}

	// Validate the password then hash it
	if password != "" {
		if err := validatePassword(password); err != nil {
			return nil, false, err
		}

		passwordHash, err := auth.GeneratePassword(password)
		if err != nil {
			return nil, false, err
		}

		user.PasswordHash = passwordHash
	}

	roleChanged := false
	if role != "" && user.SiteRole.String() != role {
		user.SiteRole = types.NewSiteRole(role)
		roleChanged = true
	}

	// Update the user
	if err := u.dao.UpdateUser(ctx, user); err != nil {
		return nil, false, err
	}

	return userResponseBuilder(user), roleChanged, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Delete deletes a user
func (u *Users) Delete(ctx context.Context, userID string) error {
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: userID})
	return u.dao.DeleteUsers(ctx, dbOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Response builders
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// usersResponseBuilder builds a slice of UserResponse from a slice of user models
func usersResponseBuilder(users []*models.User) []*UserResponse {
	responses := make([]*UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, userResponseBuilder(user))
	}

	return responses
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userResponseBuilder builds a UserResponse from a user model
func userResponseBuilder(user *models.User) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		SiteRole:    user.SiteRole,
	}
}
