package dashboard

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"siato_lms/internal/middleware"
	"siato_lms/internal/modules/rbac"
	"siato_lms/internal/modules/users"
	"siato_lms/internal/shared/view"
)

type Handler struct {
	roles *rbac.Service
	sm    *middleware.SessionManager
	views *view.Renderer
	db    *gorm.DB
}

func NewHandler(r *rbac.Service, sm *middleware.SessionManager, v *view.Renderer, db *gorm.DB) *Handler {
	return &Handler{roles: r, sm: sm, views: v, db: db}
}
func (h *Handler) show(c *gin.Context, portal, title string) {
	id, _ := h.sm.UserID(c)
	roles, _ := h.roles.Roles(id)
	var user users.User
	_ = h.db.Select("id", "name", "email").First(&user, id).Error
	h.views.HTML(c, http.StatusOK, "dashboard/index.html", view.Data{Title: title, Roles: roles, User: user, Content: map[string]any{"Portal": portal}})
}
func (h *Handler) Admin(c *gin.Context)    { h.show(c, "admin", "Dashboard Admin") }
func (h *Handler) Lecturer(c *gin.Context) { h.show(c, "lecturer", "Dashboard Dosen") }
func (h *Handler) Student(c *gin.Context)  { h.show(c, "student", "Dashboard Mahasiswa") }
func RegisterRoutes(r *gin.Engine, h *Handler, sm *middleware.SessionManager, roles *rbac.Service) {
	g := r.Group("")
	g.Use(middleware.RequireAuth(sm))
	g.GET("/admin/dashboard", middleware.RequirePermission(roles.HasPermission, "admin.dashboard"), h.Admin)
	g.GET("/lecturer/dashboard", middleware.RequirePermission(roles.HasPermission, "lecturer.dashboard"), h.Lecturer)
	g.GET("/student/dashboard", middleware.RequirePermission(roles.HasPermission, "student.dashboard"), h.Student)
}
