package main

import (
	"fmt"
	"net/http"

	"github.com/eduardopeixe/instago/controllers"
	"github.com/eduardopeixe/instago/models"
	"github.com/gorilla/mux"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "peixe"
	dbname = "instago_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)

	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}

	must(err)

	defer services.Close()
	// services.DestructiveReset()
	services.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGallery(services.Gallery)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	// Gallery routes
	r.Handle("/galleries/new", galleriesC.New).Methods("GET")
	r.HandleFunc("/galleries/create", galleriesC.Create).Methods("POST")
	fmt.Println("Serving port", 3000)
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
