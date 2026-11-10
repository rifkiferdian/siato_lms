package classes

import (
	"github.com/gin-gonic/gin"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
)

func RegisterRoutes(r *gin.Engine, h *Handler, sm *middleware.SessionManager, roles *rbac.Service) {
	g := r.Group("")
	g.Use(middleware.RequireAuth(sm))
	g.GET("/admin/classes", middleware.RequirePermission(roles.HasPermission, "classes.manage"), h.Admin)
	g.POST("/admin/classes", middleware.RequirePermission(roles.HasPermission, "classes.manage"), h.Create)
	g.GET("/admin/classes/:id", middleware.RequirePermission(roles.HasPermission, "classes.manage"), h.AdminDetail)
	g.POST("/admin/classes/:id/delete", middleware.RequirePermission(roles.HasPermission, "classes.manage"), h.Delete)
	g.GET("/lecturer/classes", middleware.RequirePermission(roles.HasPermission, "classes.view"), h.Lecturer)
	g.GET("/lecturer/classes/:id", middleware.RequirePermission(roles.HasPermission, "classes.view"), h.LecturerDetail)
	g.GET("/student/classes", middleware.RequirePermission(roles.HasPermission, "classes.view"), h.Student)
	g.GET("/student/classes/:id", middleware.RequirePermission(roles.HasPermission, "classes.view"), h.StudentDetail)
}
