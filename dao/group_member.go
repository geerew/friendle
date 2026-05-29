package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroupMember inserts a group membership row
func (dao *DAO) CreateGroupMember(ctx context.Context, m *models.GroupMember) error {
	if m.ID == "" {
		m.RefreshId()
	}
	m.RefreshCreatedAt()
	m.RefreshUpdatedAt()

	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithData(map[string]interface{}{
			models.BASE_ID:                   m.ID,
			models.GROUP_MEMBER_GROUP_ID:     m.GroupID,
			models.GROUP_MEMBER_USER_ID:      m.UserID,
			models.GROUP_MEMBER_GROUP_ROLE:   m.GroupRole,
			models.GROUP_MEMBER_TIMES_PICKED: m.TimesPicked,
			models.GROUP_MEMBER_PICKER_SKIPS: m.PickerSkips,
			models.BASE_CREATED_AT:           m.CreatedAt,
			models.BASE_UPDATED_AT:           m.UpdatedAt,
		})

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroupMember returns a group member matching dbOpts
func (dao *DAO) GetGroupMember(ctx context.Context, dbOpts *Options) (*models.GroupMember, error) {
	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns(models.GroupMemberColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.GroupMember](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroupMembers returns group members matching dbOpts
func (dao *DAO) ListGroupMembers(ctx context.Context, dbOpts *Options) ([]*models.GroupMember, error) {
	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns(models.GroupMemberColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.GroupMember](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteGroupMembers deletes records from the group_members table
//
// Errors when a where clause is not provided
func (dao *DAO) DeleteGroupMembers(ctx context.Context, dbOpts *Options) error {
	if dbOpts == nil || dbOpts.Where == nil {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).SetDbOpts(dbOpts)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)
	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateGroupMember updates mutable group member fields
func (dao *DAO) UpdateGroupMember(ctx context.Context, m *models.GroupMember) error {
	if m.ID == "" {
		return utils.ErrId
	}

	m.RefreshUpdatedAt()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: m.ID})
	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithData(map[string]interface{}{
			models.GROUP_MEMBER_GROUP_ROLE:   m.GroupRole,
			models.GROUP_MEMBER_TIMES_PICKED: m.TimesPicked,
			models.GROUP_MEMBER_PICKER_SKIPS: m.PickerSkips,
			models.BASE_UPDATED_AT:           m.UpdatedAt,
		}).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// MembersWithMinPicks returns members tied for fewest picks in a group
func (dao *DAO) MembersWithMinPicks(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	query := `SELECT id, group_id, user_id, group_role, times_picked, picker_skips, created_at, updated_at
FROM group_members WHERE group_id = ? AND times_picked = (
  SELECT MIN(times_picked) FROM group_members WHERE group_id = ?)`
	rows, err := dao.db.QueryContext(ctx, query, groupID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.GroupMember
	for rows.Next() {
		m := &models.GroupMember{}
		if err := rows.Scan(&m.ID, &m.GroupID, &m.UserID, &m.GroupRole, &m.TimesPicked, &m.PickerSkips, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}

	return members, rows.Err()
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CountGroupMembers returns the number of members matching dbOpts
func (dao *DAO) CountGroupMembers(ctx context.Context, dbOpts *Options) (int, error) {
	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).SetDbOpts(dbOpts)

	return countGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListMemberGroupIDsForUser returns group IDs the user belongs to from the given set
func (dao *DAO) ListMemberGroupIDsForUser(ctx context.Context, userID string, groupIDs []string) ([]string, error) {
	if len(groupIDs) == 0 {
		return []string{}, nil
	}

	type groupIDRow struct {
		GroupID string `db:"group_id"`
	}

	dbOpts := NewOptions().WithWhere(squirrel.Eq{
		models.GROUP_MEMBER_USER_ID:  userID,
		models.GROUP_MEMBER_GROUP_ID: groupIDs,
	})
	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns(models.GROUP_MEMBER_GROUP_ID).
		SetDbOpts(dbOpts)

	rows, err := listGeneric[groupIDRow](ctx, dao, *builderOpts)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.GroupID)
	}

	return ids, nil
}
