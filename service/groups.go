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

const maxGroupNameLength = 64

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroupRequest represents a group create request
type CreateGroupRequest struct {
	Name string `json:"name"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupResponse represents a group response
type GroupResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"createdBy,omitempty"`
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

	return groupResponseBuilder(group, true), nil
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

	groups, err := g.dao.ListGroups(ctx, dao.NewOptions().WithPagination(page).WithWhere(where))
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return []*GroupResponse{}, nil
	}

	return groupsResponseBuilder(groups, false), nil
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

	return groupResponseBuilder(group, false), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupsResponseBuilder maps group models to API response slices
func groupsResponseBuilder(groups []*models.Group, includeCreatedBy bool) []*GroupResponse {
	out := make([]*GroupResponse, len(groups))
	for i, group := range groups {
		out[i] = groupResponseBuilder(group, includeCreatedBy)
	}

	return out
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupResponseBuilder maps a group model to a single API response
func groupResponseBuilder(group *models.Group, includeCreatedBy bool) *GroupResponse {
	resp := &GroupResponse{
		ID:   group.ID,
		Name: group.Name,
	}

	if includeCreatedBy && group.CreatedBy != "" {
		resp.CreatedBy = group.CreatedBy
	}

	return resp
}
