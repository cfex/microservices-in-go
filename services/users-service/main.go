package main

import (
	"time"

	api "github.com/cfex/microservices-in-go/services/users-service/cmd"
	"github.com/cfex/microservices-in-go/services/users-service/cmd/config"
	"github.com/cfex/microservices-in-go/services/users-service/internal/database"
	"github.com/cfex/microservices-in-go/services/users-service/internal/env"
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

	dbConn := database.ConnectDB(*cfg)
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection")
		}
	}()

	jwt.NewJwt(cfg)

	router := api.CreateRouter(dbConn, cfg).NewRouter()

	grpcServer := server.NewGRPCServer(env.GetEnv("USR_GRPC_ADDR"))
	grpcServer.Run()

	server.NewServer(log, router, cfg).Serve()
}
