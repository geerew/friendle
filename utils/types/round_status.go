package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundStatus is the lifecycle state of a daily group round
type RoundStatus string

const (
	RoundAwaitingWord RoundStatus = "awaiting_word"
	RoundActive       RoundStatus = "active"
	RoundCompleted    RoundStatus = "completed"
	RoundSkipped      RoundStatus = "skipped"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsValid reports whether the round status is a known value
func (s RoundStatus) IsValid() bool {
	switch s {
	case RoundAwaitingWord, RoundActive, RoundCompleted, RoundSkipped:
		return true
	}

	return false
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// String returns the round status as a string
func (s RoundStatus) String() string {
	return string(s)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MarshalJSON implements the json.Marshaler interface
func (s RoundStatus) MarshalJSON() ([]byte, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid round status: %s", s)
	}

	return json.Marshal(string(s))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UnmarshalJSON implements the json.Unmarshaler interface
func (s *RoundStatus) UnmarshalJSON(data []byte) error {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}

	*s = RoundStatus(status)
	if !s.IsValid() {
		return fmt.Errorf("invalid round status: %s", status)
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Value implements the driver.Valuer interface
func (s RoundStatus) Value() (driver.Value, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid round status: %s", s)
	}

	return string(s), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Scan implements the sql.Scanner interface
func (s *RoundStatus) Scan(value interface{}) error {
	status, ok := value.(string)
	if !ok {
		return errors.New("invalid data type for RoundStatus")
	}

	*s = RoundStatus(status)
	if !s.IsValid() {
		return fmt.Errorf("invalid round status: %s", status)
	}

	return nil
}
