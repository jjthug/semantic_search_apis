package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"log"
	"semantic_api/api"
	db "semantic_api/db/sqlc"
	"semantic_api/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	// Postgres connection
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

	// Milvus connection

	// Initialize Milvus client
	milvusClient, err := initMilvusClient()
	if err != nil {
		log.Fatal("cannot create client", err)
	}

	defer func(milvusClient client.Client) {
		err := milvusClient.Close()
		if err != nil {
			log.Fatal("error closing client", err)
		}
	}(milvusClient)
}

func initMilvusClient() (client.Client, error) {
	milvusClient, err := client.NewClient(context.Background(), client.Config{
		Address: "localhost:19530",
	})
	return milvusClient, err
}
