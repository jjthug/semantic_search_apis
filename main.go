package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"google.golang.org/grpc"
	"log"
	"semantic_api/api"
	db "semantic_api/db/sqlc"
	"semantic_api/pb"
	"semantic_api/util"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load config")
		log.Fatal("cannot load config")
	}

	fmt.Println("loaded config")

	// Postgres connection
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		fmt.Println("cannot connect to db")
		log.Fatal("cannot connect to db")
	}

	fmt.Println("connected to pg db")

	// Set up a connection to the server
	conn, err := grpc.Dial(config.VectorGrpcAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect")
		log.Fatalf("Failed to connect: %v", err)
	}

	fmt.Println("grpc connection established")

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Failed to close connection", err)
		}
	}(conn)

	// Create a client using the generated code
	grpcClient := pb.NewVectorManagerClient(conn)

	store := db.NewStore(connPool)

	// Milvus connection

	// Initialize Milvus client
	//milvusClient, err := initMilvusClient(config.MilvusAddr)
	//if err != nil {
	//	log.Fatal("cannot create client", err)
	//}
	//
	//fmt.Println("milvus connection established")
	//
	//defer func(milvusClient client.Client) {
	//	err := milvusClient.Close()
	//	if err != nil {
	//		log.Fatal("error closing client", err)
	//	}
	//}(milvusClient)

	//server, err := api.NewServer(config, store, &grpcClient, &milvusClient)
	server, err := api.NewServer(config, store, &grpcClient)

	if err != nil {
		log.Fatal("error creating server", err)
	}
	err = server.RunHTTPServer(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}

}

func initMilvusClient(milvusAddress string) (client.Client, error) {
	milvusClient, err := client.NewGrpcClient(
		context.Background(),
		milvusAddress,
	)
	return milvusClient, err
}
