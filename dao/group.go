package dao

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var defaultGroupsListOrderBy = []string{models.GROUP_TABLE_NAME + " asc"}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupColumns defines the columns to select
var groupColumns = []string{
	fmt.Sprintf("%s AS %s", models.GROUP_TABLE_ID, models.BASE_ID),
	fmt.Sprintf("%s AS %s", models.GROUP_TABLE_CREATED_AT, models.BASE_CREATED_AT),
	fmt.Sprintf("%s AS %s", models.GROUP_TABLE_UPDATED_AT, models.BASE_UPDATED_AT),
	fmt.Sprintf("%s AS %s", models.GROUP_TABLE_NAME, models.GROUP_NAME),
	fmt.Sprintf("%s AS %s", models.GROUP_TABLE_CREATED_BY, models.GROUP_CREATED_BY),
	// Added via LEFT JOIN
	fmt.Sprintf("COUNT(%s) AS %s", models.GROUP_MEMBER_TABLE_ID, models.GROUP_MEMBER_COUNT),
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupJoins defines the joins to use during SELECT queries
var groupJoins = []join{
	{
		Type:      joinTypeLeft,
		Table:     models.GROUP_MEMBER_TABLE,
		Condition: models.GROUP_MEMBER_TABLE_GROUP_ID + " = " + models.GROUP_TABLE_ID,
	},
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupColumnsGroupBy defines the columns to group by
//
//	This is required because we are doing a LEFT JOIN to include the group member count
var groupColumnsGroupBy = []string{
	models.GROUP_TABLE_ID,
	models.GROUP_TABLE_CREATED_AT,
	models.GROUP_TABLE_UPDATED_AT,
	models.GROUP_TABLE_NAME,
	models.GROUP_TABLE_CREATED_BY,
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroup inserts a group record
func (dao *DAO) CreateGroup(ctx context.Context, group *models.Group) error {
	if group == nil {
		return utils.ErrNilPtr
	}

	if group.Name == "" {
		return utils.ErrGroupName
	}

	if group.CreatedBy == "" {
		return utils.ErrUserId
	}

	if group.ID == "" {
		group.RefreshId()
	}

	group.RefreshCreatedAt()
	group.RefreshUpdatedAt()

	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithData(
			map[string]interface{}{
				models.BASE_ID:          group.ID,
				models.GROUP_NAME:       group.Name,
				models.GROUP_CREATED_BY: group.CreatedBy,
				models.BASE_CREATED_AT:  group.CreatedAt,
				models.BASE_UPDATED_AT:  group.UpdatedAt,
			},
		)

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroups returns group records with member counts
func (dao *DAO) ListGroups(ctx context.Context, dbOpts *Options) ([]*models.Group, error) {
	applyDefaultOrderBy(dbOpts, defaultGroupsListOrderBy)

	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithColumns(groupColumns...).
		WithJoins(groupJoins...).
		WithGroupBy(groupColumnsGroupBy...).
		SetDbOpts(dbOpts)

	return listGeneric[models.Group](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroup returns a group record with a member count
func (dao *DAO) GetGroup(ctx context.Context, dbOpts *Options) (*models.Group, error) {
	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithColumns(groupColumns...).
		WithJoins(groupJoins...).
		WithGroupBy(groupColumnsGroupBy...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.Group](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MemberGroupsWhere builds a WHERE clause for groups the user belongs to
func MemberGroupsWhere(userID string) (squirrel.Sqlizer, error) {
	subSQL, subArgs, err := squirrel.Select(models.GROUP_MEMBER_GROUP_ID).
		From(models.GROUP_MEMBER_TABLE).
		Where(squirrel.Eq{models.GROUP_MEMBER_USER_ID: userID}).
		PlaceholderFormat(squirrel.Question).
		ToSql()

	if err != nil {
		return nil, err
	}

	return squirrel.Expr(models.GROUP_TABLE_ID+" IN ("+subSQL+")", subArgs...), nil
}
