package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
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
			models.BASE_ID:         p.ID,
			"round_id":             p.RoundID,
			"user_id":              p.UserID,
			"solved":               p.Solved,
			"finished":             p.Finished,
			"score":                p.Score,
			"first_guess_at":       p.FirstGuessAt,
			"completed_at":         p.CompletedAt,
			models.BASE_CREATED_AT: p.CreatedAt,
			models.BASE_UPDATED_AT: p.UpdatedAt,
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
	p.RefreshUpdatedAt()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: p.ID})
	builderOpts := newBuilderOptions(models.ROUND_PARTICIPATION_TABLE).
		WithData(map[string]interface{}{
			"solved":               p.Solved,
			"finished":             p.Finished,
			"score":                p.Score,
			"first_guess_at":       p.FirstGuessAt,
			"completed_at":         p.CompletedAt,
			models.BASE_UPDATED_AT: p.UpdatedAt,
		}).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}
