package dao

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userColumns defines the columns to select
var userColumns = []string{
	fmt.Sprintf("%s AS %s", models.USER_TABLE_ID, models.BASE_ID),
	fmt.Sprintf("%s AS %s", models.USER_TABLE_CREATED_AT, models.BASE_CREATED_AT),
	fmt.Sprintf("%s AS %s", models.USER_TABLE_UPDATED_AT, models.BASE_UPDATED_AT),
	fmt.Sprintf("%s AS %s", models.USER_TABLE_USERNAME, models.USER_USERNAME),
	fmt.Sprintf("%s AS %s", models.USER_TABLE_DISPLAY_NAME, models.USER_DISPLAY_NAME),
	fmt.Sprintf("%s AS %s", models.USER_TABLE_PASSWORD_HASH, models.USER_PASSWORD_HASH),
	fmt.Sprintf("%s AS %s", models.USER_TABLE_SITE_ROLE, models.USER_SITE_ROLE),
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateUser inserts a user record
func (dao *DAO) CreateUser(ctx context.Context, user *models.User) error {
	if user == nil {
		return utils.ErrNilPtr
	}

	if user.Username == "" {
		return utils.ErrUsername
	}

	if user.PasswordHash == "" {
		return utils.ErrUserPassword
	}

	if user.ID == "" {
		user.RefreshId()
	}

	user.RefreshCreatedAt()
	user.RefreshUpdatedAt()

	builderOpts := newBuilderOptions(models.USER_TABLE).
		WithData(
			map[string]interface{}{
				models.BASE_ID:            user.ID,
				models.USER_USERNAME:      user.Username,
				models.USER_DISPLAY_NAME:  user.DisplayName,
				models.USER_PASSWORD_HASH: user.PasswordHash,
				models.USER_SITE_ROLE:     user.SiteRole,
				models.BASE_CREATED_AT:    user.CreatedAt,
				models.BASE_UPDATED_AT:    user.UpdatedAt,
			},
		)

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CountUsers returns the number of user records
func (dao *DAO) CountUsers(ctx context.Context, dbOpts *Options) (int, error) {
	builderOpts := newBuilderOptions(models.USER_TABLE).SetDbOpts(dbOpts)
	return countGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetUser returns a user record
func (dao *DAO) GetUser(ctx context.Context, dbOpts *Options) (*models.User, error) {
	builderOpts := newBuilderOptions(models.USER_TABLE).
		WithColumns(userColumns...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.User](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListUsers returns user records
func (dao *DAO) ListUsers(ctx context.Context, dbOpts *Options) ([]*models.User, error) {
	builderOpts := newBuilderOptions(models.USER_TABLE).
		WithColumns(userColumns...).
		SetDbOpts(dbOpts)

	return listGeneric[models.User](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateUser updates a user record
func (dao *DAO) UpdateUser(ctx context.Context, user *models.User) error {
	if user == nil {
		return utils.ErrNilPtr
	}

	if user.ID == "" {
		return utils.ErrId
	}

	if user.PasswordHash == "" {
		return utils.ErrUserPassword
	}

	user.RefreshUpdatedAt()

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: user.ID})

	builderOpts := newBuilderOptions(models.USER_TABLE).
		WithData(
			map[string]interface{}{
				models.USER_DISPLAY_NAME:  user.DisplayName,
				models.USER_PASSWORD_HASH: user.PasswordHash,
				models.USER_SITE_ROLE:     user.SiteRole,
				models.BASE_UPDATED_AT:    user.UpdatedAt,
			},
		).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)
	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteUsers deletes user records
//
// Errors when a where clause is not provided
func (dao *DAO) DeleteUsers(ctx context.Context, dbOpts *Options) error {
	if dbOpts == nil || dbOpts.Where == nil {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.USER_TABLE).SetDbOpts(dbOpts)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)
	return err
}
