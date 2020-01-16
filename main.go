package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const port = ":3000"

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch r.URL.Path {
	case "/":
	case "/contact":
	default:
	}
}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Welcome to instaGO!</h1>")

}

func contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<p>Please send an email to <a href=#>support@instago.com</a></p>")

}

func faq(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "this is the FAQ page")

}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>404 Page not found</h1>")
}

func main() {
	r := httprouter.New()
	r.NotFound = http.HandlerFunc(notFound)
	r.GET("/", home)
	r.GET("/contact", contact)
	r.GET("/faq", faq)
	fmt.Println("Serving port", port)
	http.ListenAndServe(port, r)
}
