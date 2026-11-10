package people

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
func (h *Handler) identity(c *gin.Context) ([]string, users.User) {
	id, _ := h.sm.UserID(c)
	roles, _ := h.roles.Roles(id)
	var u users.User
	_ = h.db.Select("id", "name", "email").First(&u, id).Error
	return roles, u
}
func (h *Handler) Students(c *gin.Context) {
	roles, u := h.identity(c)
	rows, programs, e := h.service.Students(c.Query("q"), c.Query("status"), c.Query("source"))
	p := PageData{Kind: "student", Search: c.Query("q"), Status: c.Query("status"), Source: c.Query("source"), Students: rows, Programs: programs, Success: c.Query("success"), Error: c.Query("error")}
	if e != nil {
		p.Error = "Data mahasiswa belum dapat dimuat."
	}
	h.views.HTML(c, http.StatusOK, "people/index.html", view.Data{Title: "Data Mahasiswa", Roles: roles, User: u, Content: p})
}
func (h *Handler) Lecturers(c *gin.Context) {
	roles, u := h.identity(c)
	rows, programs, e := h.service.Lecturers(c.Query("q"), c.Query("status"), c.Query("source"))
	p := PageData{Kind: "lecturer", Search: c.Query("q"), Status: c.Query("status"), Source: c.Query("source"), Lecturers: rows, Programs: programs, Success: c.Query("success"), Error: c.Query("error")}
	if e != nil {
		p.Error = "Data dosen belum dapat dimuat."
	}
	h.views.HTML(c, http.StatusOK, "people/index.html", view.Data{Title: "Data Dosen", Roles: roles, User: u, Content: p})
}
func (h *Handler) CreateStudent(c *gin.Context) {
	program, _ := strconv.ParseUint(c.PostForm("study_program_id"), 10, 64)
	year, _ := strconv.Atoi(c.PostForm("entry_year"))
	in := StudentInput{Name: c.PostForm("name"), Email: c.PostForm("email"), StudentNumber: c.PostForm("identity_number"), StudyProgramID: program, EntryYear: year, ExternalID: c.PostForm("external_id"), Password: c.PostForm("password")}
	if e := h.service.CreateStudent(in); e != nil {
		redirect(c, "/admin/students", "error", friendlyDBError(e))
		return
	}
	redirect(c, "/admin/students", "success", "Mahasiswa berhasil ditambahkan")
}
func (h *Handler) CreateLecturer(c *gin.Context) {
	program, _ := strconv.ParseUint(c.PostForm("study_program_id"), 10, 64)
	in := LecturerInput{Name: c.PostForm("name"), Email: c.PostForm("email"), EmployeeNumber: c.PostForm("identity_number"), StudyProgramID: program, ExternalID: c.PostForm("external_id"), Password: c.PostForm("password")}
	if e := h.service.CreateLecturer(in); e != nil {
		redirect(c, "/admin/lecturers", "error", friendlyDBError(e))
		return
	}
	redirect(c, "/admin/lecturers", "success", "Dosen berhasil ditambahkan")
}
func (h *Handler) ToggleStudent(c *gin.Context)  { h.toggle(c, "/admin/students") }
func (h *Handler) ToggleLecturer(c *gin.Context) { h.toggle(c, "/admin/lecturers") }
func (h *Handler) toggle(c *gin.Context, path string) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if e := h.service.Toggle(id); e != nil {
		redirect(c, path, "error", e.Error())
		return
	}
	redirect(c, path, "success", "Status akun berhasil diperbarui")
}
func redirect(c *gin.Context, path, key, message string) {
	c.Redirect(http.StatusSeeOther, path+"?"+key+"="+url.QueryEscape(message))
}
func friendlyDBError(e error) string {
	if strings.Contains(strings.ToLower(e.Error()), "duplicate") {
		return "Email, NIM, nomor pegawai, atau external ID sudah digunakan"
	}
	return e.Error()
}
