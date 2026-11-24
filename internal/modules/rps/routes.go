package rps

import (
	"github.com/gin-gonic/gin"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
)

func RegisterRoutes(r *gin.Engine, h *Handler, sm *middleware.SessionManager, roles *rbac.Service) {
	g := r.Group("")
	g.Use(middleware.RequireAuth(sm))
	g.GET("/admin/courses/:id", middleware.RequirePermission(roles.HasPermission, "academics.view"), h.Admin)
	g.GET("/lecturer/courses/:id", middleware.RequirePermission(roles.HasPermission, "classes.view"), h.Lecturer)
	g.POST("/courses/:id/rps", h.Save)
	g.POST("/courses/:id/rps/:rpsID", h.Save)
	g.POST("/courses/:id/rps/:rpsID/activate", h.Activate)
	g.POST("/courses/:id/rps/:rpsID/cpmk", h.SaveCPMK)
	g.POST("/courses/:id/rps/:rpsID/cpmk/:cpmkID", h.SaveCPMK)
	g.POST("/courses/:id/rps/:rpsID/cpmk/:cpmkID/delete", h.DeleteCPMK)
}
