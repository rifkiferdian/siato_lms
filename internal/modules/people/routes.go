package people

import (
	"github.com/gin-gonic/gin"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
)

func RegisterRoutes(r *gin.Engine, h *Handler, sm *middleware.SessionManager, roles *rbac.Service) {
	g := r.Group("/admin")
	g.Use(middleware.RequireAuth(sm))
	g.GET("/students", middleware.RequirePermission(roles.HasPermission, "students.view"), h.Students)
	g.POST("/students", middleware.RequirePermission(roles.HasPermission, "students.manage"), h.CreateStudent)
	g.POST("/students/:id/toggle", middleware.RequirePermission(roles.HasPermission, "students.manage"), h.ToggleStudent)
	g.GET("/lecturers", middleware.RequirePermission(roles.HasPermission, "lecturers.view"), h.Lecturers)
	g.POST("/lecturers", middleware.RequirePermission(roles.HasPermission, "lecturers.manage"), h.CreateLecturer)
	g.POST("/lecturers/:id/toggle", middleware.RequirePermission(roles.HasPermission, "lecturers.manage"), h.ToggleLecturer)
}
