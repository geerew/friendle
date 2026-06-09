package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/pagination"
	"github.com/geerew/friendle/utils/types"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var defaultGroupsListOrderBy = []string{models.GROUP_TABLE_NAME + " asc"}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

const maxGroupNameLength = 64

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroupRequest represents a group create request
type CreateGroupRequest struct {
	Name string `json:"name"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupResponse represents a group response
type GroupResponse struct {
	ID                string                   `json:"id"`
	CreatedAt         types.DateTime           `json:"createdAt"`
	UpdatedAt         types.DateTime           `json:"updatedAt"`
	Name              string                   `json:"name"`
	CreatedBy         string                   `json:"createdBy,omitempty"`
	MemberCount       int                      `json:"memberCount"`
	GroupRole         *types.GroupRole         `json:"groupRole,omitempty"`
	JoinRequestStatus *types.JoinRequestStatus `json:"joinRequestStatus,omitempty"`
	AdminSummary      *GroupAdminSummary       `json:"adminSummary,omitempty"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupAdminSummary holds join request counts for a group admin
type GroupAdminSummary struct {
	PendingJoinRequestCount  int `json:"pendingJoinRequestCount"`
	RejectedJoinRequestCount int `json:"rejectedJoinRequestCount"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type groupRolesByGroupID map[string]types.GroupRole

type joinRequestStatusesByGroupID map[string]types.JoinRequestStatus

// userMemberStatus holds the caller's group role and join request status keyed by group ID
type userMemberStatus struct {
	roles    groupRolesByGroupID
	requests joinRequestStatusesByGroupID
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Groups orchestrates friend group operations
type Groups struct {
	deps
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// newGroups creates a Groups service
func newGroups(d deps) *Groups {
	return &Groups{deps: d}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Create creates a group and sets the caller as group admin
func (g *Groups) Create(ctx context.Context, req CreateGroupRequest) (*GroupResponse, error) {
	principal, err := principalFromContext(ctx)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrGroupNameRequired
	}

	if utf8.RuneCountInString(name) > maxGroupNameLength {
		return nil, ErrGroupNameTooLong
	}

	group := &models.Group{
		Name:      name,
		CreatedBy: principal.UserID,
	}

	err = g.dao.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := g.dao.CreateGroup(txCtx, group); err != nil {
			if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
				return ErrGroupNameTaken
			}

			return err
		}

		member := &models.GroupMember{
			GroupID:   group.ID,
			UserID:    principal.UserID,
			GroupRole: types.GroupRoleAdmin,
		}

		return g.dao.CreateGroupMember(txCtx, member)
	})
	if err != nil {
		return nil, err
	}

	group.MemberCount = 1

	return groupResponseBuilder(group, nil, nil), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// List returns a paginated slice of groups
func (g *Groups) List(ctx context.Context, page *pagination.Pagination) ([]*GroupResponse, error) {
	if _, err := principalFromContext(ctx); err != nil {
		return nil, err
	}

	groups, err := g.dao.ListGroups(ctx, dao.NewOptions().
		WithPagination(page).
		WithOrderBy(defaultGroupsListOrderBy...))
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return []*GroupResponse{}, nil
	}

	return groupsResponsesBuilder(groups, userMemberStatus{}), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListSelf returns a paginated slice of groups the authenticated user belongs to
func (g *Groups) ListSelf(ctx context.Context, page *pagination.Pagination) ([]*GroupResponse, error) {
	principal, err := principalFromContext(ctx)
	if err != nil {
		return nil, err
	}

	where, err := dao.MemberGroupsWhere(principal.UserID)
	if err != nil {
		return nil, err
	}

	groups, err := g.dao.ListGroups(ctx, dao.NewOptions().
		WithPagination(page).
		WithOrderBy(defaultGroupsListOrderBy...).
		WithWhere(where))
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return []*GroupResponse{}, nil
	}

	groupIDs := utils.Map(groups, func(group *models.Group) string { return group.ID })
	status, err := g.getUserMemberStatus(ctx, principal.UserID, groupIDs)
	if err != nil {
		return nil, err
	}

	return groupsResponsesBuilder(groups, status), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Search returns a paginated slice of groups whose names contain the search term
//
// # The result is ordered by prefix matches first, then substring matches
//
// Results include the caller's group role and join request status for each group
func (g *Groups) Search(ctx context.Context, page *pagination.Pagination, name string) ([]*GroupResponse, error) {
	principal, err := principalFromContext(ctx)
	if err != nil {
		return nil, err
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrGroupSearchQueryRequired
	}

	lowerName := strings.ToLower(name)
	whereClause := squirrel.Like{"LOWER(" + models.GROUP_TABLE_NAME + ")": "%" + lowerName + "%"}
	orderByClause := squirrel.Expr(
		"CASE WHEN LOWER("+models.GROUP_TABLE_NAME+") LIKE ? THEN 0 ELSE 1 END, LOWER("+models.GROUP_TABLE_NAME+") ASC",
		lowerName+"%",
	)

	groups, err := g.dao.ListGroups(ctx, dao.NewOptions().
		WithPagination(page).
		WithOrderByClause(orderByClause).
		WithWhere(whereClause))
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return []*GroupResponse{}, nil
	}

	groupIDs := utils.Map(groups, func(group *models.Group) string { return group.ID })
	status, err := g.getUserMemberStatus(ctx, principal.UserID, groupIDs)
	if err != nil {
		return nil, err
	}

	return groupsResponsesBuilder(groups, status), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Get returns a group by ID
func (g *Groups) Get(ctx context.Context, groupID string) (*GroupResponse, error) {
	principal, err := principalFromContext(ctx)
	if err != nil {
		return nil, err
	}

	group, err := g.dao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, ErrGroupNotFound
	}

	status, err := g.getUserMemberStatus(ctx, principal.UserID, []string{groupID})
	if err != nil {
		return nil, err
	}

	groupRole, joinRequestStatus := status.forGroup(groupID)

	resp := groupResponseBuilder(group, groupRole, joinRequestStatus)
	if groupRole != nil && *groupRole == types.GroupRoleAdmin {
		adminSummary, err := g.getAdminSummary(ctx, groupID)
		if err != nil {
			return nil, err
		}

		resp.AdminSummary = adminSummary
	}

	return resp, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RequestJoin creates a pending join request for the authenticated user
func (g *Groups) RequestJoin(ctx context.Context, groupID string) (*GroupResponse, error) {
	principal, err := principalFromContext(ctx)
	if err != nil {
		return nil, err
	}

	group, err := g.dao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, ErrGroupNotFound
	}

	status, err := g.getUserMemberStatus(ctx, principal.UserID, []string{groupID})
	if err != nil {
		return nil, err
	}

	groupRole, joinRequestStatus := status.forGroup(groupID)
	if groupRole != nil && groupRole.IsValid() {
		return nil, ErrGroupAlreadyMember
	}

	if joinRequestStatus != nil {
		switch *joinRequestStatus {
		case types.JoinPending:
			return nil, ErrGroupJoinRequestPending
		case types.JoinRejected:
			return nil, ErrGroupJoinRequestRejected
		}
	}

	err = g.dao.CreateGroupJoinRequest(ctx, &models.GroupJoinRequest{
		GroupID: groupID,
		UserID:  principal.UserID,
		Status:  types.JoinPending,
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			daoOpts := dao.NewOptions().WithWhere(squirrel.And{
				squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: groupID},
				squirrel.Eq{models.JOIN_REQUEST_USER_ID: principal.UserID},
			})

			existing, getErr := g.dao.GetGroupJoinRequest(ctx, daoOpts)
			if getErr != nil {
				return nil, getErr
			}

			if existing != nil && existing.Status == types.JoinRejected {
				return nil, ErrGroupJoinRequestRejected
			}

			return nil, ErrGroupJoinRequestPending
		}

		return nil, err
	}

	pending := types.JoinPending

	return groupResponseBuilder(group, nil, &pending), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Delete deletes a group and its associated data
func (g *Groups) Delete(ctx context.Context, groupID string) error {
	principal, err := principalFromContext(ctx)
	if err != nil {
		return err
	}

	if principal.SiteRole != types.SiteRoleAdmin {
		return ErrNotSiteAdmin
	}

	group, err := g.dao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
	if err != nil {
		return err
	}

	if group == nil {
		return ErrGroupNotFound
	}

	return g.dao.DeleteGroups(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// Response builders
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupsResponsesBuilder builds a slice of GroupResponse from group models
func groupsResponsesBuilder(groups []*models.Group, status userMemberStatus) []*GroupResponse {
	out := make([]*GroupResponse, len(groups))
	for i, group := range groups {
		groupRole, joinRequestStatus := status.forGroup(group.ID)
		out[i] = groupResponseBuilder(group, groupRole, joinRequestStatus)
	}

	return out
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupResponseBuilder builds a GroupResponse from a group model
func groupResponseBuilder(group *models.Group, groupRole *types.GroupRole, joinRequestStatus *types.JoinRequestStatus) *GroupResponse {
	return &GroupResponse{
		ID:                group.ID,
		CreatedAt:         group.CreatedAt,
		UpdatedAt:         group.UpdatedAt,
		Name:              group.Name,
		CreatedBy:         group.CreatedBy,
		MemberCount:       group.MemberCount,
		GroupRole:         groupRole,
		JoinRequestStatus: joinRequestStatus,
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// User member status
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// userMemberStatus queries a user's group role and join request status for a list of group IDs
func (g *Groups) getUserMemberStatus(ctx context.Context, userID string, groupIDs []string) (userMemberStatus, error) {
	if len(groupIDs) == 0 {
		return userMemberStatus{}, nil
	}

	where := squirrel.And{
		squirrel.Eq{models.GROUP_MEMBER_USER_ID: userID},
		squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: groupIDs},
	}

	members, err := g.dao.ListGroupMembers(ctx, dao.NewOptions().WithWhere(where))
	if err != nil {
		return userMemberStatus{}, err
	}

	roles := make(groupRolesByGroupID, len(members))
	for _, member := range members {
		roles[member.GroupID] = member.GroupRole
	}

	where = squirrel.And{
		squirrel.Eq{models.JOIN_REQUEST_USER_ID: userID},
		squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: groupIDs},
	}

	joinRequests, err := g.dao.ListGroupJoinRequests(ctx, dao.NewOptions().WithWhere(where))
	if err != nil {
		return userMemberStatus{}, err
	}

	requests := make(joinRequestStatusesByGroupID, len(joinRequests))
	for _, request := range joinRequests {
		requests[request.GroupID] = request.Status
	}

	return userMemberStatus{roles: roles, requests: requests}, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// forGroup returns the group role and join request status for a group ID
func (s userMemberStatus) forGroup(groupID string) (*types.GroupRole, *types.JoinRequestStatus) {
	if role := s.roles[groupID]; role.IsValid() {
		r := role
		return &r, nil
	}

	if status := s.requests[groupID]; status == types.JoinPending || status == types.JoinRejected {
		s := status
		return nil, &s
	}

	return nil, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getAdminSummary returns a count of pending and rejected join requests for a group
func (g *Groups) getAdminSummary(ctx context.Context, groupID string) (*GroupAdminSummary, error) {
	pending, err := g.dao.CountGroupJoinRequests(ctx, dao.NewOptions().WithWhere(squirrel.And{
		squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: groupID},
		squirrel.Eq{models.JOIN_REQUEST_STATUS: types.JoinPending},
	}))
	if err != nil {
		return nil, err
	}

	rejected, err := g.dao.CountGroupJoinRequests(ctx, dao.NewOptions().WithWhere(squirrel.And{
		squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: groupID},
		squirrel.Eq{models.JOIN_REQUEST_STATUS: types.JoinRejected},
	}))
	if err != nil {
		return nil, err
	}

	return &GroupAdminSummary{
		PendingJoinRequestCount:  pending,
		RejectedJoinRequestCount: rejected,
	}, nil
}
