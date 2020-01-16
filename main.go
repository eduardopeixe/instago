package main

import (
	"fmt"
	"net/http"
)

const port = ":3000"


func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Welcome to instaGO!</h1>")
	fmt.Fprint(w, "<p>Please send an email to <a href=#>support@instago.com</a></p>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	fmt.Println("Serving port", port)
	http.ListenAndServe(port, nil)
}