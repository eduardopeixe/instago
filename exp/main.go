package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "peixe"
	dbname = "instago_dev"
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	Color string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	db.LogMode(true)
	db.AutoMigrate(&User{})

	var newUser User

	var u User = User{
		Name:  "john",
		Email: "john@smith.com",
	}
	db.Where(u).First(&u) //returns where a user.name == u.Name and user.email == u.Email
	fmt.Println(u)
	db.First(&newUser) // returns first record PK ascending
	fmt.Println(newUser)
	db.Last(&newUser) // returns first record PK descending
	fmt.Println(newUser)

}
