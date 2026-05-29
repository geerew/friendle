package dao

import "github.com/geerew/friendle/utils/types"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AdminUserRow is a user projection for the site admin list query
type AdminUserRow struct {
	ID          string         `db:"id"`
	Username    string         `db:"username"`
	DisplayName string         `db:"display_name"`
	SiteRole    types.SiteRole `db:"site_role"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AdminGroupRow is a group projection for the site admin list query
type AdminGroupRow struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	MemberCount int    `db:"member_count"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupSearchRow is a group projection for name search results
type GroupSearchRow struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	MemberCount int    `db:"member_count"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UserGroupSummaryRow is a group projection for a user's memberships
type UserGroupSummaryRow struct {
	UserID      string          `db:"user_id"`
	ID          string          `db:"id"`
	Name        string          `db:"name"`
	MemberCount int             `db:"member_count"`
	GroupRole   types.GroupRole `db:"group_role"`
}
