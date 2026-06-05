package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var defaultGroupsListOrderBy = []string{models.GROUP_TABLE_NAME + " asc"}

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
				models.BASE_ID:            group.ID,
				models.GROUP_NAME:         group.Name,
				models.GROUP_CREATED_BY:   group.CreatedBy,
				models.BASE_CREATED_AT:    group.CreatedAt,
				models.BASE_UPDATED_AT:    group.UpdatedAt,
			},
		)

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroups returns group records
func (dao *DAO) ListGroups(ctx context.Context, dbOpts *Options) ([]*models.Group, error) {
	applyDefaultOrderBy(dbOpts, defaultGroupsListOrderBy)

	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.Group](ctx, dao, *builderOpts)
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
