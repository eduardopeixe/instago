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

	for i := 1; i < 6; i++ {
		user := 1
		if i > 2 {
			user = 2
		}
		amount := i * 13
		description := fmt.Sprintf("processor XYZ%d", amount)
		_, err := db.Exec(`
		insert into orders(user_id, amount, description) 
		VALUES($1, $2, $3)`,
			user, amount, description)
		if err != nil {
			panic(err)
		}
	}

}
