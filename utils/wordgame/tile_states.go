package wordgame

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/spf13/cast"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// TileStates is a JSON-encoded slice of per-letter tile feedback stored in SQLite
type TileStates []TileState

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MarshalJSON implements the json.Marshaler interface
func (t TileStates) MarshalJSON() ([]byte, error) {
	return json.Marshal([]TileState(t))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *TileStates) UnmarshalJSON(data []byte) error {
	var states []TileState
	if err := json.Unmarshal(data, &states); err != nil {
		return err
	}

	*t = TileStates(states)

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Value implements the driver.Valuer interface
func (t TileStates) Value() (driver.Value, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Scan implements the sql.Scanner interface
func (t *TileStates) Scan(value any) error {
	if value == nil {
		*t = nil

		return nil
	}

	raw := cast.ToString(value)
	if raw == "" {
		*t = nil

		return nil
	}

	if err := json.Unmarshal([]byte(raw), t); err != nil {
		return fmt.Errorf("invalid tile states json: %w", err)
	}

	return nil
}
