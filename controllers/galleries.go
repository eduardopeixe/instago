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
		New:       views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		gs:        gs,
		r:         r,
	}
}

// Galleries is the type of Gallery
type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	r         *mux.Router
}

// GalleryForm is the model for Gallery
type GalleryForm struct {
	Title string `json:"title"`
}

// Show displays a gallery
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
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

// Delete deletes a gallery by gallery ID
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		log.Println("Not owner delete")
		http.Error(w, "You cannont delete this gallery", http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = gallery

	err = g.gs.Delete(gallery.ID)
	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}

	vd.AlertSuccess("Gallery deleted successfully")
	//TODO: redirect to index gallery page
	fmt.Fprintln(w, "deleted")
}

// Update POST /galleries/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		log.Println("Not owner update")
		http.Error(w, "You cannot edit this gallery", http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = gallery

	var form GalleryForm
	err = parseForm(r, &form)
	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}

	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}

	vd.AlertSuccess("Gallery updated successfully")
	g.ShowView.Render(w, vd)
}

// galleryByID returns a gallery that matches the param ID
func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {

	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Something went wrong", http.StatusNotFound)
		}
		return nil, err
	}

	return gallery, nil
}

// Edit GET /galleries/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		log.Println("Not owner")
		http.Error(w, "You cannot edit this gallery", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = gallery

	g.EditView.Render(w, vd)
}

// Index GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {

	user := context.User(r.Context())

	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		log.Println("Not owner")
		http.Error(w, "Cannot load galleries", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = galleries

	g.IndexView.Render(w, vd)
	// fmt.Fprintln(w, galleries)
}
