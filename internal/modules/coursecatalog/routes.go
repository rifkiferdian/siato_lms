package coursecatalog

import (
	"github.com/gin-gonic/gin"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
)

func RegisterRoutes(r *gin.Engine, h *Handler, sm *middleware.SessionManager, roles *rbac.Service) {
	g := r.Group("/admin/courses")
	g.Use(middleware.RequireAuth(sm))
	g.GET("", middleware.RequirePermission(roles.HasPermission, "academics.view"), h.Page)
	g.POST("", middleware.RequirePermission(roles.HasPermission, "academics.manage"), h.SaveCourse)
	g.POST("/:id", middleware.RequirePermission(roles.HasPermission, "academics.manage"), h.SaveCourse)
	g.POST("/:id/delete", middleware.RequirePermission(roles.HasPermission, "academics.manage"), h.DeleteCourse)
	g.POST("/curricula", middleware.RequirePermission(roles.HasPermission, "academics.manage"), h.SaveCurriculum)
	g.POST("/curricula/:id", middleware.RequirePermission(roles.HasPermission, "academics.manage"), h.SaveCurriculum)
}
