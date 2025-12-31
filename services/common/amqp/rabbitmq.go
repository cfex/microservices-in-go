package amqp

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

var ProjectCreateEvent = "projectCreateEvent"

type AmqpServer struct {
	cfg *Config
}

type Connection struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func Connect(cfg *Config) *AmqpServer {
	return &AmqpServer{cfg: cfg}
}

func (a *AmqpServer) ConnectAmqp() (*Connection, func() error) {
	addr := a.cfg.AmqpConnectinoString()
	conn, err := amqp091.Dial(addr)

	if err != nil {
		log.Fatal(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &Connection{conn, channel}, nil
}

func (c *Connection) Close() error {
	if c.ch != nil {
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return fmt.Errorf("error clossing connection: %w", err)
		}
	}

	log.Panicln("MQ Connection closed successfully")
	return nil
}
