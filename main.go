package main

import (
	"context"
	"fmt"
	"log"
	"semantic_api/api"
	db "semantic_api/db/sqlc"
	"semantic_api/util"
	"semantic_api/worker"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"google.golang.org/grpc"
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

	runDBMigration(config.MigrationURL, config.DBSource)

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
	// grpcClient := pb.NewVectorManagerClient(conn)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(redisOpt, store)
	server, err := api.NewServer(config, store, taskDistributor)

	if err != nil {
		log.Fatal("error creating server", err)
	}
	err = server.RunHTTPServer(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}

}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new db migrate instance: ", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	}
	log.Println("db migrated successfully")
}

func initMilvusClient(milvusAddress string) (client.Client, error) {
	milvusClient, err := client.NewGrpcClient(
		context.Background(),
		milvusAddress,
	)
	return milvusClient, err
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	fmt.Print("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("failed to start task processor")
	}
}
