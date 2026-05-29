package models

import "fmt"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	GROUP_TABLE = "groups"

	GROUP_NAME           = "name"
	GROUP_CREATED_BY     = "created_by"
	GROUP_INTERVAL_HOURS = "interval_hours"
	GROUP_TIMEZONE       = "timezone"

	GROUP_TABLE_ID              = GROUP_TABLE + "." + BASE_ID
	GROUP_TABLE_CREATED_AT      = GROUP_TABLE + "." + BASE_CREATED_AT
	GROUP_TABLE_UPDATED_AT      = GROUP_TABLE + "." + BASE_UPDATED_AT
	GROUP_TABLE_NAME            = GROUP_TABLE + "." + GROUP_NAME
	GROUP_TABLE_CREATED_BY      = GROUP_TABLE + "." + GROUP_CREATED_BY
	GROUP_TABLE_INTERVAL_HOURS  = GROUP_TABLE + "." + GROUP_INTERVAL_HOURS
	GROUP_TABLE_TIMEZONE        = GROUP_TABLE + "." + GROUP_TIMEZONE
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Group defines the model for a friend group
type Group struct {
	Base
	Name          string `db:"name"`           // Mutable
	CreatedBy     string `db:"created_by"`     // Immutable
	IntervalHours int    `db:"interval_hours"` // Mutable
	Timezone      string `db:"timezone"`       // Mutable

	Members        []*GroupMember       `db:"-"`
	JoinRequests   []*GroupJoinRequest  `db:"-"`
	Leaderboard    []*LeaderboardEntry  `db:"-"`
	CurrentRound   *Round               `db:"-"`
	PreviousRounds []*Round             `db:"-"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupColumns returns the columns for use in a SELECT query
func GroupColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", GROUP_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", GROUP_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", GROUP_TABLE_UPDATED_AT, BASE_UPDATED_AT),
		fmt.Sprintf("%s AS %s", GROUP_TABLE_NAME, GROUP_NAME),
		fmt.Sprintf("%s AS %s", GROUP_TABLE_CREATED_BY, GROUP_CREATED_BY),
		fmt.Sprintf("%s AS %s", GROUP_TABLE_INTERVAL_HOURS, GROUP_INTERVAL_HOURS),
		fmt.Sprintf("%s AS %s", GROUP_TABLE_TIMEZONE, GROUP_TIMEZONE),
	}
}
