package main

import (
	"context"
	"semantic_api/api"
	db "semantic_api/db/sqlc"
	"semantic_api/mail"
	"semantic_api/util"
	"semantic_api/worker"

	"github.com/rs/zerolog/log"

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
		log.Fatal().Msg("cannot load config")
	}

	log.Info().Msg("loaded config")

	// Postgres connection
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	log.Info().Msg("connected to pg db")

	// Set up a connection to the server
	conn, err := grpc.Dial(config.VectorGrpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Msgf("Failed to connect: %v", err)
	}

	log.Info().Msg("grpc connection established")

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatal().Msgf("Failed to close connection", err)
		}
	}(conn)

	// Create a client using the generated code
	// grpcClient := pb.NewVectorManagerClient(conn)

	store := db.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(config, redisOpt, store)
	server, err := api.NewServer(config, store, taskDistributor)

	if err != nil {
		log.Fatal().Msgf("error creating server", err)
	}
	err = server.RunHTTPServer(config.ServerAddress)

	if err != nil {
		log.Fatal().Msgf("cannot start server", err)
	}

}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msgf("cannot create new db migrate instance: ", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("failed to run migrate up: ", err)
	}
	log.Info().Msg("db migrated successfully")
}

func initMilvusClient(milvusAddress string) (client.Client, error) {
	milvusClient, err := client.NewGrpcClient(
		context.Background(),
		milvusAddress,
	)
	return milvusClient, err
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Msg("failed to start task processor")
	}
}
