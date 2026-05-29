package api

import (
	"context"
	"strings"

	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initGroupRoutes initializes the group routes
func (r *Router) initGroupRoutes() {
	g := r.apiGroup("groups")

	// Groups
	g.Post("/", r.requireAccess(accessAuth), r.createGroup)
	g.Get("/", r.requireAccess(accessAuth), r.getMyGroups)
	g.Get("/search", r.requireAccess(accessAuth), r.searchGroups)
	g.Get("/:id", r.requireAccess(accessGroupMemberScope), r.getGroup)
	g.Patch("/:id", r.requireAccess(accessGroupAdminScope), r.updateGroup)

	// Join requests
	g.Post("/:id/join-requests", r.requireAccess(accessAuth), r.createGroupJoinRequest)
	g.Delete("/:id/join-requests/:userId", r.requireAccess(accessAuth), r.deleteGroupJoinRequest)
	g.Get("/:id/join-requests", r.requireAccess(accessGroupAdminScope), r.getGroupJoinRequests)
	g.Post("/:id/join-requests/:rid/approve", r.requireAccess(accessGroupAdminScope), r.updateGroupJoinRequestApprove)
	g.Post("/:id/join-requests/:rid/reject", r.requireAccess(accessGroupAdminScope), r.updateGroupJoinRequestReject)

	// Members
	g.Delete("/:id/members/:userId", r.requireAccess(accessGroupAdminScope), r.deleteGroupMember)

	// Leaderboard
	g.Get("/:id/leaderboard", r.requireAccess(accessGroupMemberScope), r.getGroupLeaderboard)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroup creates a group and adds the caller as group admin
func (r *Router) createGroup(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	req := &createGroupRequest{}
	if err := c.BodyParser(req); err != nil || strings.TrimSpace(req.Name) == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Name required", nil)
	}

	group := &models.Group{
		Name:          strings.TrimSpace(req.Name),
		CreatedBy:     principal.UserID,
		IntervalHours: 24,
		Timezone:      "UTC",
	}
	if err := r.appDao.CreateGroup(ctx, group); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to create group", err)
	}

	member := &models.GroupMember{
		GroupID:   group.ID,
		UserID:    principal.UserID,
		GroupRole: types.GroupRoleAdmin,
	}
	if err := r.appDao.CreateGroupMember(ctx, member); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to add member", err)
	}

	return c.Status(fiber.StatusCreated).JSON(groupResponseHelper(group, types.GroupRoleAdmin))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getMyGroups returns groups the caller belongs to
func (r *Router) getMyGroups(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	rows, err := r.appDao.ListUserGroupSummariesForUserIDs(ctx, []string{principal.UserID})
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to list groups", err)
	}

	return c.JSON(userGroupSummaryResponsesFromRows(rows))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// searchGroups searches groups by name
func (r *Router) searchGroups(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		pResult, err := paginationFromCtx(c).BuildResult([]*groupSearchResponse{})
		if err != nil {
			return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
		}

		return c.JSON(pResult)
	}

	dbOpts := dao.NewOptions().WithPagination(paginationFromCtx(c))
	rows, err := r.appDao.SearchGroupSummaries(ctx, strings.ToLower(q), dbOpts)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Search failed", err)
	}

	groupIDs := make([]string, len(rows))
	for i, row := range rows {
		groupIDs[i] = row.ID
	}

	memberGroupIDs, err := r.appDao.ListMemberGroupIDsForUser(ctx, principal.UserID, groupIDs)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Search failed", err)
	}

	pendingGroupIDs, err := r.appDao.ListPendingJoinGroupIDsForUser(ctx, principal.UserID, groupIDs)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Search failed", err)
	}

	pResult, err := dbOpts.Pagination.BuildResult(
		groupSearchResponsesFromRows(rows, stringSet(memberGroupIDs), stringSet(pendingGroupIDs)),
	)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroup returns group detail and round summary for a member
func (r *Router) getGroup(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	groupID := c.Params("id")
	m, _ := r.appDao.GetGroupMember(ctx, groupID, principal.UserID)
	role := types.GroupRoleUser
	isAdmin := principal.SiteRole == types.SiteRoleAdmin
	if m != nil {
		role = m.GroupRole
	}

	load := dao.GroupDetailLoad{
		Members:        true,
		Leaderboard:    true,
		CurrentRound:   true,
		PreviousRounds: true,
		RoundLimit:     10,
	}
	if isAdmin || (m != nil && m.GroupRole == types.GroupRoleAdmin) {
		load.JoinRequests = true
	}

	detail, err := r.appDao.GetGroupDetail(ctx, groupID, load)
	if err != nil || detail == nil || detail.Group == nil {
		return errorResponse(c, fiber.StatusNotFound, "Group not found", nil)
	}

	resp := groupResponseHelper(detail.Group, role)
	resp.Round = r.roundSummary(ctx, detail.Group, principal.UserID)

	return c.JSON(resp)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// updateGroup updates a group
func (r *Router) updateGroup(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	groupID := c.Params("id")
	g, err := r.appDao.GetGroup(ctx, groupID)
	if err != nil || g == nil {
		return errorResponse(c, fiber.StatusNotFound, "Group not found", nil)
	}

	req := &updateGroupRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid body", nil)
	}

	if req.Name != nil {
		g.Name = strings.TrimSpace(*req.Name)
	}

	if req.IntervalHours != nil {
		g.IntervalHours = *req.IntervalHours
	}

	if err := r.appDao.UpdateGroup(ctx, g); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Update failed", err)
	}

	return c.JSON(groupResponseHelper(g, types.GroupRoleAdmin))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroupJoinRequest creates a pending join request for the caller
