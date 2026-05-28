package cron

import (
	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/utils/filesystem"
	"github.com/geerew/friendle/utils/logger"
	"github.com/robfig/cron/v3"
)

type CronConfig struct {
	DataDb  database.Database
	FS      *filesystem.FS
	DataDir string
	Logger  *logger.Logger
}

type Cron struct {
	Rounds *RoundScheduler
	cron   *cron.Cron
}

func NewCronScheduler(config *CronConfig) *Cron {
	c := &Cron{cron: cron.New()}
	c.Rounds = newRoundScheduler(config.DataDb, config.FS, config.DataDir, config.Logger)
	_, _ = c.cron.AddFunc("0 0 * * *", func() { c.Rounds.AdvanceAll() })
	_, _ = c.cron.AddFunc("@every 1m", func() { c.Rounds.CheckCompletionAll() })
	return c
}

func (c *Cron) Start() {
	go c.Rounds.AdvanceAll()
	c.cron.Start()
}
