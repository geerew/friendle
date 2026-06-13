package app

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Component names a log source for filtering and structured output
type Component string

const (
	ComponentApp  Component = "app"
	ComponentAPI  Component = "api"
	ComponentCron Component = "cron"
)
