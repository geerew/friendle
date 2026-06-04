package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var defaultGroupsListOrderBy = []string{models.GROUP_TABLE + "." + models.BASE_CREATED_AT + " desc"}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroup inserts a group record
func (dao *DAO) CreateGroup(ctx context.Context, g *models.Group) error {
	if g.ID == "" {
		g.RefreshId()
	}
	g.RefreshCreatedAt()
	g.RefreshUpdatedAt()
	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithData(map[string]interface{}{
			models.BASE_ID:          g.ID,
			models.GROUP_NAME:       g.Name,
			models.GROUP_CREATED_BY: g.CreatedBy,
			models.BASE_CREATED_AT:  g.CreatedAt,
			models.BASE_UPDATED_AT:  g.UpdatedAt,
		})

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroup returns a group record
func (dao *DAO) GetGroup(ctx context.Context, dbOpts *Options) (*models.Group, error) {
	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.Group](ctx, dao, *builderOpts)
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

// UpdateGroup updates a group record
func (dao *DAO) UpdateGroup(ctx context.Context, g *models.Group) error {
	if g.ID == "" {
		return utils.ErrId
	}

	g.RefreshUpdatedAt()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: g.ID})
	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithData(map[string]interface{}{
			models.GROUP_NAME:      g.Name,
			models.BASE_UPDATED_AT: g.UpdatedAt,
		}).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListAllGroups returns all group records
func (dao *DAO) ListAllGroups(ctx context.Context) ([]*models.Group, error) {
	return dao.ListGroups(ctx, NewOptions())
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteGroups deletes group records
//
// Errors when a where clause is not provided
func (dao *DAO) DeleteGroups(ctx context.Context, dbOpts *Options) error {
	if dbOpts == nil || dbOpts.Where == nil {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.GROUP_TABLE).SetDbOpts(dbOpts)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupNameSearchOrder ranks exact name matches first, then prefix matches, then other
// substring matches, then shorter names, then name ascending, then newest created
func groupNameSearchOrder(g, q string) squirrel.Sqlizer {
	return squirrel.Expr(
		`CASE WHEN LOWER(`+g+`.name) = ? THEN 0 WHEN LOWER(`+g+`.name) LIKE ? THEN 1 ELSE 2 END, LENGTH(`+g+`.name), LOWER(`+g+`.name), `+g+`.`+models.BASE_CREATED_AT+` DESC`,
		q, q+"%",
	)
}
