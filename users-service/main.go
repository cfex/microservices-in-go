package main

import (
	api "github.com/cfex/microservices-in-go/users-service/cmd"
	"github.com/cfex/microservices-in-go/users-service/cmd/config"
	"github.com/cfex/microservices-in-go/users-service/internal/database"
	"github.com/cfex/microservices-in-go/users-service/internal/logger"
	"github.com/cfex/microservices-in-go/users-service/server"
)

func main() {

	log := logger.GetLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	dbConn := database.ConnectDB(*cfg)
	defer dbConn.Close()

	router := api.CreateRouter(dbConn, cfg).NewRouter()

	server.NewServer(log, router, cfg).Serve()
}
