package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AdminUserRow is a user projection for the site admin list query
type AdminUserRow struct {
	ID          string         `db:"id"`
	Username    string         `db:"username"`
	DisplayName string         `db:"display_name"`
	SiteRole    types.SiteRole `db:"site_role"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AdminGroupRow is a group projection for the site admin list query
type AdminGroupRow struct {
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

// UserGroupSummaryRow is a group projection for a user's memberships
type UserGroupSummaryRow struct {
	UserID      string          `db:"user_id"`
	ID          string          `db:"id"`
	Name        string          `db:"name"`
	MemberCount int             `db:"member_count"`
	GroupRole   types.GroupRole `db:"group_role"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListAdminGroups returns groups with member counts for the site admin list
func (dao *DAO) ListAdminGroups(ctx context.Context, dbOpts *Options) ([]*AdminGroupRow, error) {
	g := models.GROUP_TABLE
	gm := models.GROUP_MEMBER_TABLE

	applyDefaultOrderBy(dbOpts, defaultAdminGroupsListOrderBy)

	return listGeneric[AdminGroupRow](ctx, dao, *newBuilderOptions(g).
		WithColumns(
			g+"."+models.BASE_ID+" AS id",
			g+".name AS name",
			"COUNT("+gm+".id) AS member_count",
		).
		WithLeftJoin(gm, gm+".group_id = "+g+"."+models.BASE_ID).
		WithGroupBy(g+"."+models.BASE_ID, g+".name").
		SetDbOpts(dbOpts))
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

	return listGeneric[GroupSearchRow](ctx, dao, *newBuilderOptions(g).
		WithColumns(
			g+"."+models.BASE_ID+" AS id",
			g+".name AS name",
			"(SELECT COUNT(*) FROM "+gm+" gm_count WHERE gm_count.group_id = "+g+"."+models.BASE_ID+") AS member_count",
		).
		SetDbOpts(searchOpts))
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListUserGroupSummariesForUserIDs returns group summaries for the given user IDs, ordered by
// user then group name
func (dao *DAO) ListUserGroupSummariesForUserIDs(ctx context.Context, userIDs []string) ([]*UserGroupSummaryRow, error) {
	if len(userIDs) == 0 {
		return []*UserGroupSummaryRow{}, nil
	}

	gm := models.GROUP_MEMBER_TABLE
	g := models.GROUP_TABLE

	return listGeneric[UserGroupSummaryRow](ctx, dao, *newBuilderOptions(gm).
		WithColumns(
			gm+".user_id AS user_id",
			g+"."+models.BASE_ID+" AS id",
			g+".name AS name",
			gm+".group_role AS group_role",
			"(SELECT COUNT(*) FROM "+gm+" gm_count WHERE gm_count.group_id = "+g+"."+models.BASE_ID+") AS member_count",
		).
		WithJoin(g, g+"."+models.BASE_ID+" = "+gm+".group_id").
		SetDbOpts(NewOptions().
			WithWhere(squirrel.Eq{gm + ".user_id": userIDs}).
			WithOrderBy(gm+".user_id ASC", g+"."+models.BASE_CREATED_AT+" desc")))
}
