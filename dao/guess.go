package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGuess inserts a single guess attempt row
func (dao *DAO) CreateGuess(ctx context.Context, g *models.Guess) error {
	if g.ID == "" {
		g.RefreshId()
	}
	g.RefreshCreatedAt()

	return createGeneric(ctx, dao, *newBuilderOptions(models.GUESS_TABLE).WithData(map[string]interface{}{
		models.BASE_ID: g.ID, "round_id": g.RoundID, "user_id": g.UserID,
		"attempt": g.Attempt, "word": g.Word, "result": g.Result, "outcome": g.Outcome,
		models.BASE_CREATED_AT: g.CreatedAt,
	}))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGuesses returns guess rows matching the given options
func (dao *DAO) ListGuesses(ctx context.Context, dbOpts *Options) ([]*models.Guess, error) {
	if dbOpts == nil {
		dbOpts = NewOptions()
	}
	dbOpts = dbOpts.WithOrderBy("attempt ASC")

	return listGeneric[models.Guess](ctx, dao, *newBuilderOptions(models.GUESS_TABLE).
		WithColumns(models.GuessColumns()...).
		SetDbOpts(dbOpts))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGuessesForRoundUser returns all guesses for a user in a round
func (dao *DAO) ListGuessesForRoundUser(ctx context.Context, roundID, userID string) ([]*models.Guess, error) {
	return dao.ListGuesses(ctx, NewOptions().WithWhere(squirrel.Eq{
		"round_id": roundID,
		"user_id":  userID,
	}))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CountGuessesForRoundUser returns how many attempts a user has made in a round
func (dao *DAO) CountGuessesForRoundUser(ctx context.Context, roundID, userID string) (int, error) {
	return countGeneric(ctx, dao, *newBuilderOptions(models.GUESS_TABLE).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"round_id": roundID, "user_id": userID})))
}
