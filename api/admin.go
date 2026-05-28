package api

import (
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/models"
	"github.com/gofiber/fiber/v2"
)

func (r *Router) initAdminRoutes() {
	a := r.apiGroup("admin")
	a.Get("/users", r.requireAuth, r.requireSiteAdmin, r.adminListUsers)
	a.Get("/groups", r.requireAuth, r.requireSiteAdmin, r.adminListGroups)
}

func (r *Router) adminListUsers(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	users, err := r.appDao.ListUsers(ctx, dao.NewOptions())
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "List failed", err)
	}
	var out []fiber.Map
	for _, u := range users {
		out = append(out, fiber.Map{
			"id": u.ID, "username": u.Username, "displayName": u.DisplayName, "siteRole": u.SiteRole,
		})
	}
	return c.JSON(out)
}

func (r *Router) adminListGroups(c *fiber.Ctx) error {
	_, ctx, err := principalCtx(c)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Unauthorized", nil)
	}
	groups, err := r.appDao.ListAllGroups(ctx)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "List failed", err)
	}
	var out []fiber.Map
	for _, g := range groups {
		out = append(out, fiber.Map{"id": g.ID, "name": g.Name})
	}
	return c.JSON(out)
}

var _ = models.Group{}
