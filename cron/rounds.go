package cron

import (
	"context"
	"math/rand"
	"time"

	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/logger"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundScheduler advances daily rounds and closes stale rounds at day boundaries
type RoundScheduler struct {
	dao    *dao.DAO
	logger *logger.Logger
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AdvanceAll creates today's round for each group that does not have one yet
func (rs *RoundScheduler) AdvanceAll() {
	ctx := context.Background()
	groups, err := rs.dao.ListAllGroups(ctx)
	if err != nil {
		rs.logger.Error().Err(err).Msg("advance rounds: list groups")
		return
	}

	for _, g := range groups {
		rs.advanceGroup(ctx, g)
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// newRoundScheduler creates a RoundScheduler backed by the data database
func newRoundScheduler(db database.Database, log *logger.Logger) *RoundScheduler {
	return &RoundScheduler{dao: dao.New(db), logger: log}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// advanceGroup opens today's round for a group when eligible
func (rs *RoundScheduler) advanceGroup(ctx context.Context, g *models.Group) {
	roundDate := today()
	if cur, _ := rs.dao.GetCurrentRound(ctx, g.ID, roundDate); cur != nil {
		return
	}

	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	if prev, _ := rs.dao.GetCurrentRound(ctx, g.ID, yesterday); prev != nil {
		rs.closeRound(ctx, prev)
	}

	count, err := rs.dao.CountGroupMembers(ctx, g.ID)
	if err != nil || count < 2 {
		return
	}

	picker, err := rs.pickMember(ctx, g.ID)
	if err != nil || picker == nil {
		return
	}

	_ = rs.dao.IncrementTimesPicked(ctx, picker.ID)
	round := &models.Round{
		GroupID:      g.ID,
		RoundDate:    roundDate,
		PickerUserID: picker.UserID,
		Status:       types.RoundAwaitingWord,
	}
	if err := rs.dao.CreateRound(ctx, round); err != nil {
		rs.logger.Error().Err(err).Str("group", g.ID).Msg("create round")
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// pickMember chooses a random member from those picked least often
func (rs *RoundScheduler) pickMember(ctx context.Context, groupID string) (*models.GroupMember, error) {
	pool, err := rs.dao.MembersWithMinPicks(ctx, groupID)
	if err != nil || len(pool) == 0 {
		return nil, err
	}

	return pool[rand.Intn(len(pool))], nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// closeRound finalizes a stale round from a previous day
func (rs *RoundScheduler) closeRound(ctx context.Context, r *models.Round) {
	switch r.Status {
	case types.RoundAwaitingWord:
		r.Status = types.RoundSkipped
		_ = rs.dao.UpdateRound(ctx, r)
	case types.RoundActive:
		rs.completeRound(ctx, r)
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// completeRound marks unfinished participations done and sets the round to completed
func (rs *RoundScheduler) completeRound(ctx context.Context, r *models.Round) {
	participations, _ := rs.dao.ListRoundParticipationsForRound(ctx, r.ID)
	for _, p := range participations {
		if p.Finished {
			continue
		}

		p.Finished = true
		_ = rs.dao.UpdateRoundParticipation(ctx, p)
	}

	r.Status = types.RoundCompleted
	_ = rs.dao.UpdateRound(ctx, r)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// today returns the current calendar date used as the round key
func today() string {
	return time.Now().Format("2006-01-02")
}
