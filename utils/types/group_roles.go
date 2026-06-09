package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupRole is the role a user holds within a group
type GroupRole string

const (
	GroupRoleAdmin GroupRole = "group_admin"
	GroupRoleUser  GroupRole = "group_user"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewGroupRole parses a group role string, defaulting to group user when unknown
func NewGroupRole(role string) GroupRole {
	switch role {
	case "group_admin":
		return GroupRoleAdmin
	case "group_user":
		return GroupRoleUser
	default:
		return GroupRoleUser
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsValid reports whether the group role is a known value
func (r GroupRole) IsValid() bool {
	return r == GroupRoleAdmin || r == GroupRoleUser
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// String returns the group role as a string
func (r GroupRole) String() string {
	return string(r)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MarshalJSON implements the json.Marshaler interface
func (r GroupRole) MarshalJSON() ([]byte, error) {
	if !r.IsValid() {
		return nil, fmt.Errorf("invalid group role: %s", r)
	}

	return json.Marshal(string(r))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UnmarshalJSON implements the json.Unmarshaler interface
func (r *GroupRole) UnmarshalJSON(data []byte) error {
	var role string
	if err := json.Unmarshal(data, &role); err != nil {
		return err
	}

	*r = NewGroupRole(role)

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Value implements the driver.Valuer interface
func (r GroupRole) Value() (driver.Value, error) {
	if !r.IsValid() {
		return nil, fmt.Errorf("invalid group role: %s", r)
	}

	return string(r), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Scan implements the sql.Scanner interface
func (r *GroupRole) Scan(value interface{}) error {
	role, ok := value.(string)
	if !ok {
		return errors.New("invalid data type for GroupRole")
	}

	*r = NewGroupRole(role)

	return nil
}
