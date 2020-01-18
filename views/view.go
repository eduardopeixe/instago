package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	layoutDir   = "views/layouts/"
	templateExt = ".gohtml"
)

// NewView creates a view
func NewView(layout string, files ...string) *View {
	files = append(files, layouts()...)

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{t, layout}
}

//View is the type of a view
type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// Render is used to render the view with predefined layout
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

// layoutFile return a slice of strings representing the layout files
func layouts() []string {
	files, err := filepath.Glob(layoutDir + "*" + templateExt)
	if err != nil {
		panic(err)
	}

	return files
}
