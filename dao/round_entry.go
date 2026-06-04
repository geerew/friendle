package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateRoundEntry inserts a round entry record
func (dao *DAO) CreateRoundEntry(ctx context.Context, e *models.RoundEntry) error {
	if e.ID == "" {
		e.RefreshId()
	}
	e.RefreshCreatedAt()
	e.RefreshUpdatedAt()

	builderOpts := newBuilderOptions(models.ROUND_ENTRY_TABLE).
		WithData(map[string]interface{}{
			models.BASE_ID:                    e.ID,
			models.ROUND_ENTRY_ROUND_ID:       e.RoundID,
			models.ROUND_ENTRY_USER_ID:        e.UserID,
			models.ROUND_ENTRY_SOLVED:         e.Solved,
			models.ROUND_ENTRY_FINISHED:       e.Finished,
			models.ROUND_ENTRY_SCORE:          e.Score,
			models.ROUND_ENTRY_FIRST_GUESS_AT: e.FirstGuessAt,
			models.ROUND_ENTRY_COMPLETED_AT:   e.CompletedAt,
			models.BASE_CREATED_AT:            e.CreatedAt,
			models.BASE_UPDATED_AT:            e.UpdatedAt,
		})

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetRoundEntry returns a round entry record
func (dao *DAO) GetRoundEntry(ctx context.Context, dbOpts *Options) (*models.RoundEntry, error) {
	builderOpts := newBuilderOptions(models.ROUND_ENTRY_TABLE).
		WithColumns(models.RoundEntryColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.RoundEntry](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRoundEntries returns round entry records
func (dao *DAO) ListRoundEntries(ctx context.Context, dbOpts *Options) ([]*models.RoundEntry, error) {
	builderOpts := newBuilderOptions(models.ROUND_ENTRY_TABLE).
		WithColumns(models.RoundEntryColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.RoundEntry](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateRoundEntry updates mutable round entry fields
func (dao *DAO) UpdateRoundEntry(ctx context.Context, e *models.RoundEntry) error {
	if e.ID == "" {
		return utils.ErrId
	}

	e.RefreshUpdatedAt()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: e.ID})
	builderOpts := newBuilderOptions(models.ROUND_ENTRY_TABLE).
		WithData(map[string]interface{}{
			models.ROUND_ENTRY_SOLVED:         e.Solved,
			models.ROUND_ENTRY_FINISHED:       e.Finished,
			models.ROUND_ENTRY_SCORE:          e.Score,
			models.ROUND_ENTRY_FIRST_GUESS_AT: e.FirstGuessAt,
			models.ROUND_ENTRY_COMPLETED_AT:   e.CompletedAt,
			models.BASE_UPDATED_AT:            e.UpdatedAt,
		}).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteRoundEntries deletes round entry records
//
// Errors when a where clause is not provided
func (dao *DAO) DeleteRoundEntries(ctx context.Context, dbOpts *Options) error {
	if dbOpts == nil || dbOpts.Where == nil {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.ROUND_ENTRY_TABLE).SetDbOpts(dbOpts)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)
	return err
}
