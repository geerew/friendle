package models

import (
	"fmt"

	"github.com/geerew/friendle/utils/security"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	GUESS_TABLE = "guesses"

	GUESS_ROUND_ID       = "round_id"
	GUESS_USER_ID        = "user_id"
	GUESS_ATTEMPT_NUMBER = "attempt_number"
	GUESS_WORD           = "word"
	GUESS_RESULT         = "result"

	GUESS_TABLE_ID         = GUESS_TABLE + "." + BASE_ID
	GUESS_TABLE_CREATED_AT = GUESS_TABLE + "." + BASE_CREATED_AT
	GUESS_TABLE_ROUND_ID   = GUESS_TABLE + "." + GUESS_ROUND_ID
	GUESS_TABLE_USER_ID    = GUESS_TABLE + "." + GUESS_USER_ID
	GUESS_TABLE_ATTEMPT    = GUESS_TABLE + "." + GUESS_ATTEMPT_NUMBER
	GUESS_TABLE_WORD       = GUESS_TABLE + "." + GUESS_WORD
	GUESS_TABLE_RESULT     = GUESS_TABLE + "." + GUESS_RESULT
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Guess defines the model for a single guess attempt in a round
type Guess struct {
	ID            string         `db:"id"`
	RoundID       string         `db:"round_id"`       // Immutable
	UserID        string         `db:"user_id"`        // Immutable
	AttemptNumber int            `db:"attempt_number"` // Immutable
	Word          string         `db:"word"`           // Immutable
	Result        string         `db:"result"`         // Immutable JSON tile states
	CreatedAt     types.DateTime `db:"created_at"`     // Immutable
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RefreshId assigns a new ID when empty
func (g *Guess) RefreshId() {
	if g.ID == "" {
		g.ID = security.PseudorandomString(10)
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RefreshCreatedAt sets created_at to now when zero
func (g *Guess) RefreshCreatedAt() {
	if g.CreatedAt.IsZero() {
		g.CreatedAt = types.NowDateTime()
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GuessColumns returns the columns for use in a SELECT query
func GuessColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", GUESS_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_ROUND_ID, GUESS_ROUND_ID),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_USER_ID, GUESS_USER_ID),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_ATTEMPT, GUESS_ATTEMPT_NUMBER),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_WORD, GUESS_WORD),
		fmt.Sprintf("%s AS %s", GUESS_TABLE_RESULT, GUESS_RESULT),
	}
}
