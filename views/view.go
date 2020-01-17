package views

import (
	"html/template"
)

// NewView creates a view
func NewView(layout string, files ...string) *View {
	files = append(files, 
		"views/layouts/bootstrap.gohtml", 
		"views/layouts/footer.gohtml",
		"views/layouts/navbar.gohtml",
	)

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
