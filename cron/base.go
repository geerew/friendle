package cron

import (
	"context"

	"github.com/geerew/friendle/service"
	"github.com/geerew/friendle/utils/logger"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Config holds dependencies for the scheduled job runner
type Config struct {
	Rounds           *service.Rounds
	DailyRoundLogger *logger.Logger
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Cron manages background jobs
type Cron struct {
	DailyRound *dailyRound
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewAndStart creates a Cron scheduler and starts registered jobs
func NewAndStart(ctx context.Context, config *Config) *Cron {
	c := &Cron{
		DailyRound: &dailyRound{
			rounds: config.Rounds,
			logger: config.DailyRoundLogger,
		},
	}

	go c.DailyRound.start(ctx)

	return c
}
