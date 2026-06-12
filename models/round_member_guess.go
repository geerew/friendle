package models

import (
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	ROUND_MEMBER_GUESS_TABLE = "round_member_guesses"

	ROUND_MEMBER_GUESS_ROUND_ID = "round_id"
	ROUND_MEMBER_GUESS_USER_ID  = "user_id"
	ROUND_MEMBER_GUESS_ATTEMPT  = "attempt"
	ROUND_MEMBER_GUESS_WORD     = "word"
	ROUND_MEMBER_GUESS_RESULT   = "result"
	ROUND_MEMBER_GUESS_OUTCOME  = "outcome"

	ROUND_MEMBER_GUESS_TABLE_ID         = ROUND_MEMBER_GUESS_TABLE + "." + BASE_ID
	ROUND_MEMBER_GUESS_TABLE_CREATED_AT = ROUND_MEMBER_GUESS_TABLE + "." + BASE_CREATED_AT
	ROUND_MEMBER_GUESS_TABLE_ROUND_ID   = ROUND_MEMBER_GUESS_TABLE + "." + ROUND_MEMBER_GUESS_ROUND_ID
	ROUND_MEMBER_GUESS_TABLE_USER_ID    = ROUND_MEMBER_GUESS_TABLE + "." + ROUND_MEMBER_GUESS_USER_ID
	ROUND_MEMBER_GUESS_TABLE_ATTEMPT    = ROUND_MEMBER_GUESS_TABLE + "." + ROUND_MEMBER_GUESS_ATTEMPT
	ROUND_MEMBER_GUESS_TABLE_WORD       = ROUND_MEMBER_GUESS_TABLE + "." + ROUND_MEMBER_GUESS_WORD
	ROUND_MEMBER_GUESS_TABLE_RESULT     = ROUND_MEMBER_GUESS_TABLE + "." + ROUND_MEMBER_GUESS_RESULT
	ROUND_MEMBER_GUESS_TABLE_OUTCOME    = ROUND_MEMBER_GUESS_TABLE + "." + ROUND_MEMBER_GUESS_OUTCOME
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundMemberGuess defines the model for a single guess attempt for a round member in a round
type RoundMemberGuess struct {
	Base
	RoundID string             `db:"round_id"` // Immutable
	UserID  string             `db:"user_id"`  // Immutable
	Attempt int                `db:"attempt"`  // Immutable
	Word    string             `db:"word"`     // Immutable
	Result  types.TileStates   `db:"result"`   // Immutable
	Outcome types.GuessOutcome `db:"outcome"`  // Immutable
}
