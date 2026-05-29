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

// GroupListRow is a group projection for list queries with member count
type GroupListRow struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	MemberCount int    `db:"member_count"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupSearchRow is a group projection for name search results
type GroupSearchRow struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	MemberCount int    `db:"member_count"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroup inserts a group row
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

// GetGroup returns a group matching dbOpts
func (dao *DAO) GetGroup(ctx context.Context, dbOpts *Options) (*models.Group, error) {
	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.Group](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroups returns groups matching the given options
func (dao *DAO) ListGroups(ctx context.Context, dbOpts *Options) ([]*models.Group, error) {
	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.Group](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroupRows returns groups with member counts matching the given options
func (dao *DAO) ListGroupRows(ctx context.Context, dbOpts *Options) ([]*GroupListRow, error) {
	g := models.GROUP_TABLE
	gm := models.GROUP_MEMBER_TABLE

	applyDefaultOrderBy(dbOpts, defaultGroupsListOrderBy)

	builderOpts := newBuilderOptions(g).
		WithColumns(
			g+"."+models.BASE_ID+" AS id",
			g+".name AS name",
			"COUNT("+gm+".id) AS member_count",
		).
		WithLeftJoin(gm, gm+".group_id = "+g+"."+models.BASE_ID).
		WithGroupBy(g+"."+models.BASE_ID, g+".name").
		SetDbOpts(dbOpts)

	return listGeneric[GroupListRow](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SearchGroupSummaries searches groups by name with relevance ordering
func (dao *DAO) SearchGroupSummaries(ctx context.Context, q string, dbOpts *Options) ([]*GroupSearchRow, error) {
	g := models.GROUP_TABLE
	gm := models.GROUP_MEMBER_TABLE
	like := "%" + q + "%"

	if dbOpts == nil {
		dbOpts = NewOptions()
	}

	searchOpts := NewOptions().
		WithWhere(squirrel.Like{"LOWER(" + g + ".name)": like}).
		WithOrderByClause(searchGroupRelevanceOrder(g, q))

	if dbOpts.Pagination != nil {
		searchOpts = searchOpts.WithPagination(dbOpts.Pagination)
	}

	builderOpts := newBuilderOptions(g).
		WithColumns(
			g+"."+models.BASE_ID+" AS id",
			g+".name AS name",
			"(SELECT COUNT(*) FROM "+gm+" gm_count WHERE gm_count.group_id = "+g+"."+models.BASE_ID+") AS member_count",
		).
		SetDbOpts(searchOpts)

	return listGeneric[GroupSearchRow](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateGroup updates mutable group fields
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// searchGroupRelevanceOrder ranks exact name matches first, then prefix matches, then other
// substring matches, then shorter names, then name ascending, then newest created
func searchGroupRelevanceOrder(g, q string) squirrel.Sqlizer {
	return squirrel.Expr(
		`CASE WHEN LOWER(`+g+`.name) = ? THEN 0 WHEN LOWER(`+g+`.name) LIKE ? THEN 1 ELSE 2 END, LENGTH(`+g+`.name), LOWER(`+g+`.name), `+g+`.`+models.BASE_CREATED_AT+` DESC`,
		q, q+"%",
	)
}
