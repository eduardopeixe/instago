package controllers

import (
	"fmt"
	"net/http"

	"github.com/eduardopeixe/instago/views"
	"github.com/gorilla/schema"
)

// NewUsers creates a new Users controller. This functon will panic
// if the templates are not passed correctly, and should only be used
// during the initial setup
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

// Users is the type of users
type Users struct {
	NewView *views.View
}

// New creates a new user view
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	err := u.NewView.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

// SignupForm is the model for signup form
type SignupForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create is used to process the signup form
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	dec := schema.NewDecoder()
	var form SignupForm
	if err := dec.Decode(&form, r.PostForm); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}
