package views

import (
	"html/template"
)

// NewView creates a view
func NewView(files ...string) *View {
	files = append(files, "views/layouts/footer.gohtml")

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{t}
}

//View is the type of a view
type View struct {
	Template *template.Template
}
