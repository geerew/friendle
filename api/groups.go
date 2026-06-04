package api

import (
	"context"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initGroupRoutes initializes the group routes
func (r *Router) initGroupRoutes() {
	g := r.apiGroup("groups")

	// Groups
	g.Post("/", r.requireAccess(accessSiteUser), r.createGroup)
	g.Get("/", r.requireAccess(accessSiteUser), r.getGroups)
	g.Get("/mine", r.requireAccess(accessSiteUser), r.getMyGroups)
	g.Get("/search", r.requireAccess(accessSiteUser), r.searchGroups)
	g.Get("/:id", r.requireAccess(accessGroupMemberScope), r.getGroup)
	g.Patch("/:id", r.requireAccess(accessGroupAdminScope), r.updateGroup)
	g.Delete("/:id", r.requireAccess(accessGroupAdminScope), r.deleteGroup)

	// Join requests
	g.Post("/:id/join-requests", r.requireAccess(accessSiteUser), r.createGroupJoinRequest)
	g.Delete("/:id/join-requests/:userId", r.requireAccess(accessSiteUser), r.deleteGroupJoinRequest)
	g.Get("/:id/join-requests", r.requireAccess(accessGroupAdminScope), r.getGroupJoinRequests)
	g.Post("/:id/join-requests/:rid/approve", r.requireAccess(accessGroupAdminScope), r.updateGroupJoinRequestApprove)
	g.Post("/:id/join-requests/:rid/reject", r.requireAccess(accessGroupAdminScope), r.updateGroupJoinRequestReject)

	// Members
	g.Delete("/:id/members/:userId", r.requireAccess(accessGroupAdminScope), r.deleteGroupMember)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroup creates a group and adds the caller as group admin
func (r *Router) createGroup(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	req := &createGroupRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Name required", nil)
	}

	// Check if the name is already taken
	if taken, err := r.groupNameTaken(ctx, name, ""); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to create group", err)
	} else if taken {
		return errorResponse(c, fiber.StatusBadRequest, "Group name already exists", nil)
	}

	group := &models.Group{
		Name:      name,
		CreatedBy: principal.UserID,
	}

	// Create the group
	if err := r.appDao.CreateGroup(ctx, group); err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return errorResponse(c, fiber.StatusBadRequest, "Group name already exists", nil)
		}

		return errorResponse(c, fiber.StatusInternalServerError, "Failed to create group", err)
	}

	// Add the user as group admin
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

// getGroups returns a paginated site-wide group list
func (r *Router) getGroups(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	dbOpts := dao.NewOptions().
		WithOrderBy(utils.StringSplit(c.Query("orderBy", ""), ",")...).
		WithPagination(paginationFromCtx(c)).
		WithMembers().
		WithJoinRequests().
		WithMemberCount()

	groups, err := r.appDao.ListGroups(ctx, dbOpts)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "List failed", err)
	}

	pResult, err := dbOpts.Pagination.BuildResult(adminGroupResponseHelper(groups))
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getMyGroups returns groups the caller belongs to
func (r *Router) getMyGroups(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	members, err := r.appDao.ListGroupMembers(ctx, dao.NewOptions().
		WithWhere(squirrel.Eq{models.GROUP_MEMBER_USER_ID: principal.UserID}).
		WithGroup().
		WithMemberCount())
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to list groups", err)
	}

	return c.JSON(userGroupSummaryResponsesFromMembers(members))
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

	dbOpts := dao.NewOptions().
		WithPagination(paginationFromCtx(c)).
		WithGroupNameSearch(strings.ToLower(q)).
		WithMemberCount()
	groups, err := r.appDao.ListGroups(ctx, dbOpts)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Search failed", err)
	}

	groupIDs := make([]string, len(groups))
	for i, group := range groups {
		groupIDs[i] = group.ID
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
		groupSearchResponsesFromGroups(groups, stringSet(memberGroupIDs), stringSet(pendingGroupIDs)),
	)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroup returns group settings, members, and join requests for a member
