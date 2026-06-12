package api

import (
	"github.com/geerew/friendle/service"
	"github.com/gofiber/fiber/v2"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// initGroupRoutes initializes the group routes
func (r *Router) initGroupRoutes() {
	groupRoutes := r.apiGroup("groups")

	groupRoutes.Post("/", r.requireAccess(accessSiteUser), r.createGroup)
	groupRoutes.Get("/self", r.requireAccess(accessSiteUser), r.listSelfGroups)
	groupRoutes.Get("/", r.requireAccess(accessSiteUser), r.listGroups)
	groupRoutes.Get("/:id", r.requireAccess(accessSiteUser), r.getGroup)
	groupRoutes.Delete("/:id", r.requireAccess(accessSiteAdmin), r.deleteGroup)

	// Members
	groupRoutes.Get("/:id/members", r.requireAccess(accessSiteUser), r.listGroupMembers)
	groupRoutes.Patch("/:id/members/:userId", r.requireAccess(accessSiteUser), r.updateGroupMemberRole)
	groupRoutes.Delete("/:id/members/:userId", r.requireAccess(accessSiteUser), r.removeGroupMember)

	// Round
	groupRoutes.Get("/:id/round/today", r.requireAccess(accessSiteUser), r.getGroupRoundToday)

	// Join
	groupRoutes.Post("/:id/join", r.requireAccess(accessSiteUser), r.createGroupJoinRequest)

	// Pending
	groupRoutes.Get("/:id/pending", r.requireAccess(accessSiteUser), r.listGroupPendingJoinRequests)
	groupRoutes.Post("/:id/pending/:userId/approve", r.requireAccess(accessSiteUser), r.approveGroupJoinRequest)
	groupRoutes.Post("/:id/pending/:userId/decline", r.requireAccess(accessSiteUser), r.declineGroupJoinRequest)

	// Rejected
	groupRoutes.Get("/:id/rejected", r.requireAccess(accessSiteUser), r.listGroupRejectedJoinRequests)

}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroup creates a friend group for the authenticated user
func (r *Router) createGroup(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	req := &service.CreateGroupRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	group, err := r.appSvc.Groups.Create(ctx, principal.UserID, *req)
	if err != nil {
		return serviceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroup returns a group by ID
func (r *Router) getGroup(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	groupID := c.Params("id")

	group, err := r.appSvc.Groups.Get(ctx, groupID, principal.UserID)
	if err != nil {
		return serviceError(c, err)
	}

	return c.JSON(group)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// listGroupMembers returns paginated group members for group members
func (r *Router) listGroupMembers(c *fiber.Ctx) error {
	groupID := c.Params("id")

	if !r.isGroupMember(c, groupID) {
		return nil
	}

	_, ctx := principalAndCtx(c)
	page := paginationFromCtx(c)

	members, err := r.appSvc.Groups.ListMembers(ctx, groupID, page)
	if err != nil {
		return serviceError(c, err)
	}

	pResult, err := page.BuildResult(members)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// updateGroupMemberRole updates a group member role for group admins
func (r *Router) updateGroupMemberRole(c *fiber.Ctx) error {
	groupID := c.Params("id")

	if !r.isGroupAdmin(c, groupID) {
		return nil
	}

	principal, ctx := principalAndCtx(c)

	req := &service.UpdateMemberRoleRequest{}
	if err := c.BodyParser(req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Error parsing data", err)
	}

	userID := c.Params("userId")

	if userID == principal.UserID {
		return errorResponse(c, fiber.StatusBadRequest, "Cannot modify your own group membership", nil)
	}

	member, err := r.appSvc.Groups.UpdateMemberRole(ctx, groupID, userID, *req)
	if err != nil {
		return serviceError(c, err)
	}

	return c.JSON(member)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// removeGroupMember removes a group member for group admins
func (r *Router) removeGroupMember(c *fiber.Ctx) error {
	groupID := c.Params("id")

	if !r.isGroupAdmin(c, groupID) {
		return nil
	}

	principal, ctx := principalAndCtx(c)
	userID := c.Params("userId")

	if userID == principal.UserID {
		return errorResponse(c, fiber.StatusBadRequest, "Cannot modify your own group membership", nil)
	}

	if err := r.appSvc.Groups.RemoveMember(ctx, groupID, userID); err != nil {
		return serviceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// getGroupRoundToday returns today's round state for a group member
func (r *Router) getGroupRoundToday(c *fiber.Ctx) error {
	groupID := c.Params("id")

	if !r.isGroupMember(c, groupID) {
		return nil
	}

	principal, ctx := principalAndCtx(c)

	roundToday, err := r.appSvc.Rounds.Today(ctx, groupID, principal.UserID)
	if err != nil {
		return serviceError(c, err)
	}

	return c.JSON(roundToday)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// listGroupPendingJoinRequests returns paginated pending join requests for group admins
func (r *Router) listGroupPendingJoinRequests(c *fiber.Ctx) error {
	groupID := c.Params("id")

	if !r.isGroupAdmin(c, groupID) {
		return nil
	}

	_, ctx := principalAndCtx(c)
	page := paginationFromCtx(c)

	requests, err := r.appSvc.Groups.ListPendingJoinRequests(ctx, groupID, page)
	if err != nil {
		return serviceError(c, err)
	}

	pResult, err := page.BuildResult(requests)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// listGroupRejectedJoinRequests returns paginated rejected join requests for group admins
func (r *Router) listGroupRejectedJoinRequests(c *fiber.Ctx) error {
	groupID := c.Params("id")

	if !r.isGroupAdmin(c, groupID) {
		return nil
	}

	_, ctx := principalAndCtx(c)
	page := paginationFromCtx(c)

	requests, err := r.appSvc.Groups.ListRejectedJoinRequests(ctx, groupID, page)
	if err != nil {
		return serviceError(c, err)
	}

	pResult, err := page.BuildResult(requests)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// approveGroupJoinRequest approves a pending join request for a group admin
func (r *Router) approveGroupJoinRequest(c *fiber.Ctx) error {
	groupID := c.Params("id")

	if !r.isGroupAdmin(c, groupID) {
		return nil
	}

	_, ctx := principalAndCtx(c)
	userID := c.Params("userId")

	if err := r.appSvc.Groups.ApproveJoinRequest(ctx, groupID, userID); err != nil {
		return serviceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// declineGroupJoinRequest rejects a pending join request for a group admin
func (r *Router) declineGroupJoinRequest(c *fiber.Ctx) error {
	groupID := c.Params("id")

	if !r.isGroupAdmin(c, groupID) {
		return nil
	}

	_, ctx := principalAndCtx(c)
	userID := c.Params("userId")

	if err := r.appSvc.Groups.DeclineJoinRequest(ctx, groupID, userID); err != nil {
		return serviceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// listGroups returns paginated groups, optionally filtered by name
func (r *Router) listGroups(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)
	page := paginationFromCtx(c)

	var (
		groups any
		err    error
	)

	if c.Context().QueryArgs().Has("name") {
		groups, err = r.appSvc.Groups.Search(ctx, principal.UserID, page, c.Query("name", ""))
	} else {
		groups, err = r.appSvc.Groups.List(ctx, page)
	}

	if err != nil {
		return serviceError(c, err)
	}

	pResult, err := page.BuildResult(groups)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// listSelfGroups returns paginated groups the authenticated user belongs to
func (r *Router) listSelfGroups(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	page := paginationFromCtx(c)
	groups, err := r.appSvc.Groups.ListSelf(ctx, principal.UserID, page)
	if err != nil {
		return serviceError(c, err)
	}

	pResult, err := page.BuildResult(groups)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Error building pagination result", err)
	}

	return c.JSON(pResult)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// deleteGroup deletes a group
func (r *Router) deleteGroup(c *fiber.Ctx) error {
	_, ctx := principalAndCtx(c)

	if err := r.appSvc.Groups.Delete(ctx, c.Params("id")); err != nil {
		return serviceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// createGroupJoinRequest creates a pending join request for the authenticated user
func (r *Router) createGroupJoinRequest(c *fiber.Ctx) error {
	principal, ctx := principalAndCtx(c)

	group, err := r.appSvc.Groups.RequestJoin(ctx, principal.UserID, c.Params("id"))
	if err != nil {
		return serviceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// isGroupMember writes an error response and returns false when the caller is not a group
// member
func (r *Router) isGroupMember(c *fiber.Ctx, groupID string) bool {
	principal, ctx := principalAndCtx(c)

	isMember, err := r.appSvc.Groups.IsMember(ctx, groupID, principal.UserID)
	if err != nil {
		_ = serviceError(c, err)
		return false
	}

	if !isMember {
		_ = errorResponse(c, fiber.StatusForbidden, "Forbidden", nil)
		return false
	}

	return true
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// isGroupAdmin writes an error response and returns false when the caller is not a group
// admin
func (r *Router) isGroupAdmin(c *fiber.Ctx, groupID string) bool {
	principal, ctx := principalAndCtx(c)

	isAdmin, err := r.appSvc.Groups.IsAdmin(ctx, groupID, principal.UserID)
	if err != nil {
		_ = serviceError(c, err)
		return false
	}

	if !isAdmin {
		_ = errorResponse(c, fiber.StatusForbidden, "Forbidden", nil)
		return false
	}

	return true
}
