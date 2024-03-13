package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vvv9912/sddb"
	"log"
)

func main() {

	database_dsn := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sqlx.Connect("postgres", database_dsn)
	if err != nil {
		log.Println("error connect config.Get().DatabaseDSN\n ", database_dsn, "db err:", db)
		return
	}
	orders := sddb.NewOrdersPostgresStorage(db)
	massive, err := orders.GetOrderByStatus(context.Background(), 2)
	_ = massive
	if err := sddb.Migrate(db); err != nil {
		log.Fatal(err)
	}
}
