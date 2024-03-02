package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vvv9912/sddb"
	"log"
)

func main() {

	database_dsn := "postgres://postgres:postgres@localhost:5432/tgbot?sslmode=disable"
	db, err := sqlx.Connect("postgres", database_dsn)
	if err != nil {
		log.Println("error connect config.Get().DatabaseDSN\n ", database_dsn, "db err:", db)
		return
	}

	if err := sddb.Migrate(db); err != nil {
		log.Fatal(err)
	}
}
