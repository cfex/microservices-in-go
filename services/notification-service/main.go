package main

import (
	"context"
	"log"
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
)

func main() {
	cfg := config.InitConfig()

	amqpConfig := amqpConfig.Load()
	conn, err := amqp.ConnectAmqp(amqpConfig)
	if err != nil {
		log.Fatalf("failed to connect to amqp: %v", err)
	}
	defer conn.Close()

	client, err := amqp.NewClient(conn)
	if err != nil {
		log.Fatalf("failed to create a client: %v", err)
	}

	consumer := amqp.NewConsumer(client)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	s := services.NewSender(cfg)
	h := handlers.NewEmailHandler(*s, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = consumer.Consume(ctx, amqpConsts.EmailQueue, "notification-worker", h.HandleMessage)
	if err != nil {
		log.Fatal("failed to start consumer: %w", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %v, shutting down gracefully...", sig)

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	select {
	case <-shutdownCtx.Done():
		log.Println("Shutdown timeout exceeded, force closing")
	case <-time.After(2 * time.Second):
		log.Println("Consumer shutdown complete")
	}

	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}

	log.Println("Notification service stopped")
}
