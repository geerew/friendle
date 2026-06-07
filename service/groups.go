package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var defaultGroupsListOrderBy = []string{models.GROUP_TABLE_NAME + " asc"}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const maxGroupNameLength = 64

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroupRequest represents a group create request
type CreateGroupRequest struct {
	Name string `json:"name"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupResponse represents a group response
type GroupResponse struct {
	ID          string         `json:"id"`
	CreatedAt   types.DateTime `json:"createdAt"`
	UpdatedAt   types.DateTime `json:"updatedAt"`
	Name        string         `json:"name"`
	CreatedBy   string         `json:"createdBy,omitempty"`
	MemberCount int            `json:"memberCount"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Groups orchestrates friend group operations
type Groups struct {
	deps
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// newGroups creates a Groups service
func newGroups(d deps) *Groups {
	return &Groups{deps: d}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroup creates a group and adds the caller as group admin
func (g *Groups) CreateGroup(ctx context.Context, req CreateGroupRequest) (*GroupResponse, error) {
	principal, err := principalFromContext(ctx)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrGroupNameRequired
	}

	if utf8.RuneCountInString(name) > maxGroupNameLength {
		return nil, ErrGroupNameTooLong
	}

	group := &models.Group{
		Name:      name,
		CreatedBy: principal.UserID,
	}

	err = g.dao.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := g.dao.CreateGroup(txCtx, group); err != nil {
			if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
				return ErrGroupNameTaken
			}

			return err
		}

		member := &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.UserID,
			GroupRole: types.GroupRoleAdmin,
		}

		return g.dao.CreateGroupMember(txCtx, member)
	})
	if err != nil {
		return nil, err
	}

	group.MemberCount = 1

	return groupResponseBuilder(group), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroups returns paginated groups
func (g *Groups) ListGroups(ctx context.Context, page *pagination.Pagination) ([]*GroupResponse, error) {
	if _, err := principalFromContext(ctx); err != nil {
		return nil, err
	}

	groups, err := g.dao.ListGroups(ctx, dao.NewOptions().
		WithPagination(page).
		WithOrderBy(defaultGroupsListOrderBy...))
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return []*GroupResponse{}, nil
	}

	return groupsResponseBuilder(groups), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListSelfGroups returns paginated groups the authenticated user belongs to
func (g *Groups) ListSelfGroups(ctx context.Context, page *pagination.Pagination) ([]*GroupResponse, error) {
	principal, err := principalFromContext(ctx)
	if err != nil {
		return nil, err
	}

	where, err := dao.MemberGroupsWhere(principal.UserID)
	if err != nil {
		return nil, err
	}

	groups, err := g.dao.ListGroups(ctx, dao.NewOptions().
		WithPagination(page).
		WithOrderBy(defaultGroupsListOrderBy...).
		WithWhere(where))
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return []*GroupResponse{}, nil
	}

	return groupsResponseBuilder(groups), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SearchGroups returns paginated groups whose names contain the search term
func (g *Groups) SearchGroups(ctx context.Context, page *pagination.Pagination, name string) ([]*GroupResponse, error) {
	if _, err := principalFromContext(ctx); err != nil {
		return nil, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrGroupSearchQueryRequired
	}

	lowerName := strings.ToLower(name)
	whereClause := squirrel.Like{"LOWER(" + models.GROUP_TABLE_NAME + ")": "%" + lowerName + "%"}
	orderByClause := squirrel.Expr("CASE WHEN "+models.GROUP_TABLE_NAME+" LIKE ? THEN 0 ELSE 1 END, "+models.GROUP_TABLE_NAME+" ASC", lowerName+"%")

	dbOpts := dao.NewOptions().
		WithPagination(page).
		WithOrderByClause(orderByClause).
		WithWhere(whereClause)

	groups, err := g.dao.ListGroups(ctx, dbOpts)
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return []*GroupResponse{}, nil
	}

	return groupsResponseBuilder(groups), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroup returns a group the authenticated user belongs to
func (g *Groups) GetGroup(ctx context.Context, groupID string) (*GroupResponse, error) {
	principal, err := principalFromContext(ctx)
	if err != nil {
		return nil, err
	}

	memberWhere, err := dao.MemberGroupsWhere(principal.UserID)
	if err != nil {
		return nil, err
	}

	where := squirrel.And{
		squirrel.Eq{models.GROUP_TABLE_ID: groupID},
		memberWhere,
	}

	group, err := g.dao.GetGroup(ctx, dao.NewOptions().WithWhere(where))
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, ErrGroupNotFound
	}

	return groupResponseBuilder(group), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteGroup deletes a group and its associated data
func (g *Groups) DeleteGroup(ctx context.Context, groupID string) error {
	principal, err := principalFromContext(ctx)
	if err != nil {
		return err
	}

	if principal.SiteRole != types.SiteRoleAdmin {
		return ErrNotSiteAdmin
	}

	group, err := g.dao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
	if err != nil {
		return err
	}

	if group == nil {
		return ErrGroupNotFound
	}

	return g.dao.DeleteGroups(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupsResponseBuilder maps group models to API response slices
func groupsResponseBuilder(groups []*models.Group) []*GroupResponse {
	out := make([]*GroupResponse, len(groups))
	for i, group := range groups {
		out[i] = groupResponseBuilder(group)
	}

	return out
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupResponseBuilder maps a group model to a single API response
func groupResponseBuilder(group *models.Group) *GroupResponse {
	return &GroupResponse{
		ID:          group.ID,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
		Name:        group.Name,
		CreatedBy:   group.CreatedBy,
		MemberCount: group.MemberCount,
	}
}
