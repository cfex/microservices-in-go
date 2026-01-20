package main

import (
	"time"

	"github.com/cfex/microservices-in-go/services/common/amqp"
	amqpConfig "github.com/cfex/microservices-in-go/services/common/amqp/config"
	api "github.com/cfex/microservices-in-go/services/users-service/cmd"
	"github.com/cfex/microservices-in-go/services/users-service/cmd/config"
	"github.com/cfex/microservices-in-go/services/users-service/internal/database"
	"github.com/cfex/microservices-in-go/services/users-service/internal/jwt"
	"github.com/cfex/microservices-in-go/services/users-service/internal/logger"
	"github.com/cfex/microservices-in-go/services/users-service/server"
)

func main() {
	loc, err := time.LoadLocation("Europe/Belgrade")
	if err != nil {
		panic(err)
	}
	time.Local = loc
	logger.Init()
	log := logger.GetLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	err = cfg.Validate()
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	dbConn := database.ConnectDB(*cfg)
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()

	amqpConfig, err := amqpConfig.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load amqp configuration")
	}

	conn, err := amqp.ConnectAmqp(amqpConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to amqp")
	}
	defer conn.Close()

	client, err := amqp.NewClient(conn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create a client")
	}
	producer := amqp.NewProducer(client)

	jwt.NewJwt(cfg)

	router := api.CreateRouter(dbConn, cfg, producer).NewRouter()

	go func() {
		if err := server.NewGRPCServer(dbConn, cfg).Run(); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	server.NewServer(log, router, cfg).Serve()
}
