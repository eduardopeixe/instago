package controllers

import (
	"fmt"
	"net/http"

	"github.com/eduardopeixe/instago/models"
	"github.com/eduardopeixe/instago/rand"
	"github.com/eduardopeixe/instago/views"
)

// NewUsers creates a new Users controller. This functon will panic
// if the templates are not passed correctly, and should only be used
// during the initial setup
func NewUsers(us models.UserService) *Users {
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
	us        models.UserService
}

// New creates a new user view
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	type Alert struct {
		Level   string
		Message string
	}

	a := Alert{
		"success",
		"success to test alert",
	}
	err := u.NewView.Render(w, a)
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
	err = u.signIn(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cookietest", http.StatusFound)
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

	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)

	return nil
}

// CookieTest is used to display cookies set on the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "remember is:", cookie.Value)
	fmt.Fprintln(w, user)
}
