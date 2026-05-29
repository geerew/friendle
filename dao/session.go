package dao

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateOrReplaceSession inserts or replace a session record
func (dao *DAO) CreateOrReplaceSession(ctx context.Context, session *models.Session) error {
	if session == nil {
		return utils.ErrNilPtr
	}

	if session.ID == "" {
		return utils.ErrId
	}

	if session.UserId == "" {
		return utils.ErrUserId
	}

	builderOpts := newBuilderOptions(models.SESSION_TABLE).
		WithData(
			map[string]interface{}{
				models.BASE_ID:         session.ID,
				models.SESSION_USER_ID: session.UserId,
				models.SESSION_DATA:    session.Data,
				models.SESSION_EXPIRES: session.Expires,
			},
		).
		WithReplace()

	return createGeneric(ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GetSession gets a record from the sessions table based upon the where clause in the options. If
// there is no where clause, it will return the first record in the table
func (dao *DAO) GetSession(ctx context.Context, dbOpts *Options) (*models.Session, error) {
	builderOpts := newBuilderOptions(models.SESSION_TABLE).
		WithColumns(models.SessionColumns()...).
		SetDbOpts(dbOpts).
		WithLimit(1)

	return getGeneric[models.Session](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListSessions gets all records from the sessions table based upon the where clause and pagination
// in the options
func (dao *DAO) ListSessions(ctx context.Context, dbOpts *Options) ([]*models.Session, error) {
	builderOpts := newBuilderOptions(models.SESSION_TABLE).
		WithColumns(models.SessionColumns()...).
		SetDbOpts(dbOpts)

	return listGeneric[models.Session](ctx, dao, *builderOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateSession updates a session record
func (dao *DAO) UpdateSession(ctx context.Context, session *models.Session) error {
	if session == nil {
		return utils.ErrNilPtr
	}

	if session.ID == "" {
		return utils.ErrId
	}

	dbOpts := NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: session.ID})

	builderOpts := newBuilderOptions(models.SESSION_TABLE).
		WithData(
			map[string]interface{}{
				models.SESSION_DATA:    session.Data,
				models.SESSION_EXPIRES: session.Expires,
			},
		).
		SetDbOpts(dbOpts)

	_, err := updateGeneric(ctx, dao, *builderOpts)
	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// BulkUpdateSessions bulk updates session records
func (dao *DAO) BulkUpdateSessions(ctx context.Context, sessions []*models.Session) error {
	if sessions == nil {
		return utils.ErrNilPtr
	}

	if len(sessions) == 0 {
		return nil
	}

	builderOpts := newBuilderOptions(models.SESSION_TABLE).
		WithBulkUpdate(models.SESSION_DATA, models.SESSION_EXPIRES)

	for _, s := range sessions {
		if s == nil || s.ID == "" {
			return utils.ErrId
		}

		builderOpts = builderOpts.WithBulkUpdateRow(s.ID, s.Data, s.Expires)
	}

	_, err := updateGeneric(ctx, dao, *builderOpts)
	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteSessions deletes records from the sessions table
//
// Errors when a where clause is not provided
func (dao *DAO) DeleteSessions(ctx context.Context, dbOpts *Options) error {
	if dbOpts == nil || dbOpts.Where == nil {
		return utils.ErrWhere
	}

	builderOpts := newBuilderOptions(models.SESSION_TABLE).SetDbOpts(dbOpts)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)
	return err
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeleteAllSessions deletes all records from the sessions table
func (dao *DAO) DeleteAllSessions(ctx context.Context) error {
	builderOpts := newBuilderOptions(models.SESSION_TABLE)
	sqlStr, args, _ := deleteBuilder(*builderOpts)

	_, err := dao.db.ExecContext(ctx, sqlStr, args...)
	return err
}
