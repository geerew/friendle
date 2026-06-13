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

// groupMemberColumns defines the columns to select for group member queries
var groupMemberColumns = []string{
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_ID, models.BASE_ID),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_CREATED_AT, models.BASE_CREATED_AT),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_UPDATED_AT, models.BASE_UPDATED_AT),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_GROUP_ID, models.GROUP_MEMBER_GROUP_ID),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_USER_ID, models.GROUP_MEMBER_USER_ID),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_GROUP_ROLE, models.GROUP_MEMBER_GROUP_ROLE),
	fmt.Sprintf("%s AS %s", models.USER_TABLE_DISPLAY_NAME, models.USER_DISPLAY_NAME),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_TIMES_PICKED, models.GROUP_MEMBER_TIMES_PICKED),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_PICKER_SKIPS, models.GROUP_MEMBER_PICKER_SKIPS),
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var groupMemberJoins = []join{
	{
		Type:      joinTypeInner,
		Table:     models.USER_TABLE,
		Condition: models.GROUP_MEMBER_TABLE_USER_ID + " = " + models.USER_TABLE_ID,
	},
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CountGroupMembers returns the number of group membership records
func (dao *DAO) CountGroupMembers(ctx context.Context, dbOpts *Options) (int, error) {
	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).SetDbOpts(dbOpts)

	return countGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroupMember inserts a group membership record
func (dao *DAO) CreateGroupMember(ctx context.Context, member *models.GroupMember) error {
	if member == nil {
		return utils.ErrNilPtr
	}

	if member.GroupID == "" {
		return utils.ErrGroupId
	}

	if member.UserID == "" {
		return utils.ErrUserId
	}

	if !member.GroupRole.IsValid() {
		member.GroupRole = types.GroupRoleUser
	}

	if member.ID == "" {
		member.RefreshId()
	}

	member.RefreshCreatedAt()
	member.RefreshUpdatedAt()

	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithData(
			map[string]interface{}{
				models.BASE_ID:                   member.ID,
				models.GROUP_MEMBER_GROUP_ID:     member.GroupID,
				models.GROUP_MEMBER_USER_ID:      member.UserID,
				models.GROUP_MEMBER_GROUP_ROLE:   member.GroupRole,
				models.GROUP_MEMBER_TIMES_PICKED: member.TimesPicked,
				models.GROUP_MEMBER_PICKER_SKIPS: member.PickerSkips,
				models.BASE_CREATED_AT:           member.CreatedAt,
				models.BASE_UPDATED_AT:           member.UpdatedAt,
			},
		)

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetGroupMember returns a group membership record with the user display name
func (dao *DAO) GetGroupMember(ctx context.Context, dbOpts *Options) (*models.GroupMember, error) {
	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns(groupMemberColumns...).
		WithJoins(groupMemberJoins...).
		SetDbOpts(dbOpts)

	return getGeneric[models.GroupMember](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroupMembers returns group membership records with user display names
func (dao *DAO) ListGroupMembers(ctx context.Context, dbOpts *Options) ([]*models.GroupMember, error) {
	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithColumns(groupMemberColumns...).
		WithJoins(groupMemberJoins...).
		SetDbOpts(dbOpts)

	return listGeneric[models.GroupMember](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateGroupMember updates mutable fields on a group membership record
func (dao *DAO) UpdateGroupMember(ctx context.Context, member *models.GroupMember) error {
	if member == nil {
		return utils.ErrNilPtr
	}

	if member.ID == "" {
		return utils.ErrId
	}

	member.RefreshUpdatedAt()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: member.ID})

	builderOpts := newBuilderOptions(models.GROUP_MEMBER_TABLE).
		WithData(
			map[string]interface{}{
				models.GROUP_MEMBER_GROUP_ROLE:   member.GroupRole,
				models.GROUP_MEMBER_TIMES_PICKED: member.TimesPicked,
				models.GROUP_MEMBER_PICKER_SKIPS: member.PickerSkips,
				models.BASE_UPDATED_AT:           member.UpdatedAt,
			},
		).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)

	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteGroupMembers deletes group membership records
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
