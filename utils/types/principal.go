package types

import (
	"context"

	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Principal identifies the authenticated caller attached to a request context
type Principal struct {
	UserID   string
	SiteRole SiteRole
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// PrincipalFromContext returns the authenticated caller stored on ctx
func PrincipalFromContext(ctx context.Context) (Principal, error) {
	principal, ok := ctx.Value(PrincipalContextKey).(Principal)
	if !ok {
		return Principal{}, utils.ErrPrincipal
	}

	return principal, nil
}
