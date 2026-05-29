package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
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
		models.BASE_ID:         r.ID,
		"group_id":             r.GroupID,
		"round_date":           r.RoundDate,
		"picker_user_id":       r.PickerUserID,
		"status":               r.Status,
		models.BASE_CREATED_AT: r.CreatedAt,
		models.BASE_UPDATED_AT: r.UpdatedAt,
	}

	if r.WordPlain != nil {
		data["word_plain"] = *r.WordPlain
	}

	builderOpts := newBuilderOptions(models.ROUND_TABLE).WithData(data)

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetRound returns a round matching dbOpts
func (dao *DAO) GetRound(ctx context.Context, dbOpts *Options) (*models.Round, error) {
	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.Round](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRounds returns rounds matching dbOpts
func (dao *DAO) ListRounds(ctx context.Context, dbOpts *Options) ([]*models.Round, error) {
	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.Round](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateRound updates mutable round fields
func (dao *DAO) UpdateRound(ctx context.Context, r *models.Round) error {
	r.RefreshUpdatedAt()

	data := map[string]interface{}{
		"status":               r.Status,
		models.BASE_UPDATED_AT: r.UpdatedAt,
	}
	if r.WordPlain != nil {
		data["word_plain"] = *r.WordPlain
	}

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: r.ID})
	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithData(data).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}
