package models

import (
	"fmt"

	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	GROUP_MEMBER_TABLE = "group_members"

	GROUP_MEMBER_GROUP_ID     = "group_id"
	GROUP_MEMBER_USER_ID      = "user_id"
	GROUP_MEMBER_GROUP_ROLE   = "group_role"
	GROUP_MEMBER_TIMES_PICKED = "times_picked"
	GROUP_MEMBER_PICKER_SKIPS = "picker_skips"

	GROUP_MEMBER_TABLE_ID           = GROUP_MEMBER_TABLE + "." + BASE_ID
	GROUP_MEMBER_TABLE_CREATED_AT   = GROUP_MEMBER_TABLE + "." + BASE_CREATED_AT
	GROUP_MEMBER_TABLE_UPDATED_AT   = GROUP_MEMBER_TABLE + "." + BASE_UPDATED_AT
	GROUP_MEMBER_TABLE_GROUP_ID     = GROUP_MEMBER_TABLE + "." + GROUP_MEMBER_GROUP_ID
	GROUP_MEMBER_TABLE_USER_ID      = GROUP_MEMBER_TABLE + "." + GROUP_MEMBER_USER_ID
	GROUP_MEMBER_TABLE_GROUP_ROLE   = GROUP_MEMBER_TABLE + "." + GROUP_MEMBER_GROUP_ROLE
	GROUP_MEMBER_TABLE_TIMES_PICKED = GROUP_MEMBER_TABLE + "." + GROUP_MEMBER_TIMES_PICKED
	GROUP_MEMBER_TABLE_PICKER_SKIPS = GROUP_MEMBER_TABLE + "." + GROUP_MEMBER_PICKER_SKIPS
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupMember defines the model for a group membership row
type GroupMember struct {
	Base
	GroupID     string          `db:"group_id"`     // Immutable
	UserID      string          `db:"user_id"`      // Immutable
	GroupRole   types.GroupRole `db:"group_role"`   // Mutable
	TimesPicked int             `db:"times_picked"` // Mutable
	PickerSkips int             `db:"picker_skips"` // Mutable

	// Relations
	User  *User  `db:"-"`
	Group *Group `db:"-"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupMemberColumns returns the columns for use in a SELECT query
func GroupMemberColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", GROUP_MEMBER_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", GROUP_MEMBER_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", GROUP_MEMBER_TABLE_UPDATED_AT, BASE_UPDATED_AT),
		fmt.Sprintf("%s AS %s", GROUP_MEMBER_TABLE_GROUP_ID, GROUP_MEMBER_GROUP_ID),
		fmt.Sprintf("%s AS %s", GROUP_MEMBER_TABLE_USER_ID, GROUP_MEMBER_USER_ID),
		fmt.Sprintf("%s AS %s", GROUP_MEMBER_TABLE_GROUP_ROLE, GROUP_MEMBER_GROUP_ROLE),
		fmt.Sprintf("%s AS %s", GROUP_MEMBER_TABLE_TIMES_PICKED, GROUP_MEMBER_TIMES_PICKED),
		fmt.Sprintf("%s AS %s", GROUP_MEMBER_TABLE_PICKER_SKIPS, GROUP_MEMBER_PICKER_SKIPS),
	}
}
