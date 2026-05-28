package models

import "github.com/geerew/friendle/utils/types"

const GROUP_TABLE = "groups"

type Group struct {
	Base
	Name          string `db:"name"`
	CreatedBy     string `db:"created_by"`
	IntervalHours int    `db:"interval_hours"`
	Timezone      string `db:"timezone"`
}

const GROUP_MEMBER_TABLE = "group_members"

type GroupMember struct {
	Base
	GroupID     string           `db:"group_id"`
	UserID      string           `db:"user_id"`
	GroupRole   types.GroupRole  `db:"group_role"`
	TimesPicked int              `db:"times_picked"`
	PickerSkips int              `db:"picker_skips"`
}

const JOIN_REQUEST_TABLE = "group_join_requests"

type JoinRequestStatus string

const (
	JoinPending  JoinRequestStatus = "pending"
	JoinApproved JoinRequestStatus = "approved"
	JoinRejected JoinRequestStatus = "rejected"
)

type GroupJoinRequest struct {
	Base
	GroupID string            `db:"group_id"`
	UserID  string            `db:"user_id"`
	Status  JoinRequestStatus `db:"status"`
}

// AdminGroupListRow is a groups row enriched for the site admin list
type AdminGroupListRow struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	MemberCount int    `db:"member_count"`
}

// UserGroupSummaryRow is a group summary for a user's memberships
type UserGroupSummaryRow struct {
	UserID      string          `db:"user_id"`
	ID          string          `db:"id"`
	Name        string          `db:"name"`
	MemberCount int             `db:"member_count"`
	GroupRole   types.GroupRole `db:"group_role"`
}

const ROUND_TABLE = "rounds"

type RoundStatus string

const (
	RoundAwaitingWord RoundStatus = "awaiting_word"
	RoundActive       RoundStatus = "active"
	RoundCompleted    RoundStatus = "completed"
	RoundSkipped      RoundStatus = "skipped"
)

type Round struct {
	Base
	GroupID      string      `db:"group_id"`
	RoundDate    string      `db:"round_date"`
	PickerUserID string      `db:"picker_user_id"`
	WordHash     *string     `db:"word_hash"`
	WordPlain    *string     `db:"word_plain"`
	Status       RoundStatus `db:"status"`
}

const GUESS_TABLE = "guesses"

type Guess struct {
	Base
	RoundID      string `db:"round_id"`
	UserID       string `db:"user_id"`
	AttemptsUsed int    `db:"attempts_used"`
	Solved       bool   `db:"solved"`
	RowsJSON     string `db:"rows_json"`
	Score        int    `db:"score"`
	Finished     bool   `db:"finished"`
	FirstGuessAt *string `db:"first_guess_at"`
	CompletedAt  *string `db:"completed_at"`
}
