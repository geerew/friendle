package api

import (
	"context"
	"time"

	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/types"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) requireAuth(c *fiber.Ctx) error {
	if _, _, err := principalCtx(c); err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	return c.Next()
}

func (r *Router) requireSiteAdmin(c *fiber.Ctx) error {
	p, _, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if p.SiteRole != types.SiteRoleAdmin {
		return errorResponse(c, fiber.StatusForbidden, "Site admin required", nil)
	}
	return c.Next()
}

func (r *Router) requireGroupAdmin(c *fiber.Ctx) error {
	p, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	m, err := r.appDao.GetGroupMember(ctx, c.Params("id"), p.UserID)
	if err != nil || m == nil || m.GroupRole != types.GroupRoleAdmin {
		return errorResponse(c, fiber.StatusForbidden, "Group admin required", nil)
	}
	return c.Next()
}

func (r *Router) requireGroupMember(c *fiber.Ctx) error {
	p, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	if _, err := r.membership(ctx, c.Params("id"), p.UserID); err != nil {
		return errorResponse(c, fiber.StatusForbidden, "Not a member", nil)
	}
	return c.Next()
}

func (r *Router) membership(ctx context.Context, groupID, userID string) (*models.GroupMember, error) {
	m, err := r.appDao.GetGroupMember(ctx, groupID, userID)
	if err != nil || m == nil {
		return nil, err
	}
	return m, nil
}

func (r *Router) roundSummary(ctx context.Context, g *models.Group, userID string) fiber.Map {
	roundDate := time.Now().Format("2006-01-02")
	round, _ := r.appDao.GetCurrentRound(ctx, g.ID, roundDate)
	if round == nil {
		return fiber.Map{"status": "none"}
	}
	m, _ := r.appDao.GetGroupMember(ctx, g.ID, userID)
	yourRole := "guesser"
	if round.PickerUserID == userID {
		yourRole = "picker"
	}
	out := fiber.Map{
		"status": round.Status, "yourRole": yourRole,
	}
	if round.Status == models.RoundCompleted || round.Status == models.RoundSkipped {
		out["canReveal"] = true
	}
	if m != nil {
		out["groupRole"] = m.GroupRole
	}
	return out
}
