package dao

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// roundColumns defines the columns to select for round queries
var roundColumns = []string{
	fmt.Sprintf("%s AS %s", models.ROUND_TABLE_ID, models.BASE_ID),
	fmt.Sprintf("%s AS %s", models.ROUND_TABLE_CREATED_AT, models.BASE_CREATED_AT),
	fmt.Sprintf("%s AS %s", models.ROUND_TABLE_UPDATED_AT, models.BASE_UPDATED_AT),
	fmt.Sprintf("%s AS %s", models.ROUND_TABLE_GROUP_ID, models.ROUND_GROUP_ID),
	fmt.Sprintf("%s AS %s", models.ROUND_TABLE_ROUND_DATE, models.ROUND_ROUND_DATE),
	fmt.Sprintf("%s AS %s", models.ROUND_TABLE_PICKER_USER_ID, models.ROUND_PICKER_USER_ID),
	fmt.Sprintf("%s AS %s", models.ROUND_TABLE_WORD_PLAIN, models.ROUND_WORD_PLAIN),
	fmt.Sprintf("%s AS %s", models.ROUND_TABLE_STATUS, models.ROUND_STATUS),
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateRound inserts a round record
func (dao *DAO) CreateRound(ctx context.Context, round *models.Round) error {
	if round == nil {
		return utils.ErrNilPtr
	}

	if round.GroupID == "" {
		return utils.ErrGroupId
	}

	if round.RoundDate == "" {
		return utils.ErrRoundDate
	}

	if round.PickerUserID == "" {
		return utils.ErrUserId
	}

	if !round.Status.IsValid() {
		round.Status = types.RoundAwaitingWord
	}

	if round.ID == "" {
		round.RefreshId()
	}

	round.RefreshCreatedAt()
	round.RefreshUpdatedAt()

	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithData(
			map[string]interface{}{
				models.BASE_ID:              round.ID,
				models.ROUND_GROUP_ID:       round.GroupID,
				models.ROUND_ROUND_DATE:     round.RoundDate,
				models.ROUND_PICKER_USER_ID: round.PickerUserID,
				models.ROUND_WORD_PLAIN:     round.WordPlain,
				models.ROUND_STATUS:         round.Status,
				models.BASE_CREATED_AT:      round.CreatedAt,
				models.BASE_UPDATED_AT:      round.UpdatedAt,
			},
		)

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetRound returns a round record
func (dao *DAO) GetRound(ctx context.Context, dbOpts *Options) (*models.Round, error) {
	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithColumns(roundColumns...).
		SetDbOpts(dbOpts)

	return getGeneric[models.Round](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRounds returns round records
func (dao *DAO) ListRounds(ctx context.Context, dbOpts *Options) ([]*models.Round, error) {
	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithColumns(roundColumns...).
		SetDbOpts(dbOpts)

	return listGeneric[models.Round](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateRound updates a round record
func (dao *DAO) UpdateRound(ctx context.Context, round *models.Round) error {
	if round == nil {
		return utils.ErrNilPtr
	}

	if round.ID == "" {
		return utils.ErrId
	}

	if !round.Status.IsValid() {
		return utils.ErrRoundStatus
	}

	round.RefreshUpdatedAt()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: round.ID})

	builderOpts := newBuilderOptions(models.ROUND_TABLE).
		WithData(
			map[string]interface{}{
				models.ROUND_WORD_PLAIN: round.WordPlain,
				models.ROUND_STATUS:     round.Status,
				models.BASE_UPDATED_AT:  round.UpdatedAt,
			},
		).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}
