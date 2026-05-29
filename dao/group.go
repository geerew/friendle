package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var defaultAdminGroupsListOrderBy = []string{models.GROUP_TABLE + "." + models.BASE_CREATED_AT + " desc"}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroup inserts a group row
func (dao *DAO) CreateGroup(ctx context.Context, g *models.Group) error {
	if g.ID == "" {
		g.RefreshId()
	}
	g.RefreshCreatedAt()
	g.RefreshUpdatedAt()
	return createGeneric(ctx, dao, *newBuilderOptions(models.GROUP_TABLE).WithData(map[string]interface{}{
		models.BASE_ID: g.ID, "name": g.Name, "created_by": g.CreatedBy,
		models.BASE_CREATED_AT: g.CreatedAt, models.BASE_UPDATED_AT: g.UpdatedAt,
	}))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroup returns a group by ID
func (dao *DAO) GetGroup(ctx context.Context, id string) (*models.Group, error) {
	return getGeneric[models.Group](ctx, dao, *newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: id})).
		WithLimit(1))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroups returns groups matching the given options
func (dao *DAO) ListGroups(ctx context.Context, dbOpts *Options) ([]*models.Group, error) {
	return listGeneric[models.Group](ctx, dao, *newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).SetDbOpts(dbOpts))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SearchGroups searches groups by name
func (dao *DAO) SearchGroups(ctx context.Context, q string) ([]*models.Group, error) {
	like := "%" + q + "%"
	return listGeneric[models.Group](ctx, dao, *newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Like{"LOWER(name)": like})))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateGroup updates mutable group fields
func (dao *DAO) UpdateGroup(ctx context.Context, g *models.Group) error {
	g.RefreshUpdatedAt()
	_, err := updateGeneric(ctx, dao, *newBuilderOptions(models.GROUP_TABLE).WithData(map[string]interface{}{
		"name": g.Name,
		models.BASE_UPDATED_AT: g.UpdatedAt,
	}).SetDbOpts(NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: g.ID})))

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListAllGroups returns every group
func (dao *DAO) ListAllGroups(ctx context.Context) ([]*models.Group, error) {
	return dao.ListGroups(ctx, NewOptions())
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteGroups deletes records from the groups table
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
