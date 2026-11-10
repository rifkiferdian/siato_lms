package schedules

import (
	"github.com/gin-gonic/gin"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
)

func RegisterRoutes(r *gin.Engine, h *Handler, sm *middleware.SessionManager, roles *rbac.Service) {
	g := r.Group("")
	g.Use(middleware.RequireAuth(sm))
	g.GET("/admin/schedules", middleware.RequirePermission(roles.HasPermission, "schedules.manage"), h.Admin)
	g.POST("/admin/schedules", middleware.RequirePermission(roles.HasPermission, "schedules.manage"), h.Create)
	g.POST("/admin/schedules/:id/delete", middleware.RequirePermission(roles.HasPermission, "schedules.manage"), h.Delete)
	g.GET("/lecturer/schedule", middleware.RequirePermission(roles.HasPermission, "schedules.view"), h.Lecturer)
	g.GET("/student/schedule", middleware.RequirePermission(roles.HasPermission, "schedules.view"), h.Student)
}
