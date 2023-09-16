package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/zsoltardai/simple_bank/api"
	db "github.com/zsoltardai/simple_bank/db/sqlc"
	"github.com/zsoltardai/simple_bank/util"
)

func main() {

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Couldn't load envioment variables!", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to the database!", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server!", err)
	}
}
