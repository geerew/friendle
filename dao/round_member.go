package dao

import (
	"context"
	"fmt"

	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// roundMemberColumns defines the columns to select for round member queries
var roundMemberColumns = []string{
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_TABLE_ID, models.BASE_ID),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_TABLE_CREATED_AT, models.BASE_CREATED_AT),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_TABLE_UPDATED_AT, models.BASE_UPDATED_AT),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_TABLE_ROUND_ID, models.ROUND_MEMBER_ROUND_ID),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_TABLE_USER_ID, models.ROUND_MEMBER_USER_ID),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_TABLE_SOLVED, models.ROUND_MEMBER_SOLVED),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_TABLE_FINISHED, models.ROUND_MEMBER_FINISHED),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_TABLE_SCORE, models.ROUND_MEMBER_SCORE),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_TABLE_FIRST_GUESS_AT, models.ROUND_MEMBER_FIRST_GUESS_AT),
	fmt.Sprintf("%s AS %s", models.ROUND_MEMBER_TABLE_COMPLETED_AT, models.ROUND_MEMBER_COMPLETED_AT),
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateRoundMember inserts a round member record
func (dao *DAO) CreateRoundMember(ctx context.Context, member *models.RoundMember) error {
	if member == nil {
		return utils.ErrNilPtr
	}

	if member.RoundID == "" {
		return utils.ErrId
	}

	if member.UserID == "" {
		return utils.ErrUserId
	}

	if member.ID == "" {
		member.RefreshId()
	}

	member.RefreshCreatedAt()
	member.RefreshUpdatedAt()

	builderOpts := newBuilderOptions(models.ROUND_MEMBER_TABLE).
		WithData(
			map[string]interface{}{
				models.BASE_ID:                     member.ID,
				models.ROUND_MEMBER_ROUND_ID:       member.RoundID,
				models.ROUND_MEMBER_USER_ID:        member.UserID,
				models.ROUND_MEMBER_SOLVED:         member.Solved,
				models.ROUND_MEMBER_FINISHED:       member.Finished,
				models.ROUND_MEMBER_SCORE:          member.Score,
				models.ROUND_MEMBER_FIRST_GUESS_AT: member.FirstGuessAt,
				models.ROUND_MEMBER_COMPLETED_AT:   member.CompletedAt,
				models.BASE_CREATED_AT:             member.CreatedAt,
				models.BASE_UPDATED_AT:             member.UpdatedAt,
			},
		)

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetRoundMember returns a round member record
func (dao *DAO) GetRoundMember(ctx context.Context, dbOpts *Options) (*models.RoundMember, error) {
	builderOpts := newBuilderOptions(models.ROUND_MEMBER_TABLE).
		WithColumns(roundMemberColumns...).
		SetDbOpts(dbOpts)

	return getGeneric[models.RoundMember](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRoundMembers returns round member records
func (dao *DAO) ListRoundMembers(ctx context.Context, dbOpts *Options) ([]*models.RoundMember, error) {
	builderOpts := newBuilderOptions(models.ROUND_MEMBER_TABLE).
		WithColumns(roundMemberColumns...).
		SetDbOpts(dbOpts)

	return listGeneric[models.RoundMember](ctx, dao, *builderOpts)
}
