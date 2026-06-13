package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// JoinRequestStatus is the approval state of a group join request
type JoinRequestStatus string

const (
	JoinPending  JoinRequestStatus = "pending"
	JoinApproved JoinRequestStatus = "approved"
	JoinRejected JoinRequestStatus = "rejected"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsValid reports whether the join request status is a known value
func (s JoinRequestStatus) IsValid() bool {
	switch s {
	case JoinPending, JoinApproved, JoinRejected:
		return true
	}

	return false
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// String returns the join request status as a string
func (s JoinRequestStatus) String() string {
	return string(s)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MarshalJSON implements the json.Marshaler interface
func (s JoinRequestStatus) MarshalJSON() ([]byte, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid join request status: %s", s)
	}

	return json.Marshal(string(s))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UnmarshalJSON implements the json.Unmarshaler interface
func (s *JoinRequestStatus) UnmarshalJSON(data []byte) error {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}

	*s = JoinRequestStatus(status)
	if !s.IsValid() {
		return fmt.Errorf("invalid join request status: %s", status)
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Value implements the driver.Valuer interface
func (s JoinRequestStatus) Value() (driver.Value, error) {
	if !s.IsValid() {
		return nil, fmt.Errorf("invalid join request status: %s", s)
	}

	return string(s), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Scan implements the sql.Scanner interface
func (s *JoinRequestStatus) Scan(value interface{}) error {
	status, ok := value.(string)
	if !ok {
		return errors.New("invalid data type for JoinRequestStatus")
	}

	*s = JoinRequestStatus(status)
	if !s.IsValid() {
		return fmt.Errorf("invalid join request status: %s", status)
	}

	return nil
}
