package dao

import (
	"context"

	"github.com/geerew/friendle/models"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGuess inserts a single guess attempt row
func (dao *DAO) CreateGuess(ctx context.Context, g *models.Guess) error {
	if g.ID == "" {
		g.RefreshId()
	}
	g.RefreshCreatedAt()

	builderOpts := newBuilderOptions(models.GUESS_TABLE).
		WithData(map[string]interface{}{
			models.BASE_ID:         g.ID,
			models.GUESS_ROUND_ID:  g.RoundID,
			models.GUESS_USER_ID:   g.UserID,
			models.GUESS_ATTEMPT:   g.Attempt,
			models.GUESS_WORD:      g.Word,
			models.GUESS_RESULT:    g.Result,
			models.GUESS_OUTCOME:   g.Outcome,
			models.BASE_CREATED_AT: g.CreatedAt,
		})

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGuesses returns guess rows matching dbOpts
func (dao *DAO) ListGuesses(ctx context.Context, dbOpts *Options) ([]*models.Guess, error) {
	if dbOpts == nil {
		dbOpts = NewOptions()
	}

	dbOpts = dbOpts.WithOrderBy(models.GUESS_ATTEMPT + " ASC")

	builderOpts := newBuilderOptions(models.GUESS_TABLE).
		WithColumns(models.GuessColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.Guess](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CountGuesses returns the number of guesses matching dbOpts
func (dao *DAO) CountGuesses(ctx context.Context, dbOpts *Options) (int, error) {
	builderOpts := newBuilderOptions(models.GUESS_TABLE).SetDbOpts(dbOpts)

	return countGeneric(ctx, dao, *builderOpts)
}
