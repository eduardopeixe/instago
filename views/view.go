package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	layoutDir   = "views/layouts/"
	templateDir = "views/"
	templateExt = ".gohtml"
)

// NewView creates a view
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)

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
	v.Render(w, nil)
}

// Render is used to render the view with predefined layout
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	switch data.(type) {
	case Data:
		//Do nothing
	default:
		data = Data{
			Yield: data,
		}
	}
	//TODO: Fix this
	// var buf *bytes.Buffer
	return v.Template.ExecuteTemplate(w, v.Layout, data)
	// if err != nil{
	// http.Error(w, "Somenthing went really wrong", http.StatusInternalServerError)
	// return
	// }
	// io.Copy(w, &buf)
}

// layoutFile return a slice of strings representing the layout files
func layouts() []string {
	files, err := filepath.Glob(layoutDir + "*" + templateExt)
	if err != nil {
		panic(err)
	}

	return files
}

// addTemplatePath takes a slice of strings representation file paths for
// templates and prepends the templateDir to each string in the slice
//
// Eg. the input {"home"} results in {"views/home"} if TemplateDir == "views/"
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = templateDir + f
	}
}

// addTemplateExt taskes a slice of strings representation file paths for templates
// and append the templateExt to each string in the slice
//
// Eg. the input {"home"} results in {"home.gohtml"} if templateExt == ".gohtml"
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + templateExt
	}
}
