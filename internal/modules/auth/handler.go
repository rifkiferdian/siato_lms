package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
	"siato_lms/internal/shared/view"
	"strings"
)

type Handler struct {
	service *Service
	rbac    *rbac.Service
	sm      *middleware.SessionManager
	views   *view.Renderer
}

func NewHandler(s *Service, r *rbac.Service, sm *middleware.SessionManager, v *view.Renderer) *Handler {
	return &Handler{s, r, sm, v}
}
func (h *Handler) LoginPage(c *gin.Context) {
	if _, ok := h.sm.UserID(c); ok {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}
	h.views.HTML(c, http.StatusOK, "auth/login.html", view.Data{Title: "Masuk"})
}
func (h *Handler) Login(c *gin.Context) {
	u, err := h.service.Login(c.PostForm("email"), c.PostForm("password"))
	if err != nil {
		h.views.HTML(c, http.StatusUnprocessableEntity, "auth/login.html", view.Data{Title: "Masuk", Error: err.Error()})
		return
	}
	s, _ := h.sm.Get(c)
	s.Values["user_id"] = u.ID
	if err = s.Save(c.Request, c.Writer); err != nil {
		c.Status(500)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard")
}
func (h *Handler) Logout(c *gin.Context) {
	_ = h.sm.Destroy(c)
	c.Redirect(http.StatusSeeOther, "/login")
}
func (h *Handler) Dashboard(c *gin.Context) {
	id, _ := h.sm.UserID(c)
	roles, _ := h.rbac.Roles(id)
	destination := "/student/dashboard"
	if contains(roles, "super_admin") || contains(roles, "academic_admin") {
		destination = "/admin/dashboard"
	} else if contains(roles, "lecturer") {
		destination = "/lecturer/dashboard"
	}
	c.Redirect(http.StatusSeeOther, destination)
}
func contains(items []string, v string) bool {
	for _, x := range items {
		if strings.EqualFold(x, v) {
			return true
		}
	}
	return false
}
