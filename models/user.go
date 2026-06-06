package models

import "github.com/geerew/friendle/utils/types"

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
