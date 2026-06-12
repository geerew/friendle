package models

import (
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	JOIN_REQUEST_TABLE = "group_join_requests"

	JOIN_REQUEST_GROUP_ID = "group_id"
	JOIN_REQUEST_USER_ID  = "user_id"
	JOIN_REQUEST_STATUS   = "status"

	JOIN_REQUEST_TABLE_ID         = JOIN_REQUEST_TABLE + "." + BASE_ID
	JOIN_REQUEST_TABLE_CREATED_AT = JOIN_REQUEST_TABLE + "." + BASE_CREATED_AT
	JOIN_REQUEST_TABLE_UPDATED_AT = JOIN_REQUEST_TABLE + "." + BASE_UPDATED_AT
	JOIN_REQUEST_TABLE_GROUP_ID   = JOIN_REQUEST_TABLE + "." + JOIN_REQUEST_GROUP_ID
	JOIN_REQUEST_TABLE_USER_ID    = JOIN_REQUEST_TABLE + "." + JOIN_REQUEST_USER_ID
	JOIN_REQUEST_TABLE_STATUS     = JOIN_REQUEST_TABLE + "." + JOIN_REQUEST_STATUS
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupJoinRequest defines the model for a pending or resolved join request
type GroupJoinRequest struct {
	Base
	GroupID string                  `db:"group_id"` // Immutable
	UserID  string                  `db:"user_id"`  // Immutable
	Status  types.JoinRequestStatus `db:"status"`   // Mutable

	// Added via JOIN
	DisplayName string `db:"display_name"`
}
