package models

import "fmt"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	GUESS_TABLE = "guesses"

	GUESS_ROUND_ID       = "round_id"
	GUESS_USER_ID        = "user_id"
	GUESS_ATTEMPTS_USED  = "attempts_used"
	GUESS_SOLVED         = "solved"
	GUESS_ROWS_JSON      = "rows_json"
	GUESS_SCORE          = "score"
	GUESS_FINISHED       = "finished"
	GUESS_FIRST_GUESS_AT = "first_guess_at"
	GUESS_COMPLETED_AT   = "completed_at"

	GUESS_TABLE_ID              = GUESS_TABLE + "." + BASE_ID
	GUESS_TABLE_CREATED_AT      = GUESS_TABLE + "." + BASE_CREATED_AT
	GUESS_TABLE_UPDATED_AT      = GUESS_TABLE + "." + BASE_UPDATED_AT
	GUESS_TABLE_ROUND_ID        = GUESS_TABLE + "." + GUESS_ROUND_ID
	GUESS_TABLE_USER_ID         = GUESS_TABLE + "." + GUESS_USER_ID
	GUESS_TABLE_ATTEMPTS_USED   = GUESS_TABLE + "." + GUESS_ATTEMPTS_USED
	GUESS_TABLE_SOLVED          = GUESS_TABLE + "." + GUESS_SOLVED
	GUESS_TABLE_ROWS_JSON       = GUESS_TABLE + "." + GUESS_ROWS_JSON
	GUESS_TABLE_SCORE           = GUESS_TABLE + "." + GUESS_SCORE
	GUESS_TABLE_FINISHED        = GUESS_TABLE + "." + GUESS_FINISHED
	GUESS_TABLE_FIRST_GUESS_AT  = GUESS_TABLE + "." + GUESS_FIRST_GUESS_AT
	GUESS_TABLE_COMPLETED_AT    = GUESS_TABLE + "." + GUESS_COMPLETED_AT
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Guess defines the model for a user's guesses in a round
type Guess struct {
	Base
	RoundID      string  `db:"round_id"`       // Immutable
	UserID       string  `db:"user_id"`        // Immutable
	AttemptsUsed int     `db:"attempts_used"`  // Mutable
	Solved       bool    `db:"solved"`         // Mutable
	RowsJSON     string  `db:"rows_json"`      // Mutable
	Score        int     `db:"score"`          // Mutable
	Finished     bool    `db:"finished"`       // Mutable
	FirstGuessAt *string `db:"first_guess_at"` // Mutable
	CompletedAt  *string `db:"completed_at"`   // Mutable
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GuessColumns returns the columns for use in a SELECT query
func GuessColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", GUESS_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_UPDATED_AT, BASE_UPDATED_AT),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_ROUND_ID, GUESS_ROUND_ID),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_USER_ID, GUESS_USER_ID),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_ATTEMPTS_USED, GUESS_ATTEMPTS_USED),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_SOLVED, GUESS_SOLVED),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_ROWS_JSON, GUESS_ROWS_JSON),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_SCORE, GUESS_SCORE),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_FINISHED, GUESS_FINISHED),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_FIRST_GUESS_AT, GUESS_FIRST_GUESS_AT),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_COMPLETED_AT, GUESS_COMPLETED_AT),
	}
}
