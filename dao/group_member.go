package dao

import (
	"context"
	"fmt"

	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupMemberColumns defines the columns to select
var groupMemberColumns = []string{
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_ID, models.BASE_ID),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_CREATED_AT, models.BASE_CREATED_AT),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_UPDATED_AT, models.BASE_UPDATED_AT),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_GROUP_ID, models.GROUP_MEMBER_GROUP_ID),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_USER_ID, models.GROUP_MEMBER_USER_ID),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_GROUP_ROLE, models.GROUP_MEMBER_GROUP_ROLE),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_TIMES_PICKED, models.GROUP_MEMBER_TIMES_PICKED),
	fmt.Sprintf("%s AS %s", models.GROUP_MEMBER_TABLE_PICKER_SKIPS, models.GROUP_MEMBER_PICKER_SKIPS),
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
				models.BASE_ID:                    member.ID,
				models.GROUP_MEMBER_GROUP_ID:      member.GroupID,
				models.GROUP_MEMBER_USER_ID:       member.UserID,
				models.GROUP_MEMBER_GROUP_ROLE:    member.GroupRole,
				models.GROUP_MEMBER_TIMES_PICKED:  member.TimesPicked,
				models.GROUP_MEMBER_PICKER_SKIPS:  member.PickerSkips,
				models.BASE_CREATED_AT:            member.CreatedAt,
				models.BASE_UPDATED_AT:            member.UpdatedAt,
			},
		)

	return createGeneric(ctx, dao, *builderOpts)
}
