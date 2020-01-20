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

	var id int
	var name, email string

	row := db.QueryRow(`
		select id, name, email from users where id =$1`, 1)

	if err := row.Scan(&id, &name, &email); err != nil {
		panic(err)
	}

	fmt.Println("The new ID is", id)

}
