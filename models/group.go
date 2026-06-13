package models

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

	GROUP_MEMBER_COUNT = "member_count"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Group defines the model for a friend group
type Group struct {
	Base
	Name      string `db:"name"`       // Mutable
	CreatedBy string `db:"created_by"` // Immutable

	// Added via JOIN
	MemberCount int `db:"member_count"`
}
