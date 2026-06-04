package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateRound inserts a round record
func (dao *DAO) CreateRound(ctx context.Context, r *models.Round) error {
	if r.ID == "" {
		r.RefreshId()
	}
	r.RefreshCreatedAt()
	r.RefreshUpdatedAt()

	data := map[string]interface{}{
		models.BASE_ID:              r.ID,
		models.ROUND_GROUP_ID:       r.GroupID,
		models.ROUND_ROUND_DATE:     r.RoundDate,
		models.ROUND_PICKER_USER_ID: r.PickerUserID,
		models.ROUND_STATUS:         r.Status,
		models.BASE_CREATED_AT:      r.CreatedAt,
		models.BASE_UPDATED_AT:      r.UpdatedAt,
	}

	if r.WordPlain != nil {
		data[models.ROUND_WORD_PLAIN] = *r.WordPlain
	}

	builderOpts := newBuilderOptions(models.ROUND_TABLE).WithData(data)

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetRound returns a round record
func (dao *DAO) GetRound(ctx context.Context, dbOpts *Options) (*models.Round, error) {
	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.Round](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRounds returns round records
func (dao *DAO) ListRounds(ctx context.Context, dbOpts *Options) ([]*models.Round, error) {
	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.Round](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateRound updates a round record
func (dao *DAO) UpdateRound(ctx context.Context, r *models.Round) error {
	if r.ID == "" {
		return utils.ErrId
	}

	r.RefreshUpdatedAt()

	var wordPlain interface{}
	if r.WordPlain != nil {
		wordPlain = *r.WordPlain
	}

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: r.ID})
	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithData(map[string]interface{}{
			models.ROUND_STATUS:     r.Status,
			models.ROUND_WORD_PLAIN: wordPlain,
			models.BASE_UPDATED_AT:  r.UpdatedAt,
		}).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteRounds deletes round records
//
// Errors when a where clause is not provided
func (dao *DAO) DeleteRounds(ctx context.Context, dbOpts *Options) error {
	if dbOpts == nil || dbOpts.Where == nil {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.ROUND_TABLE).SetDbOpts(dbOpts)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)
	return err
}
