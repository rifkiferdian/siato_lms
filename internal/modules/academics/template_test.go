package academics

import (
	"html/template"
	"path/filepath"
	"strings"
	"testing"
)

func TestAcademicTemplateParses(t *testing.T) {
	root := filepath.Join("..", "..", "..", "web", "templates")
	_, err := template.New("base.html").Funcs(template.FuncMap{
		"hasRole":   func([]string, string) bool { return false },
		"hasPrefix": strings.HasPrefix,
	}).ParseFiles(
		filepath.Join(root, "layouts", "base.html"),
		filepath.Join(root, "components", "sidebar.html"),
		filepath.Join(root, "academics", "index.html"),
	)
	if err != nil {
		t.Fatal(err)
	}
}
