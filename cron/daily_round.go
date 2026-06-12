package cron

import (
	"context"
	"time"

	"github.com/geerew/friendle/service"
	"github.com/geerew/friendle/utils/logger"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// dailyRound creates today's rounds and closes yesterday's at each local midnight
type dailyRound struct {
	rounds *service.Rounds
	logger *logger.Logger
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// start runs catch-up immediately and waits for each local midnight
func (d *dailyRound) start(ctx context.Context) {
	d.run(ctx)

	for {
		now := time.Now()
		next := nextLocalMidnight(now)
		wait := time.Until(next)
		if wait <= 0 {
			wait = time.Second
		}

		timer := time.NewTimer(wait)

		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			d.run(ctx)
		}
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// run creates today's rounds for all playable groups and closes yesterday's rounds
func (d *dailyRound) run(ctx context.Context) {
	if err := d.rounds.EnsureDailyRounds(ctx); err != nil {
		d.logger.Error().Err(err).Msg("Failed to ensure daily rounds")
		return
	}

	d.logger.Info().Msg("Daily rounds ensured")
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// nextLocalMidnight returns the next local midnight strictly after t
func nextLocalMidnight(t time.Time) time.Time {
	local := t.In(time.Local)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
}
