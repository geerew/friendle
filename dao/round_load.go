package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundLoad selects which relations to populate on a round
type RoundLoad struct {
	Participations bool
	Guesses        bool
	Users          bool
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetRoundLoaded returns a round with optional relations loaded
func (dao *DAO) GetRoundLoaded(ctx context.Context, roundID string, load RoundLoad) (*models.Round, error) {
	round, err := dao.GetRound(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: roundID}))
	if err != nil || round == nil {
		return round, err
	}

	if err := dao.LoadRound(ctx, round, load); err != nil {
		return nil, err
	}

	return round, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// LoadRound populates relation fields on an existing round row
func (dao *DAO) LoadRound(ctx context.Context, round *models.Round, load RoundLoad) error {
	if round.PickerUserID != "" {
		round.Picker, _ = dao.GetUser(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: round.PickerUserID}))
	}

	if !load.Participations {
		return nil
	}

	if round.Status == types.RoundAwaitingWord || round.Status == types.RoundSkipped {
		round.Participations = []*models.RoundParticipation{}
		return nil
	}

	participations, err := dao.ListRoundParticipations(ctx, NewOptions().WithWhere(squirrel.Eq{
		models.ROUND_PARTICIPATION_ROUND_ID: round.ID,
	}))
	if err != nil {
		return err
	}

	for _, p := range participations {
		if load.Guesses {
			p.Guesses, err = dao.ListGuesses(ctx, NewOptions().WithWhere(squirrel.Eq{
				models.GUESS_ROUND_ID: round.ID,
				models.GUESS_USER_ID:  p.UserID,
			}))
			if err != nil {
				return err
			}
		}
		if load.Users {
			p.User, _ = dao.GetUser(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: p.UserID}))
		}
	}

	round.Participations = participations

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundIsRevealed reports whether the round word and all guesses may be shown
func RoundIsRevealed(status types.RoundStatus) bool {
	return status == types.RoundCompleted || status == types.RoundSkipped
}
