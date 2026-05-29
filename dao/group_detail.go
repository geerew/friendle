package dao

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupDetailLoad selects which relations to populate on a group detail load
type GroupDetailLoad struct {
	Members        bool
	JoinRequests   bool
	Leaderboard    bool
	CurrentRound   bool
	PreviousRounds bool
	RoundLimit     int
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroupDetail loads a group and optional related data
func (dao *DAO) GetGroupDetail(ctx context.Context, groupID string, load GroupDetailLoad) (*models.GroupDetail, error) {
	group, err := dao.GetGroup(ctx, groupID)
	if err != nil || group == nil {
		return nil, err
	}

	if load.Members {
		members, err := dao.ListGroupMembers(ctx, groupID)
		if err != nil {
			return nil, err
		}
		for _, m := range members {
			u, _ := dao.GetUser(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: m.UserID}))
			m.User = u
		}
		group.Members = members
	}

	if load.JoinRequests {
		group.JoinRequests, err = dao.ListPendingJoinRequests(ctx, groupID)
		if err != nil {
			return nil, err
		}
	}

	if load.Leaderboard {
		group.Leaderboard, err = dao.ListLeaderboardEntries(ctx, groupID)
		if err != nil {
			return nil, err
		}
	}

	roundDate := time.Now().Format("2006-01-02")
	if load.CurrentRound {
		group.CurrentRound, err = dao.GetCurrentRound(ctx, groupID, roundDate)
		if err != nil {
			return nil, err
		}
	}

	if load.PreviousRounds {
		limit := load.RoundLimit
		if limit <= 0 {
			limit = 10
		}
		p := pagination.New(1, limit)
		rounds, err := dao.ListRoundsForGroup(ctx, groupID, NewOptions().WithPagination(p))
		if err != nil {
			return nil, err
		}
		previous := make([]*models.Round, 0, len(rounds))
		for _, rd := range rounds {
			if group.CurrentRound != nil && rd.ID == group.CurrentRound.ID {
				continue
			}
			previous = append(previous, rd)
		}
		group.PreviousRounds = previous
	}

	return &models.GroupDetail{Group: group}, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundDetailLoad selects which relations to populate on a group round load
type RoundDetailLoad struct {
	Participations bool
	Guesses        bool
	Picker         bool
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroupRound loads a round with player participations and guesses
func (dao *DAO) GetGroupRound(ctx context.Context, roundID string, load RoundDetailLoad) (*models.GroupRound, error) {
	round, err := dao.GetRound(ctx, roundID)
	if err != nil || round == nil {
		return nil, err
	}

	if load.Picker {
		round.Picker, _ = dao.GetUser(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: round.PickerUserID}))
	}

	gr := &models.GroupRound{Round: round}
	if !load.Participations {
		return gr, nil
	}

	participations, err := dao.ListRoundParticipationsForRound(ctx, roundID)
	if err != nil {
		return nil, err
	}

	players := make([]*models.UserRound, 0, len(participations))
	for _, p := range participations {
		ur := &models.UserRound{Participation: p}
		if load.Guesses {
			ur.Guesses, _ = dao.ListGuessesForRoundUser(ctx, roundID, p.UserID)
			p.Guesses = ur.Guesses
		}
		u, _ := dao.GetUser(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: p.UserID}))
		p.User = u
		ur.User = u
		players = append(players, ur)
	}

	round.Participations = participations
	gr.Players = players

	return gr, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetUserRound loads one user's participation and guesses for a round
func (dao *DAO) GetUserRound(ctx context.Context, roundID, userID string, loadGuesses bool) (*models.UserRound, error) {
	participation, err := dao.GetRoundParticipation(ctx, roundID, userID)
	if err != nil {
		return nil, err
	}

	ur := &models.UserRound{Participation: participation}
	if participation == nil {
		return ur, nil
	}

	if loadGuesses {
		ur.Guesses, err = dao.ListGuessesForRoundUser(ctx, roundID, userID)
		if err != nil {
			return nil, err
		}
		participation.Guesses = ur.Guesses
	}

	ur.User, _ = dao.GetUser(ctx, NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: userID}))

	return ur, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundIsRevealed reports whether guess details may be shown to all members
func RoundIsRevealed(status types.RoundStatus) bool {
	return status == types.RoundCompleted || status == types.RoundSkipped
}
