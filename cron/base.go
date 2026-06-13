package cron

import (
	"context"

	"github.com/geerew/friendle/service"
	"github.com/geerew/friendle/utils/logger"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Config holds dependencies for the scheduled job runner
type Config struct {
	Rounds *service.Rounds
	Logger *logger.Logger
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Cron manages background jobs
type Cron struct {
	CloseStaleRounds *closeStaleRounds
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewAndStart creates a Cron scheduler and starts registered jobs
func NewAndStart(ctx context.Context, config *Config) *Cron {
	c := &Cron{
		CloseStaleRounds: &closeStaleRounds{
			rounds: config.Rounds,
			logger: config.Logger,
		},
	}

	go c.CloseStaleRounds.start(ctx)

	return c
}
