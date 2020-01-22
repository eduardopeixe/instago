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

	var user User

	err = db.Where("email = ?", "no@existing.email").First(&user).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			fmt.Println("No user found!")
		case gorm.ErrInvalidSQL:
			fmt.Println("Invalid SQL command")
		default:
			fmt.Println("Database error", err)
		}
	}

	fmt.Println(user)
}
