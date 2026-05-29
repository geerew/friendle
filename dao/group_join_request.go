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
			models.BASE_ID:         r.ID,
			"group_id":             r.GroupID,
			"user_id":              r.UserID,
			"status":               r.Status,
			models.BASE_CREATED_AT: r.CreatedAt,
			models.BASE_UPDATED_AT: r.UpdatedAt,
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

// UpdateJoinRequestStatus updates the status of a join request
func (dao *DAO) UpdateJoinRequestStatus(ctx context.Context, id string, status types.JoinRequestStatus) error {
	now := types.NowDateTime().String()
	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: id})

	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithData(map[string]interface{}{
			"status":               status,
			models.BASE_UPDATED_AT: now,
		}).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeletePendingJoinRequest removes a pending join request for the given group and user
func (dao *DAO) DeletePendingJoinRequest(ctx context.Context, groupID, userID string) error {
	if groupID == "" || userID == "" {
		return utils.ErrWhere
	}

	dbOpts := NewOptions().WithWhere(squirrel.Eq{
		"group_id": groupID,
		"user_id":  userID,
		"status":   types.JoinPending,
	})
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
		"user_id":  userID,
		"group_id": groupIDs,
		"status":   types.JoinPending,
	})
	builderOpts := newBuilderOptions(models.JOIN_REQUEST_TABLE).
		WithColumns("group_id").
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
