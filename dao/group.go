package dao

import (
	"context"

	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

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
