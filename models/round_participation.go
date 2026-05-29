package models

import "fmt"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	ROUND_PARTICIPATION_TABLE = "round_participations"

	ROUND_PARTICIPATION_ROUND_ID       = "round_id"
	ROUND_PARTICIPATION_USER_ID        = "user_id"
	ROUND_PARTICIPATION_SOLVED         = "solved"
	ROUND_PARTICIPATION_FINISHED       = "finished"
	ROUND_PARTICIPATION_SCORE          = "score"
	ROUND_PARTICIPATION_FIRST_GUESS_AT = "first_guess_at"
	ROUND_PARTICIPATION_COMPLETED_AT   = "completed_at"

	ROUND_PARTICIPATION_TABLE_ID              = ROUND_PARTICIPATION_TABLE + "." + BASE_ID
	ROUND_PARTICIPATION_TABLE_CREATED_AT      = ROUND_PARTICIPATION_TABLE + "." + BASE_CREATED_AT
	ROUND_PARTICIPATION_TABLE_UPDATED_AT      = ROUND_PARTICIPATION_TABLE + "." + BASE_UPDATED_AT
	ROUND_PARTICIPATION_TABLE_ROUND_ID        = ROUND_PARTICIPATION_TABLE + "." + ROUND_PARTICIPATION_ROUND_ID
	ROUND_PARTICIPATION_TABLE_USER_ID         = ROUND_PARTICIPATION_TABLE + "." + ROUND_PARTICIPATION_USER_ID
	ROUND_PARTICIPATION_TABLE_SOLVED          = ROUND_PARTICIPATION_TABLE + "." + ROUND_PARTICIPATION_SOLVED
	ROUND_PARTICIPATION_TABLE_FINISHED        = ROUND_PARTICIPATION_TABLE + "." + ROUND_PARTICIPATION_FINISHED
	ROUND_PARTICIPATION_TABLE_SCORE           = ROUND_PARTICIPATION_TABLE + "." + ROUND_PARTICIPATION_SCORE
	ROUND_PARTICIPATION_TABLE_FIRST_GUESS_AT  = ROUND_PARTICIPATION_TABLE + "." + ROUND_PARTICIPATION_FIRST_GUESS_AT
	ROUND_PARTICIPATION_TABLE_COMPLETED_AT    = ROUND_PARTICIPATION_TABLE + "." + ROUND_PARTICIPATION_COMPLETED_AT
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundParticipation defines a user's outcome in a single round
type RoundParticipation struct {
	Base
	RoundID      string  `db:"round_id"`       // Immutable
	UserID       string  `db:"user_id"`        // Immutable
	Solved       bool    `db:"solved"`         // Mutable
	Finished     bool    `db:"finished"`       // Mutable
	Score        int     `db:"score"`          // Mutable
	FirstGuessAt *string `db:"first_guess_at"` // Mutable
	CompletedAt  *string `db:"completed_at"`   // Mutable

	Guesses []*Guess `db:"-"`
	User    *User    `db:"-"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundParticipationColumns returns the columns for use in a SELECT query
func RoundParticipationColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", ROUND_PARTICIPATION_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", ROUND_PARTICIPATION_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", ROUND_PARTICIPATION_TABLE_UPDATED_AT, BASE_UPDATED_AT),
		fmt.Sprintf("%s AS %s", ROUND_PARTICIPATION_TABLE_ROUND_ID, ROUND_PARTICIPATION_ROUND_ID),
		fmt.Sprintf("%s AS %s", ROUND_PARTICIPATION_TABLE_USER_ID, ROUND_PARTICIPATION_USER_ID),
		fmt.Sprintf("%s AS %s", ROUND_PARTICIPATION_TABLE_SOLVED, ROUND_PARTICIPATION_SOLVED),
		fmt.Sprintf("%s AS %s", ROUND_PARTICIPATION_TABLE_FINISHED, ROUND_PARTICIPATION_FINISHED),
		fmt.Sprintf("%s AS %s", ROUND_PARTICIPATION_TABLE_SCORE, ROUND_PARTICIPATION_SCORE),
		fmt.Sprintf("%s AS %s", ROUND_PARTICIPATION_TABLE_FIRST_GUESS_AT, ROUND_PARTICIPATION_FIRST_GUESS_AT),
		fmt.Sprintf("%s AS %s", ROUND_PARTICIPATION_TABLE_COMPLETED_AT, ROUND_PARTICIPATION_COMPLETED_AT),
	}
}
