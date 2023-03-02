package main

import (
	"database/sql"
	"log"

	"github.com/Hypersus/simplebank/api"
	db "github.com/Hypersus/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://root:hypersus@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "localhost:8080"
)

func main() {

	sqlDB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("database: cannot connect to database", err)
	}
	store := db.NewStore(sqlDB)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("server: cannot start server", err)
	}
}
