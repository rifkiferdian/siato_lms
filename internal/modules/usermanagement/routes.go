package usermanagement

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
)

func RegisterRoutes(r *gin.Engine, h *Handler, sm *middleware.SessionManager, roles *rbac.Service) {
	g := r.Group("/admin")
	g.Use(middleware.RequireAuth(sm), requireSuperAdmin(roles))
	g.GET("/users", middleware.RequirePermission(roles.HasPermission, "users.view"), h.Users)
	g.POST("/users", middleware.RequirePermission(roles.HasPermission, "users.manage"), h.Create)
	g.POST("/users/:id/toggle", middleware.RequirePermission(roles.HasPermission, "users.manage"), h.Toggle)
	g.POST("/users/:id/roles", middleware.RequirePermission(roles.HasPermission, "users.manage"), h.SetRoles)
	g.GET("/roles", middleware.RequirePermission(roles.HasPermission, "roles.view"), h.Roles)
	g.POST("/roles/:id/permissions", middleware.RequirePermission(roles.HasPermission, "roles.manage"), h.SetPermissions)
}

func requireSuperAdmin(roles *rbac.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		names, err := roles.Roles(value.(uint64))
		if err == nil {
			for _, name := range names {
				if name == "super_admin" {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatus(http.StatusForbidden)
	}
}
