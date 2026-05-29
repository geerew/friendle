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

	return createGeneric(ctx, dao, *newBuilderOptions(models.ROUND_PARTICIPATION_TABLE).WithData(map[string]interface{}{
		models.BASE_ID: p.ID, "round_id": p.RoundID, "user_id": p.UserID,
		"solved": p.Solved, "finished": p.Finished, "score": p.Score,
		"first_guess_at": p.FirstGuessAt, "completed_at": p.CompletedAt,
		models.BASE_CREATED_AT: p.CreatedAt, models.BASE_UPDATED_AT: p.UpdatedAt,
	}))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetRoundParticipation returns a user's participation in a round
func (dao *DAO) GetRoundParticipation(ctx context.Context, roundID, userID string) (*models.RoundParticipation, error) {
	return getGeneric[models.RoundParticipation](ctx, dao, *newBuilderOptions(models.ROUND_PARTICIPATION_TABLE).
		WithColumns(models.RoundParticipationColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"round_id": roundID, "user_id": userID})).
		WithLimit(1))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRoundParticipations returns participations matching the given options
func (dao *DAO) ListRoundParticipations(ctx context.Context, dbOpts *Options) ([]*models.RoundParticipation, error) {
	return listGeneric[models.RoundParticipation](ctx, dao, *newBuilderOptions(models.ROUND_PARTICIPATION_TABLE).
		WithColumns(models.RoundParticipationColumns()...).
		SetDbOpts(dbOpts))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRoundParticipationsForRound returns all participations for a round
func (dao *DAO) ListRoundParticipationsForRound(ctx context.Context, roundID string) ([]*models.RoundParticipation, error) {
	return dao.ListRoundParticipations(ctx, NewOptions().WithWhere(squirrel.Eq{"round_id": roundID}))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateRoundParticipation updates mutable participation fields
func (dao *DAO) UpdateRoundParticipation(ctx context.Context, p *models.RoundParticipation) error {
	p.RefreshUpdatedAt()
	_, err := updateGeneric(ctx, dao, *newBuilderOptions(models.ROUND_PARTICIPATION_TABLE).WithData(map[string]interface{}{
		"solved": p.Solved, "finished": p.Finished, "score": p.Score,
		"first_guess_at": p.FirstGuessAt, "completed_at": p.CompletedAt,
		models.BASE_UPDATED_AT: p.UpdatedAt,
	}).SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: p.ID})))

	return err
}
