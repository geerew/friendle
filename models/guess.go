package models

import (
	"fmt"

	"github.com/geerew/friendle/utils/types"
	"github.com/geerew/friendle/utils/wordgame"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	GUESS_TABLE = "guesses"

	GUESS_ROUND_ID = "round_id"
	GUESS_USER_ID  = "user_id"
	GUESS_ATTEMPT  = "attempt"
	GUESS_WORD     = "word"
	GUESS_RESULT   = "result"
	GUESS_OUTCOME  = "outcome"

	GUESS_TABLE_ID         = GUESS_TABLE + "." + BASE_ID
	GUESS_TABLE_CREATED_AT = GUESS_TABLE + "." + BASE_CREATED_AT
	GUESS_TABLE_ROUND_ID   = GUESS_TABLE + "." + GUESS_ROUND_ID
	GUESS_TABLE_USER_ID    = GUESS_TABLE + "." + GUESS_USER_ID
	GUESS_TABLE_ATTEMPT    = GUESS_TABLE + "." + GUESS_ATTEMPT
	GUESS_TABLE_WORD       = GUESS_TABLE + "." + GUESS_WORD
	GUESS_TABLE_RESULT     = GUESS_TABLE + "." + GUESS_RESULT
	GUESS_TABLE_OUTCOME    = GUESS_TABLE + "." + GUESS_OUTCOME
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Guess defines the model for a single guess attempt in a round
type Guess struct {
	Base
	RoundID string              `db:"round_id"` // Immutable
	UserID  string              `db:"user_id"`  // Immutable
	Attempt int                 `db:"attempt"`  // Immutable
	Word    string              `db:"word"`     // Immutable
	Result  wordgame.TileStates `db:"result"`   // Immutable
	Outcome types.GuessOutcome  `db:"outcome"`  // Immutable
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GuessColumns returns the columns for use in a SELECT query
func GuessColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", GUESS_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_ROUND_ID, GUESS_ROUND_ID),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_USER_ID, GUESS_USER_ID),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_ATTEMPT, GUESS_ATTEMPT),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_WORD, GUESS_WORD),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_RESULT, GUESS_RESULT),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_OUTCOME, GUESS_OUTCOME),
	}
}
