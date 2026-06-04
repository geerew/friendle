package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateJoinRequest inserts a group join record
func (dao *DAO) CreateJoinRequest(ctx context.Context, r *models.GroupJoinRequest) error {
	if r.ID == "" {
		r.RefreshId()
	}
	r.RefreshCreatedAt()
	r.RefreshUpdatedAt()

	if r.Status == "" {
		r.Status = types.JoinPending
	}

	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithData(map[string]interface{}{
			models.BASE_ID:               r.ID,
			models.JOIN_REQUEST_GROUP_ID: r.GroupID,
			models.JOIN_REQUEST_USER_ID:  r.UserID,
			models.JOIN_REQUEST_STATUS:   r.Status,
			models.BASE_CREATED_AT:       r.CreatedAt,
			models.BASE_UPDATED_AT:       r.UpdatedAt,
		})

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetJoinRequest returns a group join record
func (dao *DAO) GetJoinRequest(ctx context.Context, dbOpts *Options) (*models.GroupJoinRequest, error) {
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.GroupJoinRequestColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.GroupJoinRequest](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListJoinRequests returns group join records
func (dao *DAO) ListJoinRequests(ctx context.Context, dbOpts *Options) ([]*models.GroupJoinRequest, error) {
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.GroupJoinRequestColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.GroupJoinRequest](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateJoinRequest updates a group join record
func (dao *DAO) UpdateJoinRequest(ctx context.Context, r *models.GroupJoinRequest) error {
	if r.ID == "" {
		return utils.ErrId
	}

	r.RefreshUpdatedAt()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: r.ID})
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithData(map[string]interface{}{
			models.JOIN_REQUEST_STATUS: r.Status,
			models.BASE_UPDATED_AT:     r.UpdatedAt,
		}).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteJoinRequests deletes group join records
//
// Errors when a where clause is not provided
func (dao *DAO) DeleteJoinRequests(ctx context.Context, dbOpts *Options) error {
	if dbOpts == nil || dbOpts.Where == nil {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).SetDbOpts(dbOpts)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)
	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListPendingJoinGroupIDsForUser returns group IDs with a pending join request for the user
func (dao *DAO) ListPendingJoinGroupIDsForUser(ctx context.Context, userID string, groupIDs []string) ([]string, error) {
	if len(groupIDs) == 0 {
		return []string{}, nil
	}

	dbOpts := NewOptions().WithWhere(squirrel.Eq{
		models.JOIN_REQUEST_USER_ID:  userID,
		models.JOIN_REQUEST_GROUP_ID: groupIDs,
		models.JOIN_REQUEST_STATUS:   types.JoinPending,
	})

	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.JOIN_REQUEST_TABLE_GROUP_ID).
		SetDbOpts(dbOpts)

	return pluck[string](ctx, dao, *builderOpts)
}
