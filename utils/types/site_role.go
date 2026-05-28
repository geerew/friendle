package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type SiteRole string

const (
	SiteRoleAdmin SiteRole = "site_admin"
	SiteRoleUser  SiteRole = "site_user"
)

type GroupRole string

const (
	GroupRoleAdmin GroupRole = "group_admin"
	GroupRoleUser  GroupRole = "group_user"
)

type Principal struct {
	UserID   string
	SiteRole SiteRole
	Role     SiteRole // alias for SiteRole
}

func NewSiteRole(role string) SiteRole {
	switch role {
	case "site_admin", "admin":
		return SiteRoleAdmin
	case "site_user", "user":
		return SiteRoleUser
	default:
		return SiteRoleUser
	}
}

func (r SiteRole) IsValid() bool {
	return r == SiteRoleAdmin || r == SiteRoleUser
}

func (r SiteRole) String() string { return string(r) }

func (r SiteRole) MarshalJSON() ([]byte, error) {
	if !r.IsValid() {
		return nil, fmt.Errorf("invalid site role: %s", r)
	}
	return json.Marshal(string(r))
}

func (r *SiteRole) UnmarshalJSON(data []byte) error {
	var role string
	if err := json.Unmarshal(data, &role); err != nil {
		return err
	}
	*r = NewSiteRole(role)
	return nil
}

func (r SiteRole) Value() (driver.Value, error) {
	if !r.IsValid() {
		return nil, fmt.Errorf("invalid site role: %s", r)
	}
	return string(r), nil
}

func (r *SiteRole) Scan(value interface{}) error {
	role, ok := value.(string)
	if !ok {
		return errors.New("invalid data type for SiteRole")
	}
	*r = NewSiteRole(role)
	return nil
}

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

func (r GroupRole) IsValid() bool {
	return r == GroupRoleAdmin || r == GroupRoleUser
}

func (r GroupRole) String() string { return string(r) }

// Back-compat aliases used during offcourse port
type UserRole = SiteRole

const (
	UserRoleAdmin = SiteRoleAdmin
	UserRoleUser  = SiteRoleUser
)

func NewUserRole(role string) SiteRole { return NewSiteRole(role) }
