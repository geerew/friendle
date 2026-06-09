package dao

import (
	"context"

	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CountGroupJoinRequests returns the number of group join request records
func (dao *DAO) CountGroupJoinRequests(ctx context.Context, dbOpts *Options) (int, error) {
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).SetDbOpts(dbOpts)

	return countGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroupJoinRequest inserts a group join request record
func (dao *DAO) CreateGroupJoinRequest(ctx context.Context, request *models.GroupJoinRequest) error {
	if request == nil {
		return utils.ErrNilPtr
	}

	if request.GroupID == "" {
		return utils.ErrGroupId
	}

	if request.UserID == "" {
		return utils.ErrUserId
	}

	if !request.Status.IsValid() {
		request.Status = types.JoinPending
	}

	if request.ID == "" {
		request.RefreshId()
	}

	request.RefreshCreatedAt()
	request.RefreshUpdatedAt()

	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithData(
			map[string]interface{}{
				models.BASE_ID:               request.ID,
				models.JOIN_REQUEST_GROUP_ID: request.GroupID,
				models.JOIN_REQUEST_USER_ID:  request.UserID,
				models.JOIN_REQUEST_STATUS:   request.Status,
				models.BASE_CREATED_AT:       request.CreatedAt,
				models.BASE_UPDATED_AT:       request.UpdatedAt,
			},
		)

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroupJoinRequests returns group join request records
func (dao *DAO) ListGroupJoinRequests(ctx context.Context, dbOpts *Options) ([]*models.GroupJoinRequest, error) {
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.GroupJoinRequestColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.GroupJoinRequest](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroupJoinRequest returns a group join request record
func (dao *DAO) GetGroupJoinRequest(ctx context.Context, dbOpts *Options) (*models.GroupJoinRequest, error) {
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.GroupJoinRequestColumns()...).
		SetDbOpts(dbOpts)

	return getGeneric[models.GroupJoinRequest](ctx, dao, *builderOpts)
}
