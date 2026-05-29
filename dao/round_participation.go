package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateRoundParticipation inserts a round participation row
func (dao *DAO) CreateRoundParticipation(ctx context.Context, p *models.RoundParticipation) error {
	if p.ID == "" {
		p.RefreshId()
	}
	p.RefreshCreatedAt()
	p.RefreshUpdatedAt()

	builderOpts := newBuilderOptions(models.ROUND_PARTICIPATION_TABLE).
		WithData(map[string]interface{}{
			models.BASE_ID:                            p.ID,
			models.ROUND_PARTICIPATION_ROUND_ID:       p.RoundID,
			models.ROUND_PARTICIPATION_USER_ID:        p.UserID,
			models.ROUND_PARTICIPATION_SOLVED:         p.Solved,
			models.ROUND_PARTICIPATION_FINISHED:       p.Finished,
			models.ROUND_PARTICIPATION_SCORE:          p.Score,
			models.ROUND_PARTICIPATION_FIRST_GUESS_AT: p.FirstGuessAt,
			models.ROUND_PARTICIPATION_COMPLETED_AT:   p.CompletedAt,
			models.BASE_CREATED_AT:                    p.CreatedAt,
			models.BASE_UPDATED_AT:                    p.UpdatedAt,
		})

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetRoundParticipation returns a round participation matching dbOpts
func (dao *DAO) GetRoundParticipation(ctx context.Context, dbOpts *Options) (*models.RoundParticipation, error) {
	builderOpts := newBuilderOptions(models.ROUND_PARTICIPATION_TABLE).
		WithColumns(models.RoundParticipationColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.RoundParticipation](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRoundParticipations returns participations matching dbOpts
func (dao *DAO) ListRoundParticipations(ctx context.Context, dbOpts *Options) ([]*models.RoundParticipation, error) {
	builderOpts := newBuilderOptions(models.ROUND_PARTICIPATION_TABLE).
		WithColumns(models.RoundParticipationColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.RoundParticipation](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateRoundParticipation updates mutable participation fields
func (dao *DAO) UpdateRoundParticipation(ctx context.Context, p *models.RoundParticipation) error {
	if p.ID == "" {
		return utils.ErrId
	}

	p.RefreshUpdatedAt()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: p.ID})
	builderOpts := newBuilderOptions(models.ROUND_PARTICIPATION_TABLE).
		WithData(map[string]interface{}{
			models.ROUND_PARTICIPATION_SOLVED:         p.Solved,
			models.ROUND_PARTICIPATION_FINISHED:       p.Finished,
			models.ROUND_PARTICIPATION_SCORE:          p.Score,
			models.ROUND_PARTICIPATION_FIRST_GUESS_AT: p.FirstGuessAt,
			models.ROUND_PARTICIPATION_COMPLETED_AT:   p.CompletedAt,
			models.BASE_UPDATED_AT:                    p.UpdatedAt,
		}).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteRoundParticipations deletes records from the round_participations table
//
// Errors when a where clause is not provided
func (dao *DAO) DeleteRoundParticipations(ctx context.Context, dbOpts *Options) error {
	if dbOpts == nil || dbOpts.Where == nil {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.ROUND_PARTICIPATION_TABLE).SetDbOpts(dbOpts)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)
	return err
}
