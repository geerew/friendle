package models

import (
	"fmt"

	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	USER_TABLE = "users"

	USER_USERNAME      = "username"
	USER_DISPLAY_NAME  = "display_name"
	USER_PASSWORD_HASH = "password_hash"
	USER_SITE_ROLE     = "site_role"

	USER_TABLE_ID            = USER_TABLE + "." + BASE_ID
	USER_TABLE_CREATED_AT    = USER_TABLE + "." + BASE_CREATED_AT
	USER_TABLE_UPDATED_AT    = USER_TABLE + "." + BASE_UPDATED_AT
	USER_TABLE_USERNAME      = USER_TABLE + "." + USER_USERNAME
	USER_TABLE_DISPLAY_NAME  = USER_TABLE + "." + USER_DISPLAY_NAME
	USER_TABLE_PASSWORD_HASH = USER_TABLE + "." + USER_PASSWORD_HASH
	USER_TABLE_SITE_ROLE     = USER_TABLE + "." + USER_SITE_ROLE
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// User defines the model for a user account
type User struct {
	Base
	Username     string         `db:"username"`      // Immutable
	DisplayName  string         `db:"display_name"`  // Mutable
	PasswordHash string         `db:"password_hash"` // Mutable
	SiteRole     types.SiteRole `db:"site_role"`     // Mutable
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UserColumns returns the columns for use in a SELECT query
func UserColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", USER_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", USER_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", USER_TABLE_UPDATED_AT, BASE_UPDATED_AT),
		fmt.Sprintf("%s AS %s", USER_TABLE_USERNAME, USER_USERNAME),
		fmt.Sprintf("%s AS %s", USER_TABLE_DISPLAY_NAME, USER_DISPLAY_NAME),
		fmt.Sprintf("%s AS %s", USER_TABLE_PASSWORD_HASH, USER_PASSWORD_HASH),
		fmt.Sprintf("%s AS %s", USER_TABLE_SITE_ROLE, USER_SITE_ROLE),
	}
}
