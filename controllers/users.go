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
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

// Users is the type of users
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Create is used to process the signup form
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	err := u.us.Create(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password" `
}

// Login us used to verify te provided emai addresss and password are
// valid to login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	form := LoginForm{}
	err := parseForm(r, &form)
	if err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address")
		case models.ErrIncorrectPassword:
			fmt.Fprintln(w, "Incorrect password provided")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
	}

	http.SetCookie(w, &cookie)
	fmt.Fprintln(w, user)
}
