package service

import (
	"context"

	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// principalFromContext returns the authenticated principal from the context
func principalFromContext(ctx context.Context) (types.Principal, error) {
	principal, err := types.PrincipalFromContext(ctx)
	if err != nil {
		return types.Principal{}, err
	}

	return principal, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupMembershipFromContext returns the group membership from the context
func groupMembershipFromContext(ctx context.Context) (types.GroupMembership, error) {
	groupMembership, err := types.GroupMembershipFromContext(ctx)
	if err != nil {
		return types.GroupMembership{}, err
	}

	return groupMembership, nil
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
