package models

import "fmt"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	GROUP_TABLE = "groups"

	GROUP_NAME       = "name"
	GROUP_CREATED_BY = "created_by"

	GROUP_TABLE_ID         = GROUP_TABLE + "." + BASE_ID
	GROUP_TABLE_CREATED_AT = GROUP_TABLE + "." + BASE_CREATED_AT
	GROUP_TABLE_UPDATED_AT = GROUP_TABLE + "." + BASE_UPDATED_AT
	GROUP_TABLE_NAME       = GROUP_TABLE + "." + GROUP_NAME
	GROUP_TABLE_CREATED_BY = GROUP_TABLE + "." + GROUP_CREATED_BY
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Group defines the model for a friend group
type Group struct {
	Base
	Name      string `db:"name"`       // Mutable
	CreatedBy string `db:"created_by"` // Immutable

	// Relations
	Members      []*GroupMember      `db:"-"`
	JoinRequests []*GroupJoinRequest `db:"-"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupColumns returns the columns for use in a SELECT query
func GroupColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", GROUP_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", GROUP_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", GROUP_TABLE_UPDATED_AT, BASE_UPDATED_AT),
		fmt.Sprintf("%s AS %s", GROUP_TABLE_NAME, GROUP_NAME),
		fmt.Sprintf("%s AS %s", GROUP_TABLE_CREATED_BY, GROUP_CREATED_BY),
	}
}
