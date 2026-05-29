package cron

import (
	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/utils/logger"
	"github.com/robfig/cron/v3"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Config holds dependencies for the scheduled job runner
type Config struct {
	DataDb database.Database
	Logger *logger.Logger
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Cron runs background jobs for round lifecycle management
type Cron struct {
	Rounds *RoundScheduler
	cron   *cron.Cron
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// New creates a cron scheduler and registers round jobs
func New(config *Config) *Cron {
	c := &Cron{cron: cron.New()}
	c.Rounds = newRoundScheduler(config.DataDb, config.Logger)

	_, _ = c.cron.AddFunc("0 0 * * *", func() { c.Rounds.AdvanceAll() })

	return c
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Start runs an immediate round advance, then starts the scheduled jobs
func (c *Cron) Start() {
	go c.Rounds.AdvanceAll()
	c.cron.Start()
}
