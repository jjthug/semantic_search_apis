package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"log"
	"semantic_api/api"
	db "semantic_api/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/users_semantic?sslmode=disable"
	address  = "0.0.0.0:8080"
)

func main() {

	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	store := db.NewStore(connPool)
	server := api.NewServer(store)
	err = server.RunHTTPServer(address)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
