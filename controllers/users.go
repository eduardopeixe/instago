package controllers

import (
	"fmt"
	"net/http"

	"github.com/eduardopeixe/instago/models"
	"github.com/eduardopeixe/instago/views"
)

// NewUsers creates a new Users controller. This functon will panic
// if the templates are not passed correctly, and should only be used
// during the initial setup
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us:      us,
	}
}

// Users is the type of users
type Users struct {
	NewView *views.View
	us      *models.UserService
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
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create is used to process the signup form
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	var form SignupForm
	if err := ParseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:  form.Name,
		Email: form.Email,
	}

	err := u.us.Create(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, form)
}
