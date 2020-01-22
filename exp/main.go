package main

import (
	"fmt"
	"github.com/eduardopeixe/instago/models"

	_ "github.com/jinzhu/gorm/dialects/postgres"
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

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}

	defer us.Close()
	us.ResetDB()

	// user, err := us.ByID(1)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println(user)

}
