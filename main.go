package main

import (
	"database/sql"
	"log"

	"github.com/dralos/simplebank/api"
	db "github.com/dralos/simplebank/db/sqlc"
	"github.com/dralos/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Default().Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Default().Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Default().Fatal("cannot start server:", err)
	}
}
