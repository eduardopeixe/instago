package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	type User struct {
		ID    int
		Name  string
		Email string
	}

	var users []User

	rows, err := db.Query(`
		select id, name, email from users`)

	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)

		users = append(users, user)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		fmt.Println("This is rows.Err()")
	}

	fmt.Println(users)

}
