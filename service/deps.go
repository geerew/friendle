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
