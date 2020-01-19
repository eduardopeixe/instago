package controllers

import "github.com/eduardopeixe/instago/views"

// NewStatic create new static page views
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}
}

// Static words
type Static struct {
	Home    *views.View
	Contact *views.View
}
