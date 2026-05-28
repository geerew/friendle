package cron

import (
	"context"
	"math/rand"
	"time"

	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/filesystem"
	"github.com/geerew/friendle/utils/logger"
)

type RoundScheduler struct {
	d       *dao.DAO
	fs      *filesystem.FS
	dataDir string
	logger  *logger.Logger
}

func newRoundScheduler(db database.Database, fs *filesystem.FS, dataDir string, log *logger.Logger) *RoundScheduler {
	return &RoundScheduler{dao.New(db), fs, dataDir, log}
}

func today() string {
	return time.Now().Format("2006-01-02")
}

func (s *RoundScheduler) AdvanceAll() {
	ctx := context.Background()
	groups, err := s.d.ListAllGroups(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("advance rounds: list groups")
		return
	}
	for _, g := range groups {
		s.advanceGroup(ctx, g)
	}
}

func (s *RoundScheduler) advanceGroup(ctx context.Context, g *models.Group) {
	roundDate := today()
	if cur, _ := s.d.GetCurrentRound(ctx, g.ID, roundDate); cur != nil {
		return
	}
	// close yesterday's round if still open
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	if prev, _ := s.d.GetCurrentRound(ctx, g.ID, yesterday); prev != nil {
		s.closeRound(ctx, prev)
	}
	count, err := s.d.CountGroupMembers(ctx, g.ID)
	if err != nil || count < 2 {
		return
	}
	picker, err := s.pickMember(ctx, g.ID)
	if err != nil || picker == nil {
		return
	}
	_ = s.d.IncrementTimesPicked(ctx, picker.ID)
	round := &models.Round{
		GroupID: g.ID, RoundDate: roundDate, PickerUserID: picker.UserID,
		Status: models.RoundAwaitingWord,
	}
	if err := s.d.CreateRound(ctx, round); err != nil {
		s.logger.Error().Err(err).Str("group", g.ID).Msg("create round")
	}
}

func (s *RoundScheduler) pickMember(ctx context.Context, groupID string) (*models.GroupMember, error) {
	pool, err := s.d.MembersWithMinPicks(ctx, groupID)
	if err != nil || len(pool) == 0 {
		return nil, err
	}
	return pool[rand.Intn(len(pool))], nil
}

func (s *RoundScheduler) closeRound(ctx context.Context, r *models.Round) {
	switch r.Status {
	case models.RoundAwaitingWord:
		r.Status = models.RoundSkipped
		_ = s.d.UpdateRound(ctx, r)
	case models.RoundActive:
		s.completeRound(ctx, r)
	}
}

func (s *RoundScheduler) CheckCompletionAll() {
	ctx := context.Background()
	rounds, err := s.d.ListActiveRounds(ctx)
	if err != nil {
		return
	}
	for _, r := range rounds {
		if r.Status != models.RoundActive {
			continue
		}
		if s.allGuessersDone(ctx, r) {
			s.completeRound(ctx, r)
		}
	}
}

func (s *RoundScheduler) allGuessersDone(ctx context.Context, r *models.Round) bool {
	members, _ := s.d.ListGroupMembers(ctx, r.GroupID)
	guesses, _ := s.d.ListGuessesForRound(ctx, r.ID)
	done := 0
	needed := 0
	for _, m := range members {
		if m.UserID == r.PickerUserID {
			continue
		}
		needed++
		for _, g := range guesses {
			if g.UserID == m.UserID && g.Finished {
				done++
				break
			}
		}
	}
	return needed > 0 && done >= needed
}

func (s *RoundScheduler) completeRound(ctx context.Context, r *models.Round) {
	guesses, _ := s.d.ListGuessesForRound(ctx, r.ID)
	for _, g := range guesses {
		if g.Finished {
			continue
		}
		g.Finished = true
		_ = s.d.UpdateGuess(ctx, g)
	}
	r.Status = models.RoundCompleted
	_ = s.d.UpdateRound(ctx, r)
}