func (r *Router) getGroup(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	groupID := c.Params("id")
	m, _ := r.appDao.GetGroupMember(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
		"group_id": groupID,
		"user_id":  principal.UserID,
	}))
	role := types.GroupRoleUser
	isSiteAdmin := principal.SiteRole == types.SiteRoleAdmin
	if m != nil {
		role = m.GroupRole
	}

	loadJoinRequests := isSiteAdmin || (m != nil && m.GroupRole == types.GroupRoleAdmin)

	groupOpts := dao.NewOptions().
		WithWhere(squirrel.Eq{models.BASE_ID: groupID}).
		WithMembers()
	if loadJoinRequests {
		groupOpts = groupOpts.WithJoinRequests()
	}

	group, err := r.appDao.GetGroup(ctx, groupOpts)
	if err != nil || group == nil {
		return errorResponse(c, fiber.StatusNotFound, "Group not found", nil)
	}

	return c.JSON(groupResponseHelper(group, role))
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// updateGroup updates a group
func (r *Router) updateGroup(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	groupID := c.Params("id")
	g, err := r.appDao.GetGroup(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: groupID}))
	if err != nil || g == nil {
		return errorResponse(c, fiber.StatusNotFound, "Group not found", nil)
	}

	req := &updateGroupRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid body", nil)
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return errorResponse(c, fiber.StatusBadRequest, "Name required", nil)
		}

		taken, err := r.groupNameTaken(ctx, name, groupID)
		if err != nil {
			return errorResponse(c, fiber.StatusInternalServerError, "Update failed", err)
		}

		if taken {
			return errorResponse(c, fiber.StatusBadRequest, "Group name already exists", nil)
		}

		g.Name = name
	}

	if err := r.appDao.UpdateGroup(ctx, g); err != nil {
		if strings.HasPrefix(err.Error(), "UNIQUE constraint failed") {
			return errorResponse(c, fiber.StatusBadRequest, "Group name already exists", nil)
		}

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

	if m, _ := r.appDao.GetGroupMember(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
		"group_id": groupID,
		"user_id":  principal.UserID,
	})); m != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Already a member", nil)
	}

	if existing, _ := r.appDao.GetJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
		"group_id": groupID,
		"user_id":  principal.UserID,
	})); existing != nil && existing.Status == types.JoinPending {
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
			m, err := r.appDao.GetGroupMember(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
				"group_id": groupID,
				"user_id":  principal.UserID,
			}))
			if err != nil || m == nil || m.GroupRole != types.GroupRoleAdmin {
				return errorResponse(c, fiber.StatusForbidden, "Forbidden", nil)
			}
		}
	}

	jr, err := r.appDao.GetJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
		"group_id": groupID,
		"user_id":  targetUserID,
	}))
	if err != nil || jr == nil || jr.Status != types.JoinPending {
		return errorResponse(c, fiber.StatusNotFound, "Pending request not found", nil)
	}

	if err := r.appDao.DeleteJoinRequests(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
		models.BASE_ID: jr.ID,
	})); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Cancel failed", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroupJoinRequests returns pending join requests for a group
func (r *Router) getGroupJoinRequests(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	list, err := r.appDao.ListJoinRequests(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
		"group_id": c.Params("id"),
		"status":   types.JoinPending,
	}))
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

	if err := r.appDao.DeleteGroupMembers(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
		models.GROUP_MEMBER_GROUP_ID: c.Params("id"),
		models.GROUP_MEMBER_USER_ID:  c.Params("userId"),
	})); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Remove failed", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteGroup deletes a group
func (r *Router) deleteGroup(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: c.Params("id")})
	if err := r.appDao.DeleteGroups(ctx, dbOpts); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error deleting group", err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// groupNameTaken reports whether another group already uses the given name
func (r *Router) groupNameTaken(ctx context.Context, name, excludeGroupID string) (bool, error) {
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.GROUP_TABLE_NAME: name})
	existing, err := r.appDao.GetGroup(ctx, dbOpts)
	if err != nil || existing == nil {
		return false, err
	}

	return existing.ID != excludeGroupID, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// resolveGroupJoinRequest updates a join request status and optionally adds a member
func (r *Router) resolveGroupJoinRequest(c *fiber.Ctx, status types.JoinRequestStatus, role types.GroupRole) error {
	_, ctx := principalAndCtx(c)

	jr, err := r.appDao.GetJoinRequest(ctx, dao.NewOptions().WithWhere(squirrel.Eq{models.BASE_ID: c.Params("rid")}))
	if err != nil || jr == nil {
		return errorResponse(c, fiber.StatusNotFound, "Request not found", nil)
	}

	jr.Status = status
	if err := r.appDao.UpdateJoinRequest(ctx, jr); err != nil {
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
	m, err := r.appDao.GetGroupMember(ctx, dao.NewOptions().WithWhere(squirrel.Eq{
		"group_id": groupID,
		"user_id":  userID,
	}))
	if err != nil || m == nil {
		return nil, err
	}

	return m, nil
}
