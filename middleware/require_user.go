package middleware

import (
	"log"
	"net/http"

	"github.com/eduardopeixe/instago/context"
	"github.com/eduardopeixe/instago/models"
)

// User is just a regular user middleware
type User struct {
	models.UserService
}

// Apply applies the middleware to http.Handler
func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn applies the middleware to http.HandlerFunc
func (mw *User) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(w, r)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		log.Println("User found:", user)
		next(w, r)
	})
}

// RequireUser is the type to use UserService
type RequireUser struct {
	User
}

// Apply applies the middleware to http.Handler
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn applies the middleware to http.HandlerFunc
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	ourHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())

		if user == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(w, r)
	})
	return mw.User.Apply(ourHandler)
}
