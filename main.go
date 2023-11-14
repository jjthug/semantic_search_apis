package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"semantic_api/api"
	db "semantic_api/db/sqlc"
	"semantic_api/pb"
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
	client, err := client.NewClient(context.Background(), client.Config{
		Address: "localhost:19530",
	})
	if err != nil {
		// handle error
	}
	defer client.Close()
	runGrpcServer(&client)

}

func runGrpcServer(client *client.Client) {

	server, err := gapi.NewServer(client)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	//server,err :=
	pb.RegisterVectorManagerServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRpc server:", err)
	}

}
