package auth

import (
	"github.com/gin-gonic/gin"
	"siato_lms/internal/middleware"
)

func RegisterRoutes(r *gin.Engine, h *Handler, sm *middleware.SessionManager) {
	r.GET("/login", h.LoginPage)
	r.POST("/login", middleware.LoginLimiter(), h.Login)
	r.GET("/forgot-password", h.ForgotPage)
	r.POST("/forgot-password", middleware.LoginLimiter(), h.Forgot)
	r.GET("/reset-password", h.ResetPage)
	r.POST("/reset-password", h.Reset)
	secured := r.Group("")
	secured.Use(middleware.RequireAuth(sm))
	secured.GET("/dashboard", h.Dashboard)
	secured.POST("/logout", h.Logout)
	secured.GET("/account", h.Account)
	secured.POST("/account/profile", h.UpdateProfile)
	secured.POST("/account/password", h.ChangePassword)
	secured.POST("/account/sessions/:id/revoke", h.RevokeSession)
	secured.POST("/account/sessions/logout-all", h.LogoutAll)
}
