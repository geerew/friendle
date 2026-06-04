package models

import (
	"fmt"

	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const (
	ROUND_TABLE = "rounds"

	ROUND_GROUP_ID       = "group_id"
	ROUND_ROUND_DATE     = "round_date"
	ROUND_PICKER_USER_ID = "picker_user_id"
	ROUND_WORD_PLAIN     = "word_plain"
	ROUND_STATUS         = "status"

	ROUND_TABLE_ID             = ROUND_TABLE + "." + BASE_ID
	ROUND_TABLE_CREATED_AT     = ROUND_TABLE + "." + BASE_CREATED_AT
	ROUND_TABLE_UPDATED_AT     = ROUND_TABLE + "." + BASE_UPDATED_AT
	ROUND_TABLE_GROUP_ID       = ROUND_TABLE + "." + ROUND_GROUP_ID
	ROUND_TABLE_ROUND_DATE     = ROUND_TABLE + "." + ROUND_ROUND_DATE
	ROUND_TABLE_PICKER_USER_ID = ROUND_TABLE + "." + ROUND_PICKER_USER_ID
	ROUND_TABLE_WORD_PLAIN     = ROUND_TABLE + "." + ROUND_WORD_PLAIN
	ROUND_TABLE_STATUS         = ROUND_TABLE + "." + ROUND_STATUS
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Round defines the model for a daily word round in a group
type Round struct {
	Base
	GroupID      string            `db:"group_id"`       // Immutable
	RoundDate    string            `db:"round_date"`     // Immutable
	PickerUserID string            `db:"picker_user_id"` // Immutable
	WordPlain    *string           `db:"word_plain"`     // Mutable
	Status       types.RoundStatus `db:"status"`         // Mutable
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RoundColumns returns the columns for use in a SELECT query
func RoundColumns() []string {
	return []string{
		fmt.Sprintf("%s AS %s", ROUND_TABLE_ID, BASE_ID),
		fmt.Sprintf("%s AS %s", ROUND_TABLE_CREATED_AT, BASE_CREATED_AT),
		fmt.Sprintf("%s AS %s", ROUND_TABLE_UPDATED_AT, BASE_UPDATED_AT),
		fmt.Sprintf("%s AS %s", ROUND_TABLE_GROUP_ID, ROUND_GROUP_ID),
		fmt.Sprintf("%s AS %s", ROUND_TABLE_ROUND_DATE, ROUND_ROUND_DATE),
		fmt.Sprintf("%s AS %s", ROUND_TABLE_PICKER_USER_ID, ROUND_PICKER_USER_ID),
		fmt.Sprintf("%s AS %s", ROUND_TABLE_WORD_PLAIN, ROUND_WORD_PLAIN),
		fmt.Sprintf("%s AS %s", ROUND_TABLE_STATUS, ROUND_STATUS),
	}
}
