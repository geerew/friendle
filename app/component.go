package app

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Component names a log source for filtering and structured output. Add a constant here when
// wiring a new subsystem logger in New
type Component string

const (
	ComponentApp  Component = "app"
	ComponentAPI  Component = "api"
	ComponentCron Component = "cron"
)
