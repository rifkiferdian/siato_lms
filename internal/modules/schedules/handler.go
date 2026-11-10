package schedules

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
	"siato_lms/internal/modules/users"
	"siato_lms/internal/shared/view"
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
func (h *Handler) page(c *gin.Context, portal, title string) {
	id, _ := h.sm.UserID(c)
	day, _ := strconv.Atoi(c.Query("day"))
	rows, classes, err := h.service.Page(portal, id, day)
	roleNames, _ := h.roles.Roles(id)
	var user users.User
	_ = h.db.Select("id", "name", "email").First(&user, id).Error
	p := PageData{Portal: portal, Schedules: rows, Classes: classes, SelectedDay: day, Days: Days(), CanManage: portal == "admin", Success: c.Query("success"), Error: c.Query("error")}
	if err != nil {
		p.Error = "Data jadwal belum dapat dimuat."
	}
	h.views.HTML(c, http.StatusOK, "schedules/index.html", view.Data{Title: title, Roles: roleNames, User: user, Content: p})
}
func (h *Handler) Admin(c *gin.Context)    { h.page(c, "admin", "Kelola Jadwal Kuliah") }
func (h *Handler) Lecturer(c *gin.Context) { h.page(c, "lecturer", "Jadwal Mengajar") }
func (h *Handler) Student(c *gin.Context)  { h.page(c, "student", "Jadwal Kuliah") }
func (h *Handler) Create(c *gin.Context) {
	classID, _ := strconv.ParseUint(c.PostForm("class_id"), 10, 64)
	day, _ := strconv.Atoi(c.PostForm("day_of_week"))
	in := CreateInput{ClassID: classID, DayOfWeek: day, StartTime: c.PostForm("start_time"), EndTime: c.PostForm("end_time"), Room: c.PostForm("room"), DeliveryMode: c.PostForm("delivery_mode"), MeetingURL: c.PostForm("meeting_url"), EffectiveFrom: c.PostForm("effective_from"), EffectiveUntil: c.PostForm("effective_until")}
	if err := h.service.Create(in); err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/schedules?error="+urlQuery(err.Error()))
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/schedules?success=Jadwal berhasil ditambahkan")
}
func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.service.Delete(id); err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/schedules?error="+urlQuery(err.Error()))
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/schedules?success=Jadwal berhasil dihapus")
}
func urlQuery(v string) string { return url.QueryEscape(v) }
