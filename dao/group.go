package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var defaultGroupsListOrderBy = []string{models.GROUP_TABLE + "." + models.BASE_CREATED_AT + " desc"}

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
//
// Members, join requests, member count, and users are not included by default. Enable them
// with WithMembers(), WithJoinRequests(), WithMemberCount(), and WithUsers() on the options
func (dao *DAO) GetGroup(ctx context.Context, dbOpts *Options) (*models.Group, error) {
	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	if !groupRelationsRequested(dbOpts) {
		return getGeneric[models.Group](ctx, dao, *builderOpts)
	}

	group, err := getGeneric[models.Group](ctx, dao, *builderOpts)
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, nil
	}

	if err := attachGroupRelations(ctx, dao, []*models.Group{group}, dbOpts); err != nil {
		return nil, err
	}

	return group, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroups returns groups matching the given options
//
// Members, join requests, member count, and users are not included by default. Enable them
// with WithMembers(), WithJoinRequests(), WithMemberCount(), and WithUsers() on the options
func (dao *DAO) ListGroups(ctx context.Context, dbOpts *Options) ([]*models.Group, error) {
	applyDefaultOrderBy(dbOpts, defaultGroupsListOrderBy)

	builderOpts := newBuilderOptions(models.GROUP_TABLE).
		WithColumns(models.GroupColumns()...).
		SetDbOpts(dbOpts)

	if !groupRelationsRequested(dbOpts) {
		return listGeneric[models.Group](ctx, dao, *builderOpts)
	}

	groups, err := listGeneric[models.Group](ctx, dao, *builderOpts)
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return groups, nil
	}

	if err := attachGroupRelations(ctx, dao, groups, dbOpts); err != nil {
		return nil, err
	}

	return groups, nil
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

// groupNameSearchOrder ranks exact name matches first, then prefix matches, then other
// substring matches, then shorter names, then name ascending, then newest created
func groupNameSearchOrder(g, q string) squirrel.Sqlizer {
	return squirrel.Expr(
		`CASE WHEN LOWER(`+g+`.name) = ? THEN 0 WHEN LOWER(`+g+`.name) LIKE ? THEN 1 ELSE 2 END, LENGTH(`+g+`.name), LOWER(`+g+`.name), `+g+`.`+models.BASE_CREATED_AT+` DESC`,
		q, q+"%",
	)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupRelationsRequested reports whether dbOpts requests group relation data
func groupRelationsRequested(dbOpts *Options) bool {
	if dbOpts == nil {
		return false
	}

	return dbOpts.IncludeMembers || dbOpts.IncludeJoinRequests || dbOpts.IncludeMemberCount
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// attachGroupRelations attaches members, join requests, and member counts to groups
func attachGroupRelations(ctx context.Context, dao *DAO, groups []*models.Group, dbOpts *Options) error {
	if len(groups) == 0 || dbOpts == nil {
		return nil
	}

	groupIDs := utils.Map(groups, func(g *models.Group) string {
		return g.ID
	})

	membersByGroup := make(map[string][]*models.GroupMember)
	if dbOpts.IncludeMembers {
		members, err := dao.ListGroupMembers(ctx, NewOptions().WithWhere(squirrel.Eq{
			models.GROUP_MEMBER_GROUP_ID: groupIDs,
		}))
		if err != nil {
			return err
		}

		for _, m := range members {
			membersByGroup[m.GroupID] = append(membersByGroup[m.GroupID], m)
		}
	}

	requestsByGroup := make(map[string][]*models.GroupJoinRequest)
	if dbOpts.IncludeJoinRequests {
		requests, err := dao.ListJoinRequests(ctx, NewOptions().WithWhere(squirrel.Eq{
			models.JOIN_REQUEST_GROUP_ID: groupIDs,
			models.JOIN_REQUEST_STATUS:   types.JoinPending,
		}))
		if err != nil {
			return err
		}

		for _, jr := range requests {
			requestsByGroup[jr.GroupID] = append(requestsByGroup[jr.GroupID], jr)
		}
	}

	memberCounts := map[string]int{}
	if dbOpts.IncludeMemberCount {
		counts, err := memberCountsForGroups(ctx, dao, groupIDs)
		if err != nil {
			return err
		}

		memberCounts = counts
	}

	if dbOpts.IncludeUsers && dbOpts.IncludeMembers {
		userIDs := make([]string, 0)
		for _, members := range membersByGroup {
			for _, m := range members {
				userIDs = append(userIDs, m.UserID)
			}
		}

		userMap, err := usersByIDs(ctx, dao, userIDs)
		if err != nil {
			return err
		}

		for _, members := range membersByGroup {
			for _, m := range members {
				m.User = userMap[m.UserID]
			}
		}
	}

	for _, g := range groups {
		if dbOpts.IncludeMembers {
			g.Members = membersByGroup[g.ID]
			if g.Members == nil {
				g.Members = []*models.GroupMember{}
			}
		}

		if dbOpts.IncludeJoinRequests {
			g.JoinRequests = requestsByGroup[g.ID]
			if g.JoinRequests == nil {
				g.JoinRequests = []*models.GroupJoinRequest{}
			}
		}

		if dbOpts.IncludeMemberCount {
			g.MemberCount = memberCounts[g.ID]
		}
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// memberCountsForGroups returns member counts keyed by group ID
func memberCountsForGroups(ctx context.Context, dao *DAO, groupIDs []string) (map[string]int, error) {
	if len(groupIDs) == 0 {
		return map[string]int{}, nil
	}

	members, err := dao.ListGroupMembers(ctx, NewOptions().WithWhere(squirrel.Eq{
		models.GROUP_MEMBER_GROUP_ID: groupIDs,
	}))
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int, len(groupIDs))
	for _, m := range members {
		counts[m.GroupID]++
	}

	return counts, nil
}
