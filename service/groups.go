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
const minPlayableMembers = 2

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroupRequest represents a group create request
type CreateGroupRequest struct {
	Name string `json:"name"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupResponse represents a group response
type GroupResponse struct {
	ID                 string                   `json:"id"`
	CreatedAt          types.DateTime           `json:"createdAt"`
	UpdatedAt          types.DateTime           `json:"updatedAt"`
	Name               string                   `json:"name"`
	CreatedBy          string                   `json:"createdBy,omitempty"`
	MemberCount        int                      `json:"memberCount"`
	MemberThresholdMet bool                     `json:"memberThresholdMet"`
	GroupRole          *types.GroupRole         `json:"groupRole,omitempty"`
	JoinRequestStatus  *types.JoinRequestStatus `json:"joinRequestStatus,omitempty"`
	AdminSummary       *GroupAdminSummary       `json:"adminSummary,omitempty"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupAdminSummary holds join request counts for a group admin
type GroupAdminSummary struct {
	PendingJoinRequestCount  int `json:"pendingJoinRequestCount"`
	RejectedJoinRequestCount int `json:"rejectedJoinRequestCount"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupMemberResponse represents a group member in list responses
type GroupMemberResponse struct {
	UserID      string          `json:"userId"`
	DisplayName string          `json:"displayName"`
	GroupRole   types.GroupRole `json:"groupRole"`
	TimesPicked int             `json:"timesPicked"`
	PickerSkips int             `json:"pickerSkips"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateMemberRoleRequest represents a group member role update request
type UpdateMemberRoleRequest struct {
	GroupRole types.GroupRole `json:"groupRole"`
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupJoinRequestResponse represents a pending join request in list responses
type GroupJoinRequestResponse struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
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

// Create creates a group and sets userID as group admin
func (g *Groups) Create(ctx context.Context, userID string, req CreateGroupRequest) (*GroupResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, ErrGroupNameRequired
	}

	if utf8.RuneCountInString(name) > maxGroupNameLength {
		return nil, ErrGroupNameTooLong
	}

	group := &models.Group{
		Name:      name,
		CreatedBy: userID,
	}

	err := g.dao.RunInTransaction(ctx, func(txCtx context.Context) error {
		if err := g.dao.CreateGroup(txCtx, group); err != nil {
			if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
				return ErrGroupNameTaken
			}

			return err
		}

		member := &models.GroupMember{
			GroupID:   group.ID,
			UserID:    userID,
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

// ListSelf returns a paginated slice of groups userID belongs to
func (g *Groups) ListSelf(ctx context.Context, userID string, page *pagination.Pagination) ([]*GroupResponse, error) {
	where, err := dao.MemberGroupsWhere(userID)
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
	status, err := g.getUserMemberStatus(ctx, userID, groupIDs)
	if err != nil {
		return nil, err
	}

	return groupsResponsesBuilder(groups, status), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Search returns a paginated slice of groups whose names contain the search term
//
// Response is ordered by prefix matches first, then substring matches.
//
// Results include userID's group role and join request status for each group
func (g *Groups) Search(ctx context.Context, userID string, page *pagination.Pagination, name string) ([]*GroupResponse, error) {
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
	status, err := g.getUserMemberStatus(ctx, userID, groupIDs)
	if err != nil {
		return nil, err
	}

	return groupsResponsesBuilder(groups, status), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Get returns a group by ID enriched with userID's membership and join request status
func (g *Groups) Get(ctx context.Context, groupID, userID string) (*GroupResponse, error) {
	group, err := g.dao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, ErrGroupNotFound
	}

	status, err := g.getUserMemberStatus(ctx, userID, []string{groupID})
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

// ListMembers returns paginated group members with display names
func (g *Groups) ListMembers(ctx context.Context, groupID string, page *pagination.Pagination) ([]*GroupMemberResponse, error) {
	group, err := g.dao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, ErrGroupNotFound
	}

	// List group members and order by admin first, then display name
	daoOpts := dao.NewOptions().
		WithWhere(squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: groupID}).
		WithPagination(page).
		WithOrderByClause(squirrel.Expr(
			"CASE WHEN " + models.GROUP_MEMBER_TABLE_GROUP_ROLE + " = 'group_admin' THEN 0 ELSE 1 END, LOWER(" + models.USER_TABLE_DISPLAY_NAME + ") ASC",
		))

	members, err := g.dao.ListGroupMembers(ctx, daoOpts)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return []*GroupMemberResponse{}, nil
	}

	return groupMemberResponsesBuilder(members), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// UpdateMemberRole updates a group member's role
//
// At least one group admin must remain
func (g *Groups) UpdateMemberRole(ctx context.Context, groupID, userID string, req UpdateMemberRoleRequest) (*GroupMemberResponse, error) {
	if !req.GroupRole.IsValid() {
		return nil, ErrNoUpdateData
	}

	dbOpts := dao.NewOptions().WithWhere(squirrel.And{
		squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: groupID},
		squirrel.Eq{models.GROUP_MEMBER_USER_ID: userID},
	})

	member, err := g.dao.GetGroupMember(ctx, dbOpts)
	if err != nil {
		return nil, err
	}

	if member == nil {
		return nil, ErrGroupMemberNotFound
	}

	if member.GroupRole == req.GroupRole {
		return groupMemberResponseBuilder(member), nil
	}

	if member.GroupRole == types.GroupRoleAdmin && req.GroupRole == types.GroupRoleUser {
		dbOpts := dao.NewOptions().WithWhere(squirrel.And{
			squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: groupID},
			squirrel.Eq{models.GROUP_MEMBER_GROUP_ROLE: types.GroupRoleAdmin},
		})

		adminCount, err := g.dao.CountGroupMembers(ctx, dbOpts)
		if err != nil {
			return nil, err
		}

		if adminCount <= 1 {
			return nil, ErrGroupLastAdmin
		}
	}

	member.GroupRole = req.GroupRole
	if err := g.dao.UpdateGroupMember(ctx, member); err != nil {
		return nil, err
	}

	return groupMemberResponseBuilder(member), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RemoveMember removes a member from a group
//
// The last group admin cannot be removed
func (g *Groups) RemoveMember(ctx context.Context, groupID, userID string) error {
	dbOpts := dao.NewOptions().WithWhere(squirrel.And{
		squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: groupID},
		squirrel.Eq{models.GROUP_MEMBER_USER_ID: userID},
	})

	member, err := g.dao.GetGroupMember(ctx, dbOpts)
	if err != nil {
		return err
	}

	if member == nil {
		return ErrGroupMemberNotFound
	}

	if member.GroupRole == types.GroupRoleAdmin {
		adminOpts := dao.NewOptions().WithWhere(squirrel.And{
			squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: groupID},
			squirrel.Eq{models.GROUP_MEMBER_GROUP_ROLE: types.GroupRoleAdmin},
		})

		adminCount, err := g.dao.CountGroupMembers(ctx, adminOpts)
		if err != nil {
			return err
		}

		if adminCount <= 1 {
			return ErrGroupLastAdmin
		}
	}

	joinOpts := dao.NewOptions().WithWhere(squirrel.And{
		squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: groupID},
		squirrel.Eq{models.JOIN_REQUEST_USER_ID: userID},
	})

	if err := g.dao.DeleteGroupJoinRequests(ctx, joinOpts); err != nil {
		return err
	}

	return g.dao.DeleteGroupMembers(ctx, dbOpts)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListPendingJoinRequests returns paginated pending join requests
func (g *Groups) ListPendingJoinRequests(ctx context.Context, groupID string, page *pagination.Pagination) ([]*GroupJoinRequestResponse, error) {
	daoOpts := dao.NewOptions().
		WithWhere(squirrel.And{
			squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: groupID},
			squirrel.Eq{models.JOIN_REQUEST_STATUS: types.JoinPending},
		}).
		WithPagination(page).
		WithOrderByClause(squirrel.Expr("LOWER(" + models.USER_TABLE_DISPLAY_NAME + ") ASC"))

	requests, err := g.dao.ListGroupJoinRequests(ctx, daoOpts)
	if err != nil {
		return nil, err
	}

	if len(requests) == 0 {
		return []*GroupJoinRequestResponse{}, nil
	}

	return groupJoinRequestResponsesBuilder(requests), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListRejectedJoinRequests returns paginated rejected join requests
func (g *Groups) ListRejectedJoinRequests(ctx context.Context, groupID string, page *pagination.Pagination) ([]*GroupJoinRequestResponse, error) {
	daoOpts := dao.NewOptions().
		WithWhere(squirrel.And{
			squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: groupID},
			squirrel.Eq{models.JOIN_REQUEST_STATUS: types.JoinRejected},
		}).
		WithPagination(page).
		WithOrderByClause(squirrel.Expr("LOWER(" + models.USER_TABLE_DISPLAY_NAME + ") ASC"))

	requests, err := g.dao.ListGroupJoinRequests(ctx, daoOpts)
	if err != nil {
		return nil, err
	}

	if len(requests) == 0 {
		return []*GroupJoinRequestResponse{}, nil
	}

	return groupJoinRequestResponsesBuilder(requests), nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ApproveJoinRequest approves a pending join request and adds the user as a group member
func (g *Groups) ApproveJoinRequest(ctx context.Context, groupID, userID string) error {
	return g.dao.RunInTransaction(ctx, func(txCtx context.Context) error {
		request, err := g.getPendingJoinRequest(txCtx, groupID, userID)
		if err != nil {
			return err
		}

		request.Status = types.JoinApproved
		if err := g.dao.UpdateGroupJoinRequest(txCtx, request); err != nil {
			return err
		}

		return g.dao.CreateGroupMember(txCtx, &models.GroupMember{
			GroupID:   groupID,
			UserID:    userID,
			GroupRole: types.GroupRoleUser,
		})
	})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// DeclineJoinRequest rejects a pending join request
func (g *Groups) DeclineJoinRequest(ctx context.Context, groupID, userID string) error {
	request, err := g.getPendingJoinRequest(ctx, groupID, userID)
	if err != nil {
		return err
	}

	request.Status = types.JoinRejected

	return g.dao.UpdateGroupJoinRequest(ctx, request)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RequestJoin creates a pending join request for userID
func (g *Groups) RequestJoin(ctx context.Context, userID, groupID string) (*GroupResponse, error) {
	group, err := g.dao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_ID: groupID}))
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, ErrGroupNotFound
	}

	status, err := g.getUserMemberStatus(ctx, userID, []string{groupID})
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
		UserID:  userID,
		Status:  types.JoinPending,
	})
	if err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			daoOpts := dao.NewOptions().WithWhere(squirrel.And{
				squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: groupID},
				squirrel.Eq{models.JOIN_REQUEST_USER_ID: userID},
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

// MemberRole returns a user's role in a group, or nil when they are not a member
func (g *Groups) MemberRole(ctx context.Context, groupID, userID string) (*types.GroupRole, error) {
	dbOpts := dao.NewOptions().WithWhere(squirrel.And{
		squirrel.Eq{models.GROUP_MEMBER_GROUP_ID: groupID},
		squirrel.Eq{models.GROUP_MEMBER_USER_ID: userID},
	})

	member, err := g.dao.GetGroupMember(ctx, dbOpts)
	if err != nil {
		return nil, err
	}

	if member == nil || !member.GroupRole.IsValid() {
		return nil, nil
	}

	role := member.GroupRole

	return &role, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsMember reports whether user is a group member
func (g *Groups) IsMember(ctx context.Context, groupID, userID string) (bool, error) {
	role, err := g.MemberRole(ctx, groupID, userID)
	if err != nil {
		return false, err
	}

	return role != nil, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsAdmin reports whether user is a group admin
func (g *Groups) IsAdmin(ctx context.Context, groupID, userID string) (bool, error) {
	role, err := g.MemberRole(ctx, groupID, userID)
	if err != nil {
		return false, err
	}

	return role != nil && *role == types.GroupRoleAdmin, nil
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
		ID:                 group.ID,
		CreatedAt:          group.CreatedAt,
		UpdatedAt:          group.UpdatedAt,
		Name:               group.Name,
		CreatedBy:          group.CreatedBy,
		MemberCount:        group.MemberCount,
		MemberThresholdMet: group.MemberCount >= minPlayableMembers,
		GroupRole:          groupRole,
		JoinRequestStatus:  joinRequestStatus,
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupMemberResponsesBuilder builds a slice of GroupMemberResponse from group member models
func groupMemberResponsesBuilder(members []*models.GroupMember) []*GroupMemberResponse {
	out := make([]*GroupMemberResponse, len(members))
	for i, member := range members {
		out[i] = groupMemberResponseBuilder(member)
	}

	return out
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupMemberResponseBuilder builds a GroupMemberResponse from a group member model
func groupMemberResponseBuilder(member *models.GroupMember) *GroupMemberResponse {
	return &GroupMemberResponse{
		UserID:      member.UserID,
		DisplayName: member.DisplayName,
		GroupRole:   member.GroupRole,
		TimesPicked: member.TimesPicked,
		PickerSkips: member.PickerSkips,
	}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupJoinRequestResponsesBuilder builds a slice of GroupJoinRequestResponse from join request
// models
func groupJoinRequestResponsesBuilder(requests []*models.GroupJoinRequest) []*GroupJoinRequestResponse {
	out := make([]*GroupJoinRequestResponse, len(requests))
	for i, request := range requests {
		out[i] = &GroupJoinRequestResponse{
			UserID:      request.UserID,
			DisplayName: request.DisplayName,
		}
	}

	return out
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

// getPendingJoinRequest loads a pending join request for a group and user
func (g *Groups) getPendingJoinRequest(ctx context.Context, groupID, userID string) (*models.GroupJoinRequest, error) {
	request, err := g.dao.GetGroupJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.And{
		squirrel.Eq{models.JOIN_REQUEST_GROUP_ID: groupID},
		squirrel.Eq{models.JOIN_REQUEST_USER_ID: userID},
		squirrel.Eq{models.JOIN_REQUEST_STATUS: types.JoinPending},
	}))
	if err != nil {
		return nil, err
	}

	if request == nil {
		return nil, ErrGroupJoinRequestNotFound
	}

	return request, nil
}
