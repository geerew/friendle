package dao

import (
	"fmt"

	"github.com/geerew/friendle/models"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// roundMemberGuessColumns defines the columns to select for round member guess queries
var roundMemberGuessColumns = []string{
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_GUESS_TABLE_ID, models.BASE_ID),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_GUESS_TABLE_CREATED_AT, models.BASE_CREATED_AT),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_GUESS_TABLE_ROUND_ID, models.ROUND_MEMBER_GUESS_ROUND_ID),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_GUESS_TABLE_USER_ID, models.ROUND_MEMBER_GUESS_USER_ID),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_GUESS_TABLE_ATTEMPT, models.ROUND_MEMBER_GUESS_ATTEMPT),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_GUESS_TABLE_WORD, models.ROUND_MEMBER_GUESS_WORD),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_GUESS_TABLE_RESULT, models.ROUND_MEMBER_GUESS_RESULT),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_GUESS_TABLE_OUTCOME, models.ROUND_MEMBER_GUESS_OUTCOME),
}
