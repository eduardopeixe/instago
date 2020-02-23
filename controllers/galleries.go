package controllers

import (
	"github.com/eduardopeixe/instago/models"
	"github.com/eduardopeixe/instago/views"
)

// NewGallery creates a new Gallery controller.
func NewGallery(gs models.GalleryService) *Gallery {
	return &Gallery{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

// Gallery is the type of Gallery
type Gallery struct {
	New *views.View
	gs  models.GalleryService
}
