package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/vvv9912/sddb/migrations"
	"log"
)

func main() {
	log.Println("Hello, World!")
	database_dsn := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sqlx.Connect("postgres", database_dsn)
	if err != nil {
		log.Println("error connect config.Get().DatabaseDSN\n ", database_dsn, "db err:", db)
		return
	}
	if err := migrations.Migrate(db); err != nil {
		log.Fatal(err)
	}
}
