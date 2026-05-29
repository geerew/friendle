package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateJoinRequest inserts a group join request
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

// GetJoinRequest returns a join request matching dbOpts
func (dao *DAO) GetJoinRequest(ctx context.Context, dbOpts *Options) (*models.GroupJoinRequest, error) {
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.GroupJoinRequestColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.GroupJoinRequest](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListJoinRequests returns join requests matching dbOpts
func (dao *DAO) ListJoinRequests(ctx context.Context, dbOpts *Options) ([]*models.GroupJoinRequest, error) {
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.GroupJoinRequestColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.GroupJoinRequest](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateJoinRequest updates the status of a join request
func (dao *DAO) UpdateJoinRequest(ctx context.Context, id string, status types.JoinRequestStatus) error {
	now := types.NowDateTime().String()
	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: id})

	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithData(map[string]interface{}{
			models.JOIN_REQUEST_STATUS: status,
			models.BASE_UPDATED_AT:     now,
		}).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteJoinRequests deletes records from the group_join_requests table
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

	type groupIDRow struct {
		GroupID string `db:"group_id"`
	}

	dbOpts := NewOptions().WithWhere(squirrel.Eq{
		models.JOIN_REQUEST_USER_ID:  userID,
		models.JOIN_REQUEST_GROUP_ID: groupIDs,
		models.JOIN_REQUEST_STATUS:   types.JoinPending,
	})
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns(models.JOIN_REQUEST_GROUP_ID).
		SetDbOpts(dbOpts)

	rows, err := listGeneric[groupIDRow](ctx, dao, *builderOpts)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.GroupID)
	}

	return ids, nil
}
