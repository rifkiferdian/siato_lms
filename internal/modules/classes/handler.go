package classes

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
func (h *Handler) identity(c *gin.Context) (uint64, []string, users.User) {
	id, _ := h.sm.UserID(c)
	roles, _ := h.roles.Roles(id)
	var u users.User
	_ = h.db.Select("id", "name", "email").First(&u, id).Error
	return id, roles, u
}
func (h *Handler) page(c *gin.Context, portal, title string) {
	id, roles, u := h.identity(c)
	search, status := c.Query("q"), c.Query("status")
	rows, courses, semesters, err := h.service.List(portal, id, search, status)
	p := PageData{Portal: portal, Search: search, Status: status, Classes: rows, Courses: courses, Semesters: semesters, CanManage: portal == "admin", Success: c.Query("success"), Error: c.Query("error")}
	if err != nil {
		p.Error = "Data kelas belum dapat dimuat."
	}
	h.views.HTML(c, http.StatusOK, "classes/index.html", view.Data{Title: title, Roles: roles, User: u, Content: p})
}
func (h *Handler) detail(c *gin.Context, portal, title string) {
	id, roles, u := h.identity(c)
	classID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	d, err := h.service.Detail(classID, id, portal)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	data := DetailData{Portal: portal, Class: d, CanManage: portal == "admin", Success: c.Query("success"), Error: c.Query("error")}
	data.CanManageParticipants = portal == "admin" && canManageClassStudents(roles)
	if data.CanManageParticipants {
		data.Students, data.AvailableStudents, err = h.service.Students(classID)
		if err != nil {
			data.Error = "Data peserta belum dapat dimuat."
		}
	}
	h.views.HTML(c, http.StatusOK, "classes/detail.html", view.Data{Title: title, Roles: roles, User: u, Content: data})
}
func (h *Handler) Admin(c *gin.Context)          { h.page(c, "admin", "Kelola Kelas") }
func (h *Handler) Lecturer(c *gin.Context)       { h.page(c, "lecturer", "Kelas Saya") }
func (h *Handler) Student(c *gin.Context)        { h.page(c, "student", "Kelas Saya") }
func (h *Handler) AdminDetail(c *gin.Context)    { h.detail(c, "admin", "Detail Kelas") }
func (h *Handler) LecturerDetail(c *gin.Context) { h.detail(c, "lecturer", "Detail Kelas") }
func (h *Handler) StudentDetail(c *gin.Context)  { h.detail(c, "student", "Detail Kelas") }
func (h *Handler) Create(c *gin.Context) {
	course, _ := strconv.ParseUint(c.PostForm("course_id"), 10, 64)
	semester, _ := strconv.ParseUint(c.PostForm("semester_id"), 10, 64)
	capacity, _ := strconv.Atoi(c.PostForm("capacity"))
	meetings, _ := strconv.Atoi(c.PostForm("max_meetings"))
	in := CreateInput{CourseID: course, SemesterID: semester, Name: c.PostForm("name"), Description: c.PostForm("description"), Capacity: capacity, MaxMeetings: meetings, Status: c.PostForm("status")}
	if err := h.service.Create(in); err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/classes?error="+url.QueryEscape(err.Error()))
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/classes?success="+url.QueryEscape("Kelas berhasil dibuat"))
}
func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/classes?error="+url.QueryEscape(err.Error()))
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/classes?success="+url.QueryEscape("Kelas berhasil dihapus"))
}
func (h *Handler) AddStudents(c *gin.Context) {
	if !h.canManageClassStudents(c) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	classID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	values := c.PostFormArray("student_ids")
	studentIDs := make([]uint64, 0, len(values))
	for _, value := range values {
		id, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			h.redirectDetail(c, classID, "error", "Pilihan mahasiswa tidak valid")
			return
		}
		studentIDs = append(studentIDs, id)
	}
	if err := h.service.AddStudents(classID, studentIDs); err != nil {
		h.redirectDetail(c, classID, "error", err.Error())
		return
	}
	h.redirectDetail(c, classID, "success", "Mahasiswa berhasil ditambahkan ke kelas")
}
func (h *Handler) RemoveStudent(c *gin.Context) {
	if !h.canManageClassStudents(c) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	classID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	studentID, _ := strconv.ParseUint(c.Param("studentID"), 10, 64)
	if err := h.service.RemoveStudent(classID, studentID); err != nil {
		h.redirectDetail(c, classID, "error", err.Error())
		return
	}
	h.redirectDetail(c, classID, "success", "Mahasiswa berhasil dikeluarkan dari kelas")
}
func (h *Handler) canManageClassStudents(c *gin.Context) bool {
	id, _ := h.sm.UserID(c)
	roles, _ := h.roles.Roles(id)
	return canManageClassStudents(roles)
}
func (h *Handler) redirectDetail(c *gin.Context, classID uint64, key, message string) {
	c.Redirect(http.StatusSeeOther, "/admin/classes/"+strconv.FormatUint(classID, 10)+"?"+key+"="+url.QueryEscape(message))
}
func hasRole(roles []string, wanted string) bool {
	for _, role := range roles {
		if role == wanted {
			return true
		}
	}
	return false
}

func canManageClassStudents(roles []string) bool {
	return hasRole(roles, "academic_admin") || hasRole(roles, "super_admin")
}
