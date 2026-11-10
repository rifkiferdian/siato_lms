package view

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"siato_lms/internal/middleware"
)

type Renderer struct {
	sm      *middleware.SessionManager
	root    string
	appName string
}
type Data struct {
	Title, AppName, CSRF, Error, Success string
	User                                 any
	Roles                                []string
	Content                              any
}

func New(sm *middleware.SessionManager, root, appName string) *Renderer {
	return &Renderer{sm: sm, root: root, appName: appName}
}
func (r *Renderer) HTML(c *gin.Context, status int, page string, d Data) {
	d.AppName = r.appName
	d.CSRF, _ = r.sm.CSRF(c)
	files := []string{filepath.Join(r.root, "layouts", "base.html"), filepath.Join(r.root, "components", "sidebar.html"), filepath.Join(r.root, filepath.FromSlash(page))}
	t, err := template.New("base.html").Funcs(template.FuncMap{"hasRole": func(roles []string, target string) bool {
		for _, role := range roles {
			if role == target {
				return true
			}
		}
		return false
	}}).ParseFiles(files...)
	if err != nil {
		c.String(http.StatusInternalServerError, "template error")
		return
	}
	c.Status(status)
	_ = t.ExecuteTemplate(c.Writer, "base.html", d)
}
