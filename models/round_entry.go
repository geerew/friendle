package models

import "fmt"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	ROUND_ENTRY_TABLE = "round_entries"

	ROUND_ENTRY_ROUND_ID       = "round_id"
	ROUND_ENTRY_USER_ID        = "user_id"
	ROUND_ENTRY_SOLVED         = "solved"
	ROUND_ENTRY_FINISHED       = "finished"
	ROUND_ENTRY_SCORE          = "score"
	ROUND_ENTRY_FIRST_GUESS_AT = "first_guess_at"
	ROUND_ENTRY_COMPLETED_AT   = "completed_at"

	ROUND_ENTRY_TABLE_ID             = ROUND_ENTRY_TABLE + "." + BASE_ID
	ROUND_ENTRY_TABLE_CREATED_AT     = ROUND_ENTRY_TABLE + "." + BASE_CREATED_AT
	ROUND_ENTRY_TABLE_UPDATED_AT     = ROUND_ENTRY_TABLE + "." + BASE_UPDATED_AT
	ROUND_ENTRY_TABLE_ROUND_ID       = ROUND_ENTRY_TABLE + "." + ROUND_ENTRY_ROUND_ID
	ROUND_ENTRY_TABLE_USER_ID        = ROUND_ENTRY_TABLE + "." + ROUND_ENTRY_USER_ID
	ROUND_ENTRY_TABLE_SOLVED         = ROUND_ENTRY_TABLE + "." + ROUND_ENTRY_SOLVED
	ROUND_ENTRY_TABLE_FINISHED       = ROUND_ENTRY_TABLE + "." + ROUND_ENTRY_FINISHED
	ROUND_ENTRY_TABLE_SCORE          = ROUND_ENTRY_TABLE + "." + ROUND_ENTRY_SCORE
	ROUND_ENTRY_TABLE_FIRST_GUESS_AT = ROUND_ENTRY_TABLE + "." + ROUND_ENTRY_FIRST_GUESS_AT
	ROUND_ENTRY_TABLE_COMPLETED_AT   = ROUND_ENTRY_TABLE + "." + ROUND_ENTRY_COMPLETED_AT
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundEntry defines one guesser's summary for a single round
type RoundEntry struct {
	Base
	RoundID      string  `db:"round_id"`       // Immutable
	UserID       string  `db:"user_id"`        // Immutable
	Solved       bool    `db:"solved"`         // Mutable
	Finished     bool    `db:"finished"`       // Mutable
	Score        int     `db:"score"`          // Mutable
	FirstGuessAt *string `db:"first_guess_at"` // Mutable
	CompletedAt  *string `db:"completed_at"`   // Mutable
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundEntryColumns returns the columns for use in a SELECT query
func RoundEntryColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", ROUND_ENTRY_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", ROUND_ENTRY_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", ROUND_ENTRY_TABLE_UPDATED_AT, BASE_UPDATED_AT),
		fmt.Sprintf("%s AS %s", ROUND_ENTRY_TABLE_ROUND_ID, ROUND_ENTRY_ROUND_ID),
		fmt.Sprintf("%s AS %s", ROUND_ENTRY_TABLE_USER_ID, ROUND_ENTRY_USER_ID),
		fmt.Sprintf("%s AS %s", ROUND_ENTRY_TABLE_SOLVED, ROUND_ENTRY_SOLVED),
		fmt.Sprintf("%s AS %s", ROUND_ENTRY_TABLE_FINISHED, ROUND_ENTRY_FINISHED),
		fmt.Sprintf("%s AS %s", ROUND_ENTRY_TABLE_SCORE, ROUND_ENTRY_SCORE),
		fmt.Sprintf("%s AS %s", ROUND_ENTRY_TABLE_FIRST_GUESS_AT, ROUND_ENTRY_FIRST_GUESS_AT),
		fmt.Sprintf("%s AS %s", ROUND_ENTRY_TABLE_COMPLETED_AT, ROUND_ENTRY_COMPLETED_AT),
	}
}
