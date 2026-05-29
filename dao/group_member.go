package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroupMember inserts a group membership row
func (dao *DAO) CreateGroupMember(ctx context.Context, m *models.GroupMember) error {
	if m.ID == "" {
		m.RefreshId()
	}
	m.RefreshCreatedAt()
	m.RefreshUpdatedAt()
	return createGeneric(ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).WithData(map[string]interface{}{
		models.BASE_ID: m.ID, "group_id": m.GroupID, "user_id": m.UserID, "group_role": m.GroupRole,
		"times_picked": m.TimesPicked, "picker_skips": m.PickerSkips,
		models.BASE_CREATED_AT: m.CreatedAt, models.BASE_UPDATED_AT: m.UpdatedAt,
	}))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroupMember returns a membership for a user in a group
func (dao *DAO) GetGroupMember(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	return getGeneric[models.GroupMember](ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns(models.GroupMemberColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"group_id": groupID, "user_id": userID})).
		WithLimit(1))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroupMembers returns all members of a group
func (dao *DAO) ListGroupMembers(ctx context.Context, groupID string) ([]*models.GroupMember, error) {
	return listGeneric[models.GroupMember](ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns(models.GroupMemberColumns()...).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"group_id": groupID})))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteGroupMember removes a user from a group
func (dao *DAO) DeleteGroupMember(ctx context.Context, groupID, userID string) error {
	if groupID == "" || userID == "" {
		return utils.ErrWhere
	}
	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).SetDbOpts(
		NewOptions().WithWhere(squirrel.Eq{"group_id": groupID, "user_id": userID}))
	sqlStr, args, _ := deleteBuilder(*builderOpts)
	_, err := dao.db.ExecContext(ctx, sqlStr, args...)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IncrementTimesPicked bumps the picker count for a member
func (dao *DAO) IncrementTimesPicked(ctx context.Context, memberID string) error {
	now := types.NowDateTime().String()
	sql := `UPDATE group_members SET times_picked = times_picked + 1, updated_at = ? WHERE id = ?`
	_, err := dao.db.ExecContext(ctx, sql, now, memberID)

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

// CountGroupMembers returns the number of members in a group
func (dao *DAO) CountGroupMembers(ctx context.Context, groupID string) (int, error) {
	return countGeneric(ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{"group_id": groupID})))
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

	rows, err := listGeneric[groupIDRow](ctx, dao, *newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns("group_id").
		SetDbOpts(NewOptions().WithWhere(squirrel.Eq{
			"user_id":  userID,
			"group_id": groupIDs,
		})))
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.GroupID)
	}

	return ids, nil
}
