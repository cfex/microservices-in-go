package amqp

import (
	"github.com/cfex/microservices-in-go/services/common/amqp/config"
	"github.com/rabbitmq/amqp091-go"
)

var ProjectCreateEvent = "projectCreateEvent"

type AmqpClient struct {
	conn *amqp091.Connection
	ch   *amqp091.Channel
}

func ConnectAmqp(cfg *config.Config) (*amqp091.Connection, error) {
	addr := cfg.AmqpConnectinoString()

	return amqp091.Dial(addr)
}

func NewClient(conn *amqp091.Connection) (*AmqpClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return &AmqpClient{}, err
	}

	return &AmqpClient{conn: conn, ch: ch}, nil
}

func (ac *AmqpClient) Close() error {
	return ac.ch.Close()
}
