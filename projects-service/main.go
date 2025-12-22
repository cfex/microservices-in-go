package main

import (
	api "github.com/cfex/microservices-in-go/projects-service/cmd"
	"github.com/cfex/microservices-in-go/projects-service/cmd/config"
	"github.com/cfex/microservices-in-go/projects-service/internal/database"
	"github.com/cfex/microservices-in-go/projects-service/internal/logger"
	"github.com/cfex/microservices-in-go/projects-service/server"
)

func main() {
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

	router := api.CreateRouter(dbConn, cfg).NewRouter()

	server.NewServer(log, router, cfg).Serve()
}
