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

	rows, err := db.Query(`
	select * from users
	inner join orders On users.id = orders.user_id
	`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var userID, orderID, amount int
		var name, email, description string
		if err := rows.Scan(&userID, &name, &email, &orderID, &userID, &amount, &description); err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		fmt.Println(rows.Err())
	}

}
