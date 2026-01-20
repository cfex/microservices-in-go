package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cfex/microservices-in-go/services/common/amqp"
	amqpConfig "github.com/cfex/microservices-in-go/services/common/amqp/config"
	amqpConsts "github.com/cfex/microservices-in-go/services/common/amqp/consts"
	"github.com/cfex/microservices-in-go/services/notification-service/config"
	"github.com/cfex/microservices-in-go/services/notification-service/internal/handlers"
	"github.com/cfex/microservices-in-go/services/notification-service/internal/services"
	"github.com/cfex/microservices-in-go/services/notification-service/logger"
)

func main() {
	log := logger.GetLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
	}

	err = cfg.Validate()
	if err != nil {
		log.Fatal().Err(err).Msg("configuration validation failed")
	}

	amqpConfig, err := amqpConfig.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load amqp configuration")
	}

	conn, err := amqp.ConnectAmqp(amqpConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to amqp")
	}
	defer conn.Close()

	client, err := amqp.NewClient(conn)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create amqp client")
	}

	consumer := amqp.NewConsumer(client)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create amqp consumer")
	}

	s := services.NewSender(cfg)
	h := handlers.NewEmailHandler(*s, cfg, &log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = consumer.Consume(ctx, amqpConsts.EmailQueue, "notification-worker", h.HandleMessage)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start consuming messages")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Info().Msgf("Received signal %s, shutting down...", sig)

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	select {
	case <-shutdownCtx.Done():
		log.Info().Msg("Shutdown timed out, forcing exit")
	case <-time.After(2 * time.Second):
		log.Info().Msg("Shutdown complete")
	}

	if err := conn.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close AMQP connection")
	}

	log.Info().Str("signal", sig.String()).Msg("Received signal, shutting down")
}
