package rps

import (
	"html/template"
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplateParses(t *testing.T) {
	root := filepath.Join("..", "..", "..", "web", "templates")
	_, e := template.New("base.html").Funcs(template.FuncMap{"hasRole": func([]string, string) bool { return false }, "hasPrefix": strings.HasPrefix}).ParseFiles(filepath.Join(root, "layouts", "base.html"), filepath.Join(root, "components", "sidebar.html"), filepath.Join(root, "rps", "index.html"))
	if e != nil {
		t.Fatal(e)
	}
}
