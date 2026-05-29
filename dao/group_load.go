package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupLoad selects which relations to populate on a group
type GroupLoad struct {
	Members      bool
	JoinRequests bool
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroupLoaded returns a group with optional relations loaded
func (dao *DAO) GetGroupLoaded(ctx context.Context, groupID string, load GroupLoad) (*models.Group, error) {
	group, err := dao.GetGroup(ctx, groupID)
	if err != nil || group == nil {
		return group, err
	}

	if err := dao.LoadGroup(ctx, group, load); err != nil {
		return nil, err
	}

	return group, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// LoadGroup populates relation fields on an existing group row
func (dao *DAO) LoadGroup(ctx context.Context, group *models.Group, load GroupLoad) error {
	if load.Members {
		members, err := dao.ListGroupMembers(ctx, group.ID)
		if err != nil {
			return err
		}
		for _, m := range members {
			u, _ := dao.GetUser(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: m.UserID}))
			m.User = u
		}
		group.Members = members
	}

	if load.JoinRequests {
		requests, err := dao.ListPendingJoinRequests(ctx, group.ID)
		if err != nil {
			return err
		}
		group.JoinRequests = requests
	}

	return nil
}
