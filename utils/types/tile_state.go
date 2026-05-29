package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TileState is the per-letter feedback for one position in a guess row
type TileState string

const (
	TileCorrect TileState = "correct"
	TilePresent TileState = "present"
	TileAbsent  TileState = "absent"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsValid reports whether the tile state is a known value
func (t TileState) IsValid() bool {
	switch t {
	case TileCorrect, TilePresent, TileAbsent:
		return true
	}

	return false
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// String returns the tile state as a string
func (t TileState) String() string {
	return string(t)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MarshalJSON implements the json.Marshaler interface
func (t TileState) MarshalJSON() ([]byte, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("invalid tile state: %s", t)
	}

	return json.Marshal(string(t))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *TileState) UnmarshalJSON(data []byte) error {
	var state string
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	*t = TileState(state)
	if !t.IsValid() {
		return fmt.Errorf("invalid tile state: %s", state)
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Value implements the driver.Valuer interface
func (t TileState) Value() (driver.Value, error) {
	if !t.IsValid() {
		return nil, fmt.Errorf("invalid tile state: %s", t)
	}

	return string(t), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Scan implements the sql.Scanner interface
func (t *TileState) Scan(value any) error {
	state, ok := value.(string)
	if !ok {
		return errors.New("invalid data type for TileState")
	}

	*t = TileState(state)
	if !t.IsValid() {
		return fmt.Errorf("invalid tile state: %s", state)
	}

	return nil
}
