package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
	"siato_lms/internal/shared/view"
	"strconv"
	"strings"
)

type Handler struct {
	service              *Service
	rbac                 *rbac.Service
	sm                   *middleware.SessionManager
	views                *view.Renderer
	baseURL, environment string
}

func NewHandler(s *Service, r *rbac.Service, sm *middleware.SessionManager, v *view.Renderer, baseURL, environment string) *Handler {
	return &Handler{s, r, sm, v, baseURL, environment}
}

type AccountData struct {
	Sessions       []SessionRow
	Success, Error string
}
type ResetData struct{ Token, DevelopmentURL, Success, Error string }

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
	if err = h.sm.Start(c, u.ID); err != nil {
		c.Status(500)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard")
}
func (h *Handler) ForgotPage(c *gin.Context) {
	h.views.HTML(c, http.StatusOK, "auth/forgot.html", view.Data{Title: "Lupa Kata Sandi"})
}
func (h *Handler) Forgot(c *gin.Context) {
	token, e := h.service.RequestReset(c.PostForm("email"))
	d := ResetData{Success: "Jika email terdaftar, tautan reset telah dibuat."}
	if e != nil {
		d.Error = "Permintaan reset belum dapat diproses."
	}
	if token != "" && h.environment != "production" {
		d.DevelopmentURL = h.baseURL + "/reset-password?token=" + token
	}
	h.views.HTML(c, http.StatusOK, "auth/forgot.html", view.Data{Title: "Lupa Kata Sandi", Content: d})
}
func (h *Handler) ResetPage(c *gin.Context) {
	h.views.HTML(c, http.StatusOK, "auth/reset.html", view.Data{Title: "Reset Kata Sandi", Content: ResetData{Token: c.Query("token")}})
}
func (h *Handler) Reset(c *gin.Context) {
	token := c.PostForm("token")
	if e := h.service.ResetPassword(token, c.PostForm("password"), c.PostForm("password_confirmation")); e != nil {
		h.views.HTML(c, http.StatusUnprocessableEntity, "auth/reset.html", view.Data{Title: "Reset Kata Sandi", Content: ResetData{Token: token, Error: e.Error()}})
		return
	}
	c.Redirect(http.StatusSeeOther, "/login?reset=success")
}
func (h *Handler) Account(c *gin.Context) {
	id, _ := h.sm.UserID(c)
	u, e := h.service.Profile(id)
	if e != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	roles, _ := h.rbac.Roles(id)
	sessions, _ := h.service.Sessions(id, h.sm.CurrentTokenHash(c))
	d := AccountData{Sessions: sessions, Success: c.Query("success"), Error: c.Query("error")}
	h.views.HTML(c, http.StatusOK, "account/index.html", view.Data{Title: "Profil & Keamanan", User: u, Roles: roles, Content: d})
}
func (h *Handler) UpdateProfile(c *gin.Context) {
	id, _ := h.sm.UserID(c)
	if e := h.service.UpdateProfile(id, c.PostForm("name"), c.PostForm("email")); e != nil {
		c.Redirect(http.StatusSeeOther, "/account?error="+url.QueryEscape(e.Error()))
		return
	}
	c.Redirect(http.StatusSeeOther, "/account?success="+url.QueryEscape("Profil berhasil diperbarui"))
}
func (h *Handler) ChangePassword(c *gin.Context) {
	id, _ := h.sm.UserID(c)
	if e := h.service.ChangePassword(id, c.PostForm("current_password"), c.PostForm("password"), c.PostForm("password_confirmation")); e != nil {
		c.Redirect(http.StatusSeeOther, "/account?error="+url.QueryEscape(e.Error()))
		return
	}
	_ = h.sm.Destroy(c)
	c.Redirect(http.StatusSeeOther, "/login?password=changed")
}
func (h *Handler) RevokeSession(c *gin.Context) {
	id, _ := h.sm.UserID(c)
	sessionID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	_ = h.service.RevokeSession(id, sessionID)
	c.Redirect(http.StatusSeeOther, "/account?success="+url.QueryEscape("Sesi berhasil dikeluarkan"))
}
func (h *Handler) LogoutAll(c *gin.Context) {
	id, _ := h.sm.UserID(c)
	_ = h.sm.DestroyAll(c, id)
	c.Redirect(http.StatusSeeOther, "/login")
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
