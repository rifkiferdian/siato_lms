package coursecatalog

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
func u64(v string) uint64  { x, _ := strconv.ParseUint(v, 10, 64); return x }
func integer(v string) int { x, _ := strconv.Atoi(v); return x }
func (h *Handler) Page(c *gin.Context) {
	uid, _ := h.sm.UserID(c)
	rs, _ := h.r.Roles(uid)
	var u users.User
	_ = h.db.Select("id", "name", "email").First(&u, uid).Error
	d, e := h.s.Page(c.Query("q"), u64(c.Query("program")), u64(c.Query("curriculum")), integer(c.DefaultQuery("page", "1")))
	d.Success = c.Query("success")
	d.Error = c.Query("error")
	d.CanManage = h.r.HasPermission(uid, "academics.manage")
	if e != nil {
		d.Error = "Data mata kuliah belum dapat dimuat"
	}
	h.v.HTML(c, http.StatusOK, "courses/index.html", view.Data{Title: "Mata Kuliah", Roles: rs, User: u, Content: d})
}
func (h *Handler) SaveCourse(c *gin.Context) {
	in := CourseInput{CurriculumID: u64(c.PostForm("curriculum_id")), Code: c.PostForm("code"), Name: c.PostForm("name"), Credits: integer(c.PostForm("credits")), RecommendedSemester: integer(c.PostForm("recommended_semester")), CourseType: c.PostForm("course_type"), IsActive: c.PostForm("is_active") == "1", TheoryCredits: integer(c.PostForm("theory_credits")), PracticumCredits: integer(c.PostForm("practicum_credits")), PracticeCredits: integer(c.PostForm("practice_credits")), FieldCredits: integer(c.PostForm("field_credits"))}
	h.done(c, h.s.SaveCourse(u64(c.Param("id")), in), "Mata kuliah berhasil disimpan")
}
func (h *Handler) DeleteCourse(c *gin.Context) {
	h.done(c, h.s.DeleteCourse(u64(c.Param("id"))), "Mata kuliah berhasil dihapus")
}
func (h *Handler) SaveCurriculum(c *gin.Context) {
	in := CurriculumInput{StudyProgramID: u64(c.PostForm("study_program_id")), Code: c.PostForm("code"), Name: c.PostForm("name"), EffectiveYear: integer(c.PostForm("effective_year"))}
	h.done(c, h.s.SaveCurriculum(u64(c.Param("id")), in), "Kurikulum berhasil disimpan")
}
func (h *Handler) done(c *gin.Context, e error, ok string) {
	k, m := "success", ok
	if e != nil {
		k = "error"
		m = e.Error()
		if strings.Contains(strings.ToLower(m), "duplicate") {
			m = "Kode tersebut sudah digunakan pada kurikulum yang sama"
		}
	}
	c.Redirect(http.StatusSeeOther, "/admin/courses?"+k+"="+url.QueryEscape(m))
}
