package academics

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
	service *Service
	roles   *rbac.Service
	sm      *middleware.SessionManager
	views   *view.Renderer
	db      *gorm.DB
}

func NewHandler(s *Service, r *rbac.Service, sm *middleware.SessionManager, v *view.Renderer, db *gorm.DB) *Handler {
	return &Handler{s, r, sm, v, db}
}
func (h *Handler) Page(c *gin.Context) {
	id, _ := h.sm.UserID(c)
	roles, _ := h.roles.Roles(id)
	var u users.User
	_ = h.db.Select("id", "name", "email").First(&u, id).Error
	d, e := h.service.Page()
	d.Tab = tab(c.DefaultQuery("tab", "years"))
	d.Success = c.Query("success")
	d.Error = c.Query("error")
	d.CanManage = h.roles.HasPermission(id, "academics.manage")
	if e != nil {
		d.Error = "Data master akademik belum dapat dimuat."
	}
	h.views.HTML(c, http.StatusOK, "academics/index.html", view.Data{Title: "Master Akademik", Roles: roles, User: u, Content: d})
}
func tab(v string) string {
	switch v {
	case "years", "semesters", "programs":
		return v
	}
	return "years"
}
func id(c *gin.Context) uint64 { v, _ := strconv.ParseUint(c.Param("id"), 10, 64); return v }
func uintForm(c *gin.Context, key string) uint64 {
	v, _ := strconv.ParseUint(c.PostForm(key), 10, 64)
	return v
}
func intForm(c *gin.Context, key string) int { v, _ := strconv.Atoi(c.PostForm(key)); return v }
func (h *Handler) SaveYear(c *gin.Context) {
	h.done(c, "years", h.service.SaveYear(id(c), AcademicYearInput{c.PostForm("name"), c.PostForm("starts_on"), c.PostForm("ends_on")}), "Tahun akademik berhasil disimpan")
}
func (h *Handler) DeleteYear(c *gin.Context) {
	h.done(c, "years", h.service.DeleteYear(id(c)), "Tahun akademik berhasil dihapus")
}
func (h *Handler) SaveSemester(c *gin.Context) {
	h.done(c, "semesters", h.service.SaveSemester(id(c), SemesterInput{uintForm(c, "academic_year_id"), c.PostForm("name"), c.PostForm("code"), c.PostForm("starts_on"), c.PostForm("ends_on")}), "Semester berhasil disimpan")
}
func (h *Handler) DeleteSemester(c *gin.Context) {
	h.done(c, "semesters", h.service.DeleteSemester(id(c)), "Semester berhasil dihapus")
}
func (h *Handler) ActivateSemester(c *gin.Context) {
	h.done(c, "semesters", h.service.ActivateSemester(id(c)), "Semester aktif berhasil ditetapkan")
}
func (h *Handler) SaveProgram(c *gin.Context) {
	h.done(c, "programs", h.service.SaveProgram(id(c), StudyProgramInput{c.PostForm("code"), c.PostForm("name"), c.PostForm("degree")}), "Program studi berhasil disimpan")
}
func (h *Handler) DeleteProgram(c *gin.Context) {
	h.done(c, "programs", h.service.DeleteProgram(id(c)), "Program studi berhasil dihapus")
}
func (h *Handler) SaveCourse(c *gin.Context) {
	h.done(c, "courses", h.service.SaveCourse(id(c), CourseInput{uintForm(c, "study_program_id"), c.PostForm("code"), c.PostForm("name"), intForm(c, "credits")}), "Mata kuliah berhasil disimpan")
}
func (h *Handler) DeleteCourse(c *gin.Context) {
	h.done(c, "courses", h.service.DeleteCourse(id(c)), "Mata kuliah berhasil dihapus")
}
func (h *Handler) done(c *gin.Context, tab string, e error, success string) {
	key, msg := "success", success
	if e != nil {
		key = "error"
		msg = friendly(e)
	}
	c.Redirect(http.StatusSeeOther, "/admin/academics?tab="+tab+"&"+key+"="+url.QueryEscape(msg))
}
func friendly(e error) string {
	v := strings.ToLower(e.Error())
	if strings.Contains(v, "duplicate") {
		return "Kode atau nama tersebut sudah digunakan"
	}
	if strings.Contains(v, "foreign key constraint") || strings.Contains(v, "cannot delete") {
		return "Data tidak dapat dihapus karena masih digunakan oleh data lain"
	}
	return e.Error()
}
