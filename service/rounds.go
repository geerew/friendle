package service

import (
	"context"
	"sort"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundTodayResponse represents today's round state for a user in a group
type RoundTodayResponse struct {
	RoundDate   string            `json:"roundDate"`
	Status      types.RoundStatus `json:"status"`
	IsPicker    bool              `json:"isPicker"`
	RoundEndsAt types.DateTime    `json:"roundEndsAt"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Rounds orchestrates daily group round operations
type Rounds struct {
	deps
	now func() time.Time
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// newRounds creates a Rounds service
func newRounds(d deps) *Rounds {
	return &Rounds{
		deps: d,
		now:  time.Now,
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Create creates today's round for a group
func (r *Rounds) Create(ctx context.Context, groupID string) (*models.Round, error) {
	group, err := r.dao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, ErrGroupNotFound
	}

	if group.MemberCount < minPlayableMembers {
		return nil, ErrGroupTooFewMembers
	}

	today := utils.DateString(r.now())

	existing, err := r.getRoundForDate(ctx, groupID, today)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return existing, nil
	}

	members, err := r.dao.ListGroupMembers(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: groupID}))
	if err != nil {
		return nil, err
	}

	picker, err := r.selectPicker(ctx, groupID, members)
	if err != nil {
		return nil, err
	}

	if picker == nil {
		return nil, ErrGroupTooFewMembers
	}

	round := &models.Round{
		GroupID:      groupID,
		RoundDate:    today,
		PickerUserID: picker.UserID,
		Status:       types.RoundAwaitingWord,
	}

	err = r.dao.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := r.dao.CreateRound(txCtx, round); err != nil {
			return err
		}

		picker.TimesPicked++
		return r.dao.UpdateGroupMember(txCtx, picker)
	})
	if err != nil {
		return nil, err
	}

	return round, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Today returns today's round state for userID in a group
func (r *Rounds) Today(ctx context.Context, groupID, userID string) (*RoundTodayResponse, error) {
	group, err := r.dao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, ErrGroupNotFound
	}

	if group.MemberCount < minPlayableMembers {
		return nil, ErrGroupTooFewMembers
	}

	now := r.now()
	today := utils.DateString(now)
	round, err := r.getRoundForDate(ctx, groupID, today)
	if err != nil {
		return nil, err
	}

	if round == nil {
		return nil, ErrRoundNotFound
	}

	return &RoundTodayResponse{
		RoundDate:   round.RoundDate,
		Status:      round.Status,
		IsPicker:    round.PickerUserID == userID,
		RoundEndsAt: types.DateTime(utils.NextMidnight(now)),
	}, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MarkCompleted sets a round to completed, incrementing picker skips when still awaiting a word
func (r *Rounds) MarkCompleted(ctx context.Context, round *models.Round) error {
	if round == nil || round.Status == types.RoundCompleted {
		return nil
	}

	switch round.Status {
	case types.RoundAwaitingWord:
		return r.dao.RunInTransaction(ctx, func(txCtx context.Context) error {
			round.Status = types.RoundCompleted
			if err := r.dao.UpdateRound(txCtx, round); err != nil {
				return err
			}

			member, err := r.dao.GetGroupMember(txCtx, dao.NewOptions().WithWhere(squirrel.And{
				squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: round.GroupID},
				squirrel.Eq{models.GROUP_MEMBER_USER_ID: round.PickerUserID},
			}))
			if err != nil {
				return err
			}

			if member == nil {
				return nil
			}

			member.PickerSkips++
			return r.dao.UpdateGroupMember(txCtx, member)
		})
	case types.RoundActive:
		round.Status = types.RoundCompleted
		return r.dao.UpdateRound(ctx, round)
	default:
		return nil
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Close finds open rounds from before today and marks each completed
func (r *Rounds) Close(ctx context.Context) error {
	today := utils.DateString(r.now())

	dbOpts := dao.NewOptions().WithWhere(squirrel.And{
		squirrel.Lt{models.ROUND_ROUND_DATE: today},
		squirrel.Or{
			squirrel.Eq{models.ROUND_STATUS: types.RoundActive},
			squirrel.Eq{models.ROUND_STATUS: types.RoundAwaitingWord},
		},
	})

	rounds, err := r.dao.ListRounds(ctx, dbOpts)
	if err != nil {
		return err
	}

	for _, round := range rounds {
		if err := r.MarkCompleted(ctx, round); err != nil {
			return err
		}
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Picker
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// selectPicker chooses the next picker, excluding yesterday's picker when possible
func (r *Rounds) selectPicker(ctx context.Context, groupID string, members []*models.GroupMember) (*models.GroupMember, error) {
	if len(members) == 0 {
		return nil, nil
	}

	yesterday := utils.PreviousDateString(r.now())
	round, err := r.getRoundForDate(ctx, groupID, yesterday)
	if err != nil {
		return nil, err
	}

	yesterdayPickerID := ""
	if round != nil {
		yesterdayPickerID = round.PickerUserID
	}

	candidates := utils.FilterMap(members, func(member *models.GroupMember) (*models.GroupMember, bool) {
		if yesterdayPickerID != "" && member.UserID == yesterdayPickerID {
			return nil, false
		}

		return member, true
	})

	if len(candidates) == 0 {
		candidates = members
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		left := candidates[i]
		right := candidates[j]

		if left.TimesPicked != right.TimesPicked {
			return left.TimesPicked < right.TimesPicked
		}

		if left.PickerSkips != right.PickerSkips {
			return left.PickerSkips < right.PickerSkips
		}

		return left.UserID < right.UserID
	})

	return candidates[0], nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getRoundForDate returns a round for a group for a given date
func (r *Rounds) getRoundForDate(ctx context.Context, groupID, date string) (*models.Round, error) {
	dbOpts := dao.NewOptions().WithWhere(squirrel.And{
		squirrel.Eq{models.ROUND_GROUP_ID: groupID},
		squirrel.Eq{models.ROUND_ROUND_DATE: date},
	})

	return r.dao.GetRound(ctx, dbOpts)
}
