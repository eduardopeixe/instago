package main

import (
	"fmt"
	"net/http"

	"github.com/eduardopeixe/instago/controllers"
	"github.com/gorilla/mux"
)

const port = ":3000"

func main() {
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers()

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	fmt.Println("Serving port", port)
	http.ListenAndServe(port, r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
