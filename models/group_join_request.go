package models

import (
	"fmt"

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
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupJoinRequestColumns returns the columns for use in a SELECT query
func GroupJoinRequestColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", JOIN_REQUEST_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", JOIN_REQUEST_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", JOIN_REQUEST_TABLE_UPDATED_AT, BASE_UPDATED_AT),
		fmt.Sprintf("%s AS %s", JOIN_REQUEST_TABLE_GROUP_ID, JOIN_REQUEST_GROUP_ID),
		fmt.Sprintf("%s AS %s", JOIN_REQUEST_TABLE_USER_ID, JOIN_REQUEST_USER_ID),
		fmt.Sprintf("%s AS %s", JOIN_REQUEST_TABLE_STATUS, JOIN_REQUEST_STATUS),
	}
}
