package views

import (
	"html/template"
	"path/filepath"
)

const (
	layoutDir = "views/layouts/"
	templateExt = ".gohtml" 
)

// NewView creates a view
func NewView(layout string, files ...string) *View {
	files = append(files, layoutfiles()...) 

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

// layoutFile return a slice of strings representing the layout files
func layouts()  []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}

	return files
}
