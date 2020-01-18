package controllers

import "github.com/eduardopeixe/instago/views"

// NewStatic create new static page views
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap", "views/static/contact.gohtml"),
	}
}

// Static words
type Static struct {
	Home    *views.View
	Contact *views.View
}
