package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eduardopeixe/instago/models"
	"github.com/eduardopeixe/instago/views"
)

// NewGallery creates a new Gallery controller.
func NewGallery(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

// Galleries is the type of Gallery
type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

// GalleryForm is the model for Gallery
type GalleryForm struct {
	Title string `json:"title"`
}

// Create creates a new Gallery
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {

	var vd views.Data

	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: 1,
	}

	err := g.gs.Create(&gallery)
	log.Println("setting error", err)

	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	//TODO: Redirect to just created Gallery
	fmt.Fprintln(w, gallery)
}
