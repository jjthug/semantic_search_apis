package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"semantic_api/api"
	db "semantic_api/db/sqlc"
	"semantic_api/util"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/users_semantic?sslmode=disable"
	address  = "0.0.0.0:8080"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store)
	err = server.RunHTTPServer(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
