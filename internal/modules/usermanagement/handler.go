package usermanagement

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

func NewHandler(service *Service, roles *rbac.Service, sm *middleware.SessionManager, views *view.Renderer, db *gorm.DB) *Handler {
	return &Handler{service, roles, sm, views, db}
}
func (h *Handler) identity(c *gin.Context) (uint64, []string, users.User) {
	id, _ := h.sm.UserID(c)
	roleNames, _ := h.roles.Roles(id)
	var user users.User
	_ = h.db.Select("id", "name", "email").First(&user, id).Error
	return id, roleNames, user
}
func (h *Handler) Users(c *gin.Context) {
	_, roles, user := h.identity(c)
	search, status, role := c.Query("q"), c.Query("status"), c.Query("role")
	rows, roleOptions, adminRoles, err := h.service.Users(search, status, role)
	data := UsersPage{Search: search, Status: status, Role: role, Users: rows, Roles: roleOptions, AdminRoles: adminRoles, Success: c.Query("success"), Error: c.Query("error")}
	if err != nil {
		data.Error = "Data pengguna belum dapat dimuat."
	}
	h.views.HTML(c, http.StatusOK, "users/index.html", view.Data{Title: "User Management", Roles: roles, User: user, Content: data})
}
func (h *Handler) Create(c *gin.Context) {
	roleID, _ := strconv.ParseUint(c.PostForm("role_id"), 10, 64)
	h.redirect(c, "/admin/users", h.service.Create(c.PostForm("name"), c.PostForm("email"), c.PostForm("password"), roleID), "Akun administratif berhasil dibuat")
}
func (h *Handler) Toggle(c *gin.Context) {
	actorID, _ := h.sm.UserID(c)
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h.redirect(c, "/admin/users", h.service.Toggle(actorID, userID), "Status pengguna berhasil diperbarui")
}
func (h *Handler) SetRoles(c *gin.Context) {
	actorID, _ := h.sm.UserID(c)
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	roleIDs, err := parseIDs(c.PostFormArray("role_ids"))
	if err == nil {
		err = h.service.SetRoles(actorID, userID, roleIDs)
	}
	h.redirect(c, "/admin/users", err, "Role pengguna berhasil diperbarui")
}
func (h *Handler) Roles(c *gin.Context) {
	_, roles, user := h.identity(c)
	selected, _ := strconv.ParseUint(c.Query("role"), 10, 64)
	rows, permissions, current, err := h.service.Roles(selected)
	data := RolesPage{Roles: rows, Permissions: permissions, SelectedRole: current, Success: c.Query("success"), Error: c.Query("error")}
	if err != nil {
		data.Error = err.Error()
	}
	h.views.HTML(c, http.StatusOK, "users/roles.html", view.Data{Title: "Role & Permission", Roles: roles, User: user, Content: data})
}
func (h *Handler) SetPermissions(c *gin.Context) {
	roleID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	permissionIDs, err := parseIDs(c.PostFormArray("permission_ids"))
	if err == nil {
		err = h.service.SetPermissions(roleID, permissionIDs)
	}
	h.redirect(c, "/admin/roles?role="+strconv.FormatUint(roleID, 10), err, "Permission role berhasil diperbarui")
}
func parseIDs(values []string) ([]uint64, error) {
	ids := make([]uint64, 0, len(values))
	for _, value := range values {
		id, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
func (h *Handler) redirect(c *gin.Context, path string, err error, success string) {
	separator := "?"
	if strings.Contains(path, "?") {
		separator = "&"
	}
	if err != nil {
		c.Redirect(http.StatusSeeOther, path+separator+"error="+url.QueryEscape(err.Error()))
		return
	}
	c.Redirect(http.StatusSeeOther, path+separator+"success="+url.QueryEscape(success))
}
