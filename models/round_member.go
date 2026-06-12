package models

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	ROUND_MEMBER_TABLE = "round_members"

	ROUND_MEMBER_ROUND_ID       = "round_id"
	ROUND_MEMBER_USER_ID        = "user_id"
	ROUND_MEMBER_SOLVED         = "solved"
	ROUND_MEMBER_FINISHED       = "finished"
	ROUND_MEMBER_SCORE          = "score"
	ROUND_MEMBER_FIRST_GUESS_AT = "first_guess_at"
	ROUND_MEMBER_COMPLETED_AT   = "completed_at"

	ROUND_MEMBER_TABLE_ID             = ROUND_MEMBER_TABLE + "." + BASE_ID
	ROUND_MEMBER_TABLE_CREATED_AT     = ROUND_MEMBER_TABLE + "." + BASE_CREATED_AT
	ROUND_MEMBER_TABLE_UPDATED_AT     = ROUND_MEMBER_TABLE + "." + BASE_UPDATED_AT
	ROUND_MEMBER_TABLE_ROUND_ID       = ROUND_MEMBER_TABLE + "." + ROUND_MEMBER_ROUND_ID
	ROUND_MEMBER_TABLE_USER_ID        = ROUND_MEMBER_TABLE + "." + ROUND_MEMBER_USER_ID
	ROUND_MEMBER_TABLE_SOLVED         = ROUND_MEMBER_TABLE + "." + ROUND_MEMBER_SOLVED
	ROUND_MEMBER_TABLE_FINISHED       = ROUND_MEMBER_TABLE + "." + ROUND_MEMBER_FINISHED
	ROUND_MEMBER_TABLE_SCORE          = ROUND_MEMBER_TABLE + "." + ROUND_MEMBER_SCORE
	ROUND_MEMBER_TABLE_FIRST_GUESS_AT = ROUND_MEMBER_TABLE + "." + ROUND_MEMBER_FIRST_GUESS_AT
	ROUND_MEMBER_TABLE_COMPLETED_AT   = ROUND_MEMBER_TABLE + "." + ROUND_MEMBER_COMPLETED_AT
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundMember defines one member's participation in a single round
type RoundMember struct {
	Base
	RoundID      string  `db:"round_id"`       // Immutable
	UserID       string  `db:"user_id"`        // Immutable
	Solved       bool    `db:"solved"`         // Mutable
	Finished     bool    `db:"finished"`       // Mutable
	Score        int     `db:"score"`          // Mutable
	FirstGuessAt *string `db:"first_guess_at"` // Mutable
	CompletedAt  *string `db:"completed_at"`   // Mutable
}
