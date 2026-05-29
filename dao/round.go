package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
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

// GetRound returns a round matching dbOpts
//
// Participations, guesses, and users are not included by default. Enable them with
// WithParticipations(), WithGuesses(), and WithUsers() on the options
func (dao *DAO) GetRound(ctx context.Context, dbOpts *Options) (*models.Round, error) {
	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	includeParticipations := dbOpts != nil && dbOpts.IncludeParticipations
	includeGuesses := dbOpts != nil && dbOpts.IncludeGuesses
	includeUsers := dbOpts != nil && dbOpts.IncludeUsers

	if !includeParticipations && !includeGuesses && !includeUsers {
		return getGeneric[models.Round](ctx, dao, *builderOpts)
	}

	round, err := getGeneric[models.Round](ctx, dao, *builderOpts)
	if err != nil {
		return nil, err
	}

	if round == nil {
		return nil, nil
	}

	if err := attachRoundRelations(ctx, dao, []*models.Round{round}, includeParticipations, includeGuesses, includeUsers); err != nil {
		return nil, err
	}

	return round, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRounds returns rounds matching dbOpts
//
// Participations, guesses, and users are not included by default. Enable them with
// WithParticipations(), WithGuesses(), and WithUsers() on the options
func (dao *DAO) ListRounds(ctx context.Context, dbOpts *Options) ([]*models.Round, error) {
	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithColumns(models.RoundColumns()...).
		SetDbOpts(dbOpts)

	includeParticipations := dbOpts != nil && dbOpts.IncludeParticipations
	includeGuesses := dbOpts != nil && dbOpts.IncludeGuesses
	includeUsers := dbOpts != nil && dbOpts.IncludeUsers

	if !includeParticipations && !includeGuesses && !includeUsers {
		return listGeneric[models.Round](ctx, dao, *builderOpts)
	}

	rounds, err := listGeneric[models.Round](ctx, dao, *builderOpts)
	if err != nil {
		return nil, err
	}

	if len(rounds) == 0 {
		return rounds, nil
	}

	if err := attachRoundRelations(ctx, dao, rounds, includeParticipations, includeGuesses, includeUsers); err != nil {
		return nil, err
	}

	return rounds, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateRound updates mutable round fields
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

// DeleteRounds deletes records from the rounds table
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// attachRoundRelations attaches participations, guesses, and users to the given rounds
func attachRoundRelations(ctx context.Context, dao *DAO, rounds []*models.Round, includeParticipations, includeGuesses, includeUsers bool) error {
	if len(rounds) == 0 {
		return nil
	}

	activeRoundIDs := make([]string, 0, len(rounds))
	for _, round := range rounds {
		if round.Status == types.RoundAwaitingWord || round.Status == types.RoundSkipped {
			if includeParticipations {
				round.Participations = []*models.RoundParticipation{}
			}
			continue
		}

		activeRoundIDs = append(activeRoundIDs, round.ID)
	}

	participationsByRound := make(map[string][]*models.RoundParticipation)
	if includeParticipations && len(activeRoundIDs) > 0 {
		participations, err := dao.ListRoundParticipations(ctx, NewOptions().WithWhere(squirrel.Eq{
			models.ROUND_PARTICIPATION_ROUND_ID: activeRoundIDs,
		}))
		if err != nil {
			return err
		}

		for _, p := range participations {
			participationsByRound[p.RoundID] = append(participationsByRound[p.RoundID], p)
		}
	}

	guessesByRoundUser := make(map[string]map[string][]*models.Guess)
	if includeGuesses && len(activeRoundIDs) > 0 {
		guesses, err := dao.ListGuesses(ctx, NewOptions().WithWhere(squirrel.Eq{
			models.GUESS_ROUND_ID: activeRoundIDs,
		}))
		if err != nil {
			return err
		}

		for _, g := range guesses {
			if guessesByRoundUser[g.RoundID] == nil {
				guessesByRoundUser[g.RoundID] = make(map[string][]*models.Guess)
			}

			guessesByRoundUser[g.RoundID][g.UserID] = append(guessesByRoundUser[g.RoundID][g.UserID], g)
		}
	}

	var userMap map[string]*models.User
	if includeUsers {
		userIDs := make([]string, 0)
		for _, round := range rounds {
			if round.PickerUserID != "" {
				userIDs = append(userIDs, round.PickerUserID)
			}
		}

		if includeParticipations {
			for _, participations := range participationsByRound {
				for _, p := range participations {
					userIDs = append(userIDs, p.UserID)
				}
			}
		}

		var err error
		userMap, err = usersByIDs(ctx, dao, userIDs)
		if err != nil {
			return err
		}
	}

	for _, round := range rounds {
		if round.Status == types.RoundAwaitingWord || round.Status == types.RoundSkipped {
			continue
		}

		if includeParticipations {
			round.Participations = participationsByRound[round.ID]
			if round.Participations == nil {
				round.Participations = []*models.RoundParticipation{}
			}
		}

		if includeGuesses {
			for _, p := range round.Participations {
				if byUser := guessesByRoundUser[round.ID]; byUser != nil {
					p.Guesses = byUser[p.UserID]
				}
			}
		}

		if includeUsers {
			round.Picker = userMap[round.PickerUserID]
			for _, p := range round.Participations {
				p.User = userMap[p.UserID]
			}
		}
	}

	return nil
}
