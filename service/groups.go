package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/geerew/friendle/models"
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
	CreatedBy string `json:"createdBy"`
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

// Create creates a group and adds the caller as group admin
func (g *Groups) Create(ctx context.Context, req CreateGroupRequest) (*GroupResponse, error) {
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

	return groupResponseBuilder(group), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupResponseBuilder maps a group model to a response
func groupResponseBuilder(group *models.Group) *GroupResponse {
	return &GroupResponse{
		ID:        group.ID,
		Name:      group.Name,
		CreatedBy: group.CreatedBy,
	}
}
