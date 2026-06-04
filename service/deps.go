package service

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userContext holds the caller identity and group-scoped access for read paths
type userContext struct {
	UserID     string
	SiteAdmin  bool
	GroupAdmin bool
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userContext loads group-scoped access for the caller on ctx
func (d deps) userContext(ctx context.Context, groupID string) (userContext, error) {
	principal, err := types.PrincipalFromContext(ctx)
	if err != nil {
		return userContext{}, err
	}

	viewer := userContext{
		UserID:    principal.UserID,
		SiteAdmin: principal.SiteRole == types.SiteRoleAdmin,
	}

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{
		models.GROUP_MEMBER_GROUP_ID: groupID,
		models.GROUP_MEMBER_USER_ID:  principal.UserID,
	})
	if m, _ := d.dao.GetGroupMember(ctx, dbOpts); m != nil {
		viewer.GroupAdmin = m.GroupRole == types.GroupRoleAdmin
	}

	return viewer, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// validatePassword checks password length constraints for admin user writes
func validatePassword(password string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	if len(password) > 128 {
		return ErrPasswordTooLong
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// uniqueStrings returns deduplicated non-empty strings preserving first-seen order
func uniqueStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	unique := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}

		if _, ok := seen[value]; ok {
			continue
		}

		seen[value] = struct{}{}
		unique = append(unique, value)
	}

	return unique
}
