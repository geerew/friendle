package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateRound inserts a round row
func (dao *DAO) CreateRound(ctx context.Context, r *models.Round) error {
	if r.ID == "" {
		r.RefreshId()
	}
	r.RefreshCreatedAt()
	r.RefreshUpdatedAt()
	data := map[string]interface{}{
		models.BASE_ID: r.ID, "group_id": r.GroupID, "round_date": r.RoundDate,
		"picker_user_id": r.PickerUserID, "status": r.Status,
		models.BASE_CREATED_AT: r.CreatedAt, models.BASE_UPDATED_AT: r.UpdatedAt,
	}
	if r.WordHash != nil {
		data["word_hash"] = *r.WordHash
	}
	if r.WordPlain != nil {
		data["word_plain"] = *r.WordPlain
	}

	return createGeneric(ctx, dao, *newBuilderOptions(models.ROUND_TABLE).WithData(data))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetRound returns a round by ID
func (dao *DAO) GetRound(ctx context.Context, id string) (*models.Round, error) {
	return getGeneric[models.Round](ctx, dao, *newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: id})).WithLimit(1))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetCurrentRound returns the round for a group on a given date
func (dao *DAO) GetCurrentRound(ctx context.Context, groupID, roundDate string) (*models.Round, error) {
	return getGeneric[models.Round](ctx, dao, *newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"group_id": groupID, "round_date": roundDate})).
		WithLimit(1))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRoundsForGroup returns rounds for a group ordered by date descending
func (dao *DAO) ListRoundsForGroup(ctx context.Context, groupID string, dbOpts *Options) ([]*models.Round, error) {
	opts := NewOptions().WithWhere(squirrel.Eq{"group_id": groupID}).WithOrderBy(models.ROUND_ROUND_DATE + " DESC")
	if dbOpts != nil {
		if dbOpts.Pagination != nil {
			opts = opts.WithPagination(dbOpts.Pagination)
		}
		if dbOpts.Where != nil {
			opts = opts.WithWhere(squirrel.And{opts.Where, dbOpts.Where})
		}
	}

	return listGeneric[models.Round](ctx, dao, *newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(opts))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateRound updates mutable round fields
func (dao *DAO) UpdateRound(ctx context.Context, r *models.Round) error {
	r.RefreshUpdatedAt()
	data := map[string]interface{}{
		"status": r.Status, models.BASE_UPDATED_AT: r.UpdatedAt,
	}
	if r.WordHash != nil {
		data["word_hash"] = *r.WordHash
	}
	if r.WordPlain != nil {
		data["word_plain"] = *r.WordPlain
	}
	_, err := updateGeneric(ctx, dao, *newBuilderOptions(models.ROUND_TABLE).WithData(data).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: r.ID})))

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListActiveRounds returns rounds awaiting a word or in progress
func (dao *DAO) ListActiveRounds(ctx context.Context) ([]*models.Round, error) {
	return listGeneric[models.Round](ctx, dao, *newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Or{
			squirrel.Eq{"status": types.RoundAwaitingWord},
			squirrel.Eq{"status": types.RoundActive},
		})))
}
