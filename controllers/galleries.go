package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/eduardopeixe/instago/context"
	"github.com/eduardopeixe/instago/models"
	"github.com/eduardopeixe/instago/views"
	"github.com/gorilla/mux"
)

// NewGallery creates a new Gallery controller.
func NewGallery(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		New:      views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
		gs:       gs,
		r:        r,
	}
}

// Galleries is the type of Gallery
type Galleries struct {
	New      *views.View
	ShowView *views.View
	gs       models.GalleryService
	r        *mux.Router
}

// GalleryForm is the model for Gallery
type GalleryForm struct {
	Title string `json:"title"`
}

func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
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

	user := context.User(r.Context())

	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	err := g.gs.Create(&gallery)
	log.Println("setting error", err)

	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	url, err := g.r.Get("show_gallery").URL("id", fmt.Sprintf("%d", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)

}
