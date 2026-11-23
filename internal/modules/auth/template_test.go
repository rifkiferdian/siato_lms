package auth

import (
	"html/template"
	"path/filepath"
	"strings"
	"testing"
)

func TestAccountTemplatesParse(t *testing.T) {
	root := filepath.Join("..", "..", "..", "web", "templates")
	for _, page := range []string{"auth/forgot.html", "auth/reset.html", "account/index.html"} {
		_, err := template.New("base.html").Funcs(template.FuncMap{"hasRole": func([]string, string) bool { return false }, "hasPrefix": strings.HasPrefix}).ParseFiles(filepath.Join(root, "layouts", "base.html"), filepath.Join(root, "components", "sidebar.html"), filepath.Join(root, filepath.FromSlash(page)))
		if err != nil {
			t.Fatalf("%s: %v", page, err)
		}
	}
}
