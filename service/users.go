package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/queryparser"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var (
	userListApiAllowedFilters = []string{"name", "role", "siteRole"}

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

// CreateUserRequest represents a user create request
type CreateUserRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	SiteRole    string `json:"siteRole"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	SiteRole    string `json:"siteRole"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Users orchestrates site-admin user reads and writes
type Users struct {
	deps
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// newUsers creates a Users service
func newUsers(d deps) *Users {
	return &Users{deps: d}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListUsers returns paginated users
func (u *Users) ListUsers(ctx context.Context, page *pagination.Pagination) ([]*UserResponse, error) {
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

// CreateUser creates a user from the site admin API
func (u *Users) CreateUser(ctx context.Context, req CreateUserRequest) error {
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

// UpdateUser updates a user and reports whether the site role changed
func (u *Users) UpdateUser(ctx context.Context, userID string, req UpdateUserRequest) (*UserResponse, bool, error) {
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

// DeleteUser deletes a user
func (u *Users) DeleteUser(ctx context.Context, userID string) error {
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_ID: userID})
	return u.dao.DeleteUsers(ctx, dbOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// usersResponseBuilder maps user models to API responses
func usersResponseBuilder(users []*models.User) []*UserResponse {
	responses := make([]*UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, userResponseBuilder(user))
	}

	return responses
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userResponseBuilder maps a user model to an API response
func userResponseBuilder(user *models.User) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		SiteRole:    user.SiteRole,
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// usersWhereFromApiQuery parses a user list API query into a WHERE clause
func usersWhereFromApiQuery(apiQuery string) (squirrel.Sqlizer, error) {
	parsed, err := queryparser.Parse(apiQuery, userListApiAllowedFilters)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", utils.ErrApiQueryParse, err)
	}

	if parsed == nil {
		return nil, nil
	}

	return usersWhereBuilder(parsed.Expr), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// usersWhereBuilder builds a squirrel WHERE expression from a queryparser.QueryExpr
func usersWhereBuilder(expr queryparser.QueryExpr) squirrel.Sqlizer {
	switch node := expr.(type) {
	case *queryparser.FilterExpr:
		switch node.Key {
		case "name":
			return squirrel.Or{
				squirrel.Like{"LOWER(" + models.USER_TABLE_USERNAME + ")": "%" + node.Value + "%"},
				squirrel.Like{"LOWER(" + models.USER_TABLE_DISPLAY_NAME + ")": "%" + node.Value + "%"},
			}
		case "role", "siteRole":
			val := node.Value
			if val == "admin" {
				val = string(types.SiteRoleAdmin)
			} else if val == "user" {
				val = string(types.SiteRoleUser)
			}

			return squirrel.Eq{models.USER_TABLE_SITE_ROLE: val}
		default:
			return nil
		}
	case *queryparser.AndExpr:
		var andSlice []squirrel.Sqlizer
		for _, child := range node.Children {
			andSlice = append(andSlice, usersWhereBuilder(child))
		}

		return squirrel.And(andSlice)
	case *queryparser.OrExpr:
		var orSlice []squirrel.Sqlizer
		for _, child := range node.Children {
			orSlice = append(orSlice, usersWhereBuilder(child))
		}

		return squirrel.Or(orSlice)
	default:
		return nil
	}
}
