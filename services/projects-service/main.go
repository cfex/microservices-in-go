package main

import (
	api "github.com/cfex/microservices-in-go/services/projects-service/cmd"
	"github.com/cfex/microservices-in-go/services/projects-service/cmd/config"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/clients"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/database"
	"github.com/cfex/microservices-in-go/services/projects-service/internal/logger"
	"github.com/cfex/microservices-in-go/services/projects-service/server"
)

func main() {
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

	cl, err := clients.NewUserClient(cfg)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	router := api.CreateRouter(dbConn, cfg, cl).NewRouter()

	server.NewServer(log, router, cfg).Serve()
}
