package main

import (
	"fmt"
	"net/http"

	"github.com/eduardopeixe/instago/controllers"
	"github.com/eduardopeixe/instago/views"
	"github.com/gorilla/mux"
)

const port = ":3000"

var (
	homeView    *views.View
	contactView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404 Page not found</h1>")
}

func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/signup", usersC.New)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	fmt.Println("Serving port", port)
	http.ListenAndServe(port, r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
