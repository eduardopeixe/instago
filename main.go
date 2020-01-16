package main

import (
	"fmt"
	"net/http"
)

const port = ":3000"

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, "<h1>Welcome to instaGO!</h1>")
	case "/contact":
		fmt.Fprint(w, "<p>Please send an email to <a href=#>support@instago.com</a></p>")
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>404 Page not found</h1>")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Serving port", port)
	http.ListenAndServe(port, nil)
}
