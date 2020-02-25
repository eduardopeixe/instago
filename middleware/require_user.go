package middleware

import (
	"log"
	"net/http"

	"github.com/eduardopeixe/instago/models"
)

// RequireUser is the type to use UserService
type RequireUser struct {
	models.UserService
}

// Apply applies the middleware to http.Handler
func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}

// ApplyFn applies the middleware to http.HandlerFunc
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		// if user is logged in, just continue
		log.Println("User found:", user)
		next(w, r)
	})
}
