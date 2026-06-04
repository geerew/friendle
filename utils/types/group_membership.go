package types

import (
	"context"

	"github.com/geerew/friendle/utils"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupMembership is the caller's membership for the group on the current request, set by API
// middleware after it verifies the user belongs to that group
type GroupMembership struct {
	GroupRole GroupRole
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const GroupMembershipContextKey contextKey = "groupMembership"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupMembershipFromContext returns membership set by API auth middleware for this request
func GroupMembershipFromContext(ctx context.Context) (GroupMembership, error) {
	groupMembership, ok := ctx.Value(GroupMembershipContextKey).(GroupMembership)
	if !ok {
		return GroupMembership{}, utils.ErrGroupMembership
	}

	return groupMembership, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithGroupMembership returns a context carrying group membership for the caller
func WithGroupMembership(ctx context.Context, groupRole GroupRole) context.Context {
	return context.WithValue(ctx, GroupMembershipContextKey, GroupMembership{
		GroupRole: groupRole,
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// isAdmin returns whether the caller is a group admin
func (m GroupMembership) IsAdmin() bool {
	return m.GroupRole == GroupRoleAdmin
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// isMember returns whether the caller is a group member
func (m GroupMembership) IsMember() bool {
	return m.GroupRole == GroupRoleUser
}
