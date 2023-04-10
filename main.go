package main

import (
	"database/sql"
	"log"

	"github.com/Hypersus/simplebank/api"
	db "github.com/Hypersus/simplebank/db/sqlc"
	"github.com/Hypersus/simplebank/util"
	_ "github.com/lib/pq"
)

// const (
// 	dbDriver      = "postgres"
// 	dbSource      = "postgres://root:hypersus@localhost:5432/simple_bank?sslmode=disable"
// 	serverAddress = "localhost:8080"
// )

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("config: cannot load config", err)
	}
	sqlDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("database: cannot connect to database", err)
	}
	store := db.NewStore(sqlDB)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("server: cannot start server", err)
	}

}
