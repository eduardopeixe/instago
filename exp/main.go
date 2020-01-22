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
	Name   string
	Email  string `gorm:"not null;unique_index"`
	Color  string
	Orders []Order
}

type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
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
	db.AutoMigrate(&User{}, &Order{})

	var u User
	if err := db.First(&u).Error; err != nil {
		panic(err)
	}

	err = createOrder(db, u, 1199, "Description #1")
	if err != nil {
		panic(err)
	}
	err = createOrder(db, u, 999, "Description #2")
	if err != nil {
		panic(err)
	}
	err = createOrder(db, u, 4999, "Description #3")
	if err != nil {
		panic(err)
	}
}

func createOrder(db *gorm.DB, user User, amount int, desc string) error {
	return db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	}).Error
}
