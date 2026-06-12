package dao

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupJoinRequestColumns defines the columns to select
var groupJoinRequestColumns = []string{
	fmt.Sprintf("%s AS %s", models.JOIN_REQUEST_TABLE_ID, models.BASE_ID),
	fmt.Sprintf("%s AS %s", models.JOIN_REQUEST_TABLE_CREATED_AT, models.BASE_CREATED_AT),
	fmt.Sprintf("%s AS %s", models.JOIN_REQUEST_TABLE_UPDATED_AT, models.BASE_UPDATED_AT),
	fmt.Sprintf("%s AS %s", models.JOIN_REQUEST_TABLE_GROUP_ID, models.JOIN_REQUEST_GROUP_ID),
	fmt.Sprintf("%s AS %s", models.JOIN_REQUEST_TABLE_USER_ID, models.JOIN_REQUEST_USER_ID),
	fmt.Sprintf("%s AS %s", models.JOIN_REQUEST_TABLE_STATUS, models.JOIN_REQUEST_STATUS),
	fmt.Sprintf("%s AS %s", models.USER_TABLE_DISPLAY_NAME, models.USER_DISPLAY_NAME),
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var groupJoinRequestJoins = []join{
	{
		Type:      joinTypeInner,
		Table:     models.USER_TABLE,
		Condition: models.JOIN_REQUEST_TABLE_USER_ID + " = " + models.USER_TABLE_ID,
	},
}

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

// ListGroupJoinRequests returns group join request records with user display names
func (dao *DAO) ListGroupJoinRequests(ctx context.Context, dbOpts *Options) ([]*models.GroupJoinRequest, error) {
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(groupJoinRequestColumns...).
		WithJoins(groupJoinRequestJoins...).
		SetDbOpts(dbOpts)

	return listGeneric[models.GroupJoinRequest](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroupJoinRequest returns a group join request record
func (dao *DAO) GetGroupJoinRequest(ctx context.Context, dbOpts *Options) (*models.GroupJoinRequest, error) {
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(groupJoinRequestColumns...).
		WithJoins(groupJoinRequestJoins...).
		SetDbOpts(dbOpts)

	return getGeneric[models.GroupJoinRequest](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateGroupJoinRequest updates a group join request record
func (dao *DAO) UpdateGroupJoinRequest(ctx context.Context, request *models.GroupJoinRequest) error {
	if request == nil {
		return utils.ErrNilPtr
	}

	if request.ID == "" {
		return utils.ErrId
	}

	request.RefreshUpdatedAt()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: request.ID})

	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithData(
			map[string]interface{}{
				models.JOIN_REQUEST_STATUS: request.Status,
				models.BASE_UPDATED_AT:     request.UpdatedAt,
			},
		).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteGroupJoinRequests deletes group join request records
//
// Errors when a where clause is not provided
func (dao *DAO) DeleteGroupJoinRequests(ctx context.Context, dbOpts *Options) error {
	if dbOpts == nil || dbOpts.Where == nil {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).SetDbOpts(dbOpts)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)

	return err
}
