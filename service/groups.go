package service

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/queryparser"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var (
	groupListApiAllowedFilters = []string{"name"}

	defaultGroupsListOrderBy = []string{models.GROUP_TABLE_NAME + " asc"}
)

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

// SearchGroups returns paginated groups matching an API query
func (g *Groups) SearchGroups(ctx context.Context, page *pagination.Pagination, apiQuery string) ([]*GroupResponse, error) {
	if _, err := principalFromContext(ctx); err != nil {
		return nil, err
	}

	if strings.TrimSpace(apiQuery) == "" {
		return nil, ErrGroupSearchQueryRequired
	}

	where, err := groupsWhereFromApiQuery(apiQuery)
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupsWhereFromApiQuery parses a group list API query into a WHERE clause
func groupsWhereFromApiQuery(apiQuery string) (squirrel.Sqlizer, error) {
	parsed, err := queryparser.Parse(apiQuery, groupListApiAllowedFilters)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", utils.ErrApiQueryParse, err)
	}

	if parsed == nil {
		return nil, nil
	}

	return groupsWhereBuilder(parsed.Expr), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupsWhereBuilder builds a squirrel WHERE expression from a queryparser.QueryExpr
func groupsWhereBuilder(expr queryparser.QueryExpr) squirrel.Sqlizer {
	switch node := expr.(type) {
	case *queryparser.FilterExpr:
		switch node.Key {
		case "name":
			pattern := strings.ToLower(strings.TrimSpace(node.Value))
			pattern = strings.ReplaceAll(pattern, "*", "%")

			if !strings.Contains(pattern, "%") {
				pattern = "%" + pattern + "%"
			}

			return squirrel.Like{"LOWER(" + models.GROUP_TABLE_NAME + ")": pattern}
		default:
			return nil
		}
	case *queryparser.AndExpr:
		var andSlice []squirrel.Sqlizer
		for _, child := range node.Children {
			andSlice = append(andSlice, groupsWhereBuilder(child))
		}

		return squirrel.And(andSlice)
	case *queryparser.OrExpr:
		var orSlice []squirrel.Sqlizer
		for _, child := range node.Children {
			orSlice = append(orSlice, groupsWhereBuilder(child))
		}

		return squirrel.Or(orSlice)
	default:
		return nil
	}
}
