package academics

import (
	"github.com/gin-gonic/gin"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
)

func RegisterRoutes(r *gin.Engine, h *Handler, sm *middleware.SessionManager, roles *rbac.Service) {
	g := r.Group("/admin/academics")
	g.Use(middleware.RequireAuth(sm))
	g.GET("", middleware.RequirePermission(roles.HasPermission, "academics.view"), h.Page)
	m := g.Group("")
	m.Use(middleware.RequirePermission(roles.HasPermission, "academics.manage"))
	m.POST("/years", h.SaveYear)
	m.POST("/years/:id", h.SaveYear)
	m.POST("/years/:id/delete", h.DeleteYear)
	m.POST("/semesters", h.SaveSemester)
	m.POST("/semesters/:id", h.SaveSemester)
	m.POST("/semesters/:id/delete", h.DeleteSemester)
	m.POST("/semesters/:id/activate", h.ActivateSemester)
	m.POST("/programs", h.SaveProgram)
	m.POST("/programs/:id", h.SaveProgram)
	m.POST("/programs/:id/delete", h.DeleteProgram)
	m.POST("/courses", h.SaveCourse)
	m.POST("/courses/:id", h.SaveCourse)
	m.POST("/courses/:id/delete", h.DeleteCourse)
}
