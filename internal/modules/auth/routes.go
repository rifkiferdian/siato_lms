package auth

import (
	"github.com/gin-gonic/gin"
	"siato_lms/internal/middleware"
)

func RegisterRoutes(r *gin.Engine, h *Handler, sm *middleware.SessionManager) {
	r.GET("/login", h.LoginPage)
	r.POST("/login", middleware.LoginLimiter(), h.Login)
	secured := r.Group("")
	secured.Use(middleware.RequireAuth(sm))
	secured.GET("/dashboard", h.Dashboard)
	secured.POST("/logout", h.Logout)
}
