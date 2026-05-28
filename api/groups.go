package api

import (
	"strings"

	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) initGroupRoutes() {
	g := r.apiGroup("groups")
	g.Post("/", r.requireAuth, r.createGroup)
	g.Get("/", r.requireAuth, r.listMyGroups)
	g.Get("/search", r.requireAuth, r.searchGroups)
	g.Get("/:id", r.requireAuth, r.getGroup)
	g.Patch("/:id", r.requireAuth, r.requireGroupAdmin, r.updateGroup)
	g.Post("/:id/join-requests", r.requireAuth, r.createJoinRequest)
	g.Get("/:id/join-requests", r.requireAuth, r.requireGroupAdmin, r.listJoinRequests)
	g.Post("/:id/join-requests/:rid/approve", r.requireAuth, r.requireGroupAdmin, r.approveJoinRequest)
	g.Post("/:id/join-requests/:rid/reject", r.requireAuth, r.requireGroupAdmin, r.rejectJoinRequest)
	g.Delete("/:id/members/:userId", r.requireAuth, r.requireGroupAdmin, r.removeMember)
	g.Get("/:id/leaderboard", r.requireAuth, r.requireGroupMember, r.leaderboard)
	r.initRoundRoutes(g)
}

type createGroupRequest struct {
	Name string `json:"name"`
}

func (r *Router) createGroup(c *fiber.Ctx) error {
	principal, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	req := &createGroupRequest{}
	if err := c.BodyParser(req); err != nil || strings.TrimSpace(req.Name) == "" {
		return errorResponse(c, fiber.StatusBadRequest, "Name required", nil)
	}
	group := &models.Group{Name: strings.TrimSpace(req.Name), CreatedBy: principal.UserID, IntervalHours: 24, Timezone: "UTC"}
	if err := r.appDao.CreateGroup(ctx, group); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to create group", err)
	}
	member := &models.GroupMember{GroupID: group.ID, UserID: principal.UserID, GroupRole: types.GroupRoleAdmin}
	if err := r.appDao.CreateGroupMember(ctx, member); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to add member", err)
	}
	return c.Status(fiber.StatusCreated).JSON(groupResponse(group, types.GroupRoleAdmin))
}

func (r *Router) listMyGroups(c *fiber.Ctx) error {
	principal, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}

	rows, err := r.appDao.ListUserGroupSummariesForUserIDs(ctx, []string{principal.UserID})
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to list groups", err)
	}

	return c.JSON(userGroupSummaryResponsesFromRows(rows))
}

func (r *Router) searchGroups(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		return c.JSON([]fiber.Map{})
	}
	groups, err := r.appDao.SearchGroups(ctx, strings.ToLower(q))
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Search failed", err)
	}
	var out []fiber.Map
	for _, g := range groups {
		out = append(out, fiber.Map{"id": g.ID, "name": g.Name})
	}
	return c.JSON(out)
}

func (r *Router) getGroup(c *fiber.Ctx) error {
	principal, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	groupID := c.Params("id")
	if _, err := r.membership(ctx, groupID, principal.UserID); err != nil {
		return errorResponse(c, fiber.StatusForbidden, "Not a member", nil)
	}
	g, err := r.appDao.GetGroup(ctx, groupID)
	if err != nil || g == nil {
		return errorResponse(c, fiber.StatusNotFound, "Group not found", nil)
	}
	m, _ := r.appDao.GetGroupMember(ctx, groupID, principal.UserID)
	role := types.GroupRoleUser
	if m != nil {
		role = m.GroupRole
	}
	summary := r.roundSummary(ctx, g, principal.UserID)
	resp := groupResponse(g, role)
	resp["round"] = summary
	return c.JSON(resp)
}

func (r *Router) updateGroup(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	groupID := c.Params("id")
	g, err := r.appDao.GetGroup(ctx, groupID)
	if err != nil || g == nil {
		return errorResponse(c, fiber.StatusNotFound, "Group not found", nil)
	}
	var req struct {
		Name          *string `json:"name"`
		IntervalHours *int    `json:"intervalHours"`
	}
	if err := c.BodyParser(&req); err != nil {
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
	return c.JSON(groupResponse(g, types.GroupRoleAdmin))
}

func (r *Router) createJoinRequest(c *fiber.Ctx) error {
	principal, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	groupID := c.Params("id")
	if m, _ := r.appDao.GetGroupMember(ctx, groupID, principal.UserID); m != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Already a member", nil)
	}
	if existing, _ := r.appDao.GetJoinRequestByUser(ctx, groupID, principal.UserID); existing != nil && existing.Status == models.JoinPending {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "pending"})
	}
	jr := &models.GroupJoinRequest{GroupID: groupID, UserID: principal.UserID, Status: models.JoinPending}
	if err := r.appDao.CreateJoinRequest(ctx, jr); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Request failed", err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": jr.ID, "status": jr.Status})
}

func (r *Router) listJoinRequests(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	list, err := r.appDao.ListPendingJoinRequests(ctx, c.Params("id"))
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "List failed", err)
	}
	return c.JSON(list)
}

func (r *Router) approveJoinRequest(c *fiber.Ctx) error {
	return r.resolveJoinRequest(c, models.JoinApproved, types.GroupRoleUser)
}

func (r *Router) rejectJoinRequest(c *fiber.Ctx) error {
	return r.resolveJoinRequest(c, models.JoinRejected, "")
}

func (r *Router) resolveJoinRequest(c *fiber.Ctx, status models.JoinRequestStatus, role types.GroupRole) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	jr, err := r.appDao.GetJoinRequest(ctx, c.Params("rid"))
	if err != nil || jr == nil {
		return errorResponse(c, fiber.StatusNotFound, "Request not found", nil)
	}
	if err := r.appDao.UpdateJoinRequestStatus(ctx, jr.ID, status); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Update failed", err)
	}
	if status == models.JoinApproved && role != "" {
		_ = r.appDao.CreateGroupMember(ctx, &models.GroupMember{GroupID: jr.GroupID, UserID: jr.UserID, GroupRole: role})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (r *Router) removeMember(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if err := r.appDao.DeleteGroupMember(ctx, c.Params("id"), c.Params("userId")); err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Remove failed", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (r *Router) leaderboard(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	rows, err := r.appDao.Leaderboard(ctx, c.Params("id"))
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Leaderboard failed", err)
	}
	return c.JSON(rows)
}

func groupResponse(g *models.Group, role types.GroupRole) fiber.Map {
	return fiber.Map{
		"id": g.ID, "name": g.Name, "intervalHours": g.IntervalHours,
		"timezone": g.Timezone, "groupRole": role,
	}
}
