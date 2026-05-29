package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GuessOutcome is the overall result of a single guess attempt
type GuessOutcome string

const (
	GuessOutcomeCorrect   GuessOutcome = "correct"
	GuessOutcomePartial   GuessOutcome = "partial"
	GuessOutcomeIncorrect GuessOutcome = "incorrect"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsValid reports whether the guess outcome is a known value
func (o GuessOutcome) IsValid() bool {
	switch o {
	case GuessOutcomeCorrect, GuessOutcomePartial, GuessOutcomeIncorrect:
		return true
	}

	return false
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// String returns the guess outcome as a string
func (o GuessOutcome) String() string {
	return string(o)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MarshalJSON implements the json.Marshaler interface
func (o GuessOutcome) MarshalJSON() ([]byte, error) {
	if !o.IsValid() {
		return nil, fmt.Errorf("invalid guess outcome: %s", o)
	}

	return json.Marshal(string(o))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UnmarshalJSON implements the json.Unmarshaler interface
func (o *GuessOutcome) UnmarshalJSON(data []byte) error {
	var outcome string
	if err := json.Unmarshal(data, &outcome); err != nil {
		return err
	}

	*o = GuessOutcome(outcome)
	if !o.IsValid() {
		return fmt.Errorf("invalid guess outcome: %s", outcome)
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Value implements the driver.Valuer interface
func (o GuessOutcome) Value() (driver.Value, error) {
	if !o.IsValid() {
		return nil, fmt.Errorf("invalid guess outcome: %s", o)
	}

	return string(o), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Scan implements the sql.Scanner interface
func (o *GuessOutcome) Scan(value interface{}) error {
	outcome, ok := value.(string)
	if !ok {
		return errors.New("invalid data type for GuessOutcome")
	}

	*o = GuessOutcome(outcome)
	if !o.IsValid() {
		return fmt.Errorf("invalid guess outcome: %s", outcome)
	}

	return nil
}
