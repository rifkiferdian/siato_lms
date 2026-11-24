package rps

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
	"siato_lms/internal/modules/users"
	"siato_lms/internal/shared/view"
	"strconv"
	"strings"
)

type Handler struct {
	s  *Service
	r  *rbac.Service
	sm *middleware.SessionManager
	v  *view.Renderer
	db *gorm.DB
}

func NewHandler(s *Service, r *rbac.Service, sm *middleware.SessionManager, v *view.Renderer, db *gorm.DB) *Handler {
	return &Handler{s, r, sm, v, db}
}
func id(c *gin.Context, key string) uint64 { v, _ := strconv.ParseUint(c.Param(key), 10, 64); return v }
func queryID(c *gin.Context) uint64        { v, _ := strconv.ParseUint(c.Query("rps"), 10, 64); return v }
func number(v string) int                  { n, _ := strconv.Atoi(v); return n }
func (h *Handler) page(c *gin.Context, portal string) {
	uid, _ := h.sm.UserID(c)
	roles, _ := h.r.Roles(uid)
	var u users.User
	_ = h.db.Select("id", "name", "email").First(&u, uid).Error
	d, e := h.s.Page(id(c, "id"), queryID(c), uid, portal)
	if e != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	d.Success = c.Query("success")
	d.Error = c.Query("error")
	if c.Query("new") == "1" && d.CanManage {
		d.RPS = RPS{}
		d.CPMKs = nil
	}
	h.v.HTML(c, http.StatusOK, "rps/index.html", view.Data{Title: "Detail Mata Kuliah", Roles: roles, User: u, Content: d})
}
func (h *Handler) Admin(c *gin.Context)    { h.page(c, "admin") }
func (h *Handler) Lecturer(c *gin.Context) { h.page(c, "lecturer") }
func (h *Handler) authorized(c *gin.Context) (uint64, uint64, bool) {
	uid, _ := h.sm.UserID(c)
	courseID := id(c, "id")
	d, e := h.s.Page(courseID, 0, uid, "lecturer")
	roles, _ := h.r.Roles(uid)
	admin := false
	for _, role := range roles {
		if role == "super_admin" || role == "academic_admin" {
			admin = true
		}
	}
	return uid, courseID, (e == nil && d.CanManage) || admin
}
func (h *Handler) Save(c *gin.Context) {
	uid, courseID, ok := h.authorized(c)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	in := RPSInput{c.PostForm("version"), c.PostForm("description"), c.PostForm("learning_methods"), c.PostForm("main_references"), c.PostForm("supporting_references")}
	h.done(c, courseID, h.s.Save(courseID, id(c, "rpsID"), uid, in), "RPS berhasil disimpan")
}
func (h *Handler) Activate(c *gin.Context) {
	uid, courseID, ok := h.authorized(c)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	h.done(c, courseID, h.s.Activate(courseID, id(c, "rpsID"), uid), "RPS berhasil diaktifkan")
}
func (h *Handler) SaveCPMK(c *gin.Context) {
	_, courseID, ok := h.authorized(c)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	in := CPMKInput{c.PostForm("code"), c.PostForm("title"), c.PostForm("description"), number(c.PostForm("sort_order"))}
	h.done(c, courseID, h.s.SaveCPMK(courseID, id(c, "rpsID"), id(c, "cpmkID"), in), "CPMK berhasil disimpan")
}
func (h *Handler) DeleteCPMK(c *gin.Context) {
	_, courseID, ok := h.authorized(c)
	if !ok {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	h.done(c, courseID, h.s.DeleteCPMK(courseID, id(c, "rpsID"), id(c, "cpmkID")), "CPMK berhasil dihapus")
}
func (h *Handler) done(c *gin.Context, courseID uint64, e error, ok string) {
	key, msg := "success", ok
	if e != nil {
		key = "error"
		msg = e.Error()
		if strings.Contains(strings.ToLower(msg), "duplicate") {
			msg = "Versi RPS atau kode CPMK sudah digunakan"
		}
	}
	c.Redirect(http.StatusSeeOther, "/admin/courses/"+strconv.FormatUint(courseID, 10)+"?"+key+"="+url.QueryEscape(msg))
}