func (r *Router) createGroupJoinRequest(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	groupID := c.Params("id")

	req := &createGroupJoinRequest{}
	if len(c.Body()) > 0 {
		if err := c.BodyParser(req); err != nil {
			return errorResponse(c, fiber.StatusBadRequest, "Invalid body", nil)
		}

		if uid := strings.TrimSpace(req.UserID); uid != "" && uid != principal.UserID {
			return errorResponse(c, fiber.StatusForbidden, "Forbidden", nil)
		}
	}

	if m, _ := r.appDao.GetGroupMember(ctx, groupID, principal.UserID); m != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Already a member", nil)
	}

	if existing, _ := r.appDao.GetJoinRequestByUser(ctx, groupID, principal.UserID); existing != nil && existing.Status == types.JoinPending {
		return c.Status(fiber.StatusOK).JSON(&joinRequestResponse{Status: types.JoinPending})
	}

	jr := &models.GroupJoinRequest{GroupID: groupID, UserID: principal.UserID, Status: types.JoinPending}
	if err := r.appDao.CreateJoinRequest(ctx, jr); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Request failed", err)
	}

	return c.Status(fiber.StatusCreated).JSON(&joinRequestResponse{ID: jr.ID, Status: jr.Status})
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteGroupJoinRequest cancels a pending join request
func (r *Router) deleteGroupJoinRequest(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	groupID := c.Params("id")
	targetUserID := c.Params("userId")
	if targetUserID == "me" {
		targetUserID = principal.UserID
	}

	if targetUserID != principal.UserID {
		if principal.SiteRole != types.SiteRoleAdmin {
			m, err := r.appDao.GetGroupMember(ctx, groupID, principal.UserID)
			if err != nil || m == nil || m.GroupRole != types.GroupRoleAdmin {
				return errorResponse(c, fiber.StatusForbidden, "Forbidden", nil)
			}
		}
	}

	jr, err := r.appDao.GetJoinRequestByUser(ctx, groupID, targetUserID)
	if err != nil || jr == nil || jr.Status != types.JoinPending {
		return errorResponse(c, fiber.StatusNotFound, "Pending request not found", nil)
	}

	if err := r.appDao.DeletePendingJoinRequest(ctx, groupID, targetUserID); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Cancel failed", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroupJoinRequests returns pending join requests for a group
func (r *Router) getGroupJoinRequests(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	list, err := r.appDao.ListPendingJoinRequests(ctx, c.Params("id"))
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "List failed", err)
	}

	return c.JSON(list)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// updateGroupJoinRequestApprove approves a pending join request
func (r *Router) updateGroupJoinRequestApprove(c *fiber.Ctx) error {
	return r.resolveGroupJoinRequest(c, types.JoinApproved, types.GroupRoleUser)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// updateGroupJoinRequestReject rejects a pending join request
func (r *Router) updateGroupJoinRequestReject(c *fiber.Ctx) error {
	return r.resolveGroupJoinRequest(c, types.JoinRejected, "")
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteGroupMember removes a member from a group
func (r *Router) deleteGroupMember(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	if err := r.appDao.DeleteGroupMember(ctx, c.Params("id"), c.Params("userId")); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Remove failed", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroupLeaderboard returns the leaderboard for a group
func (r *Router) getGroupLeaderboard(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	entries, err := r.appDao.ListLeaderboardEntries(ctx, c.Params("id"))
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Leaderboard failed", err)
	}

	return c.JSON(leaderboardEntryResponsesFromModels(entries))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// resolveGroupJoinRequest updates a join request status and optionally adds a member
func (r *Router) resolveGroupJoinRequest(c *fiber.Ctx, status types.JoinRequestStatus, role types.GroupRole) error {
	_, ctx := principalAndCtx(c)

	jr, err := r.appDao.GetJoinRequest(ctx, c.Params("rid"))
	if err != nil || jr == nil {
		return errorResponse(c, fiber.StatusNotFound, "Request not found", nil)
	}

	if err := r.appDao.UpdateJoinRequestStatus(ctx, jr.ID, status); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Update failed", err)
	}

	if status == types.JoinApproved && role != "" {
		_ = r.appDao.CreateGroupMember(ctx, &models.GroupMember{
			GroupID: jr.GroupID, UserID: jr.UserID, GroupRole: role,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// membership returns group membership for a user, or an error when not a member
func (r *Router) membership(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	m, err := r.appDao.GetGroupMember(ctx, groupID, userID)
	if err != nil || m == nil {
		return nil, err
	}

	return m, nil
}
