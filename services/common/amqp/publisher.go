package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn *AmqpClient
}

func NewProducer(conn *AmqpClient) *Producer {
	return &Producer{conn: conn}
}

func (p *Producer) DeclareExchange(exchangeName, exchangeType string) error {
	err := p.conn.ch.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)

	if err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", exchangeName, err)
	}

	log.Printf("Exchange '%s' declared successfully", exchangeName)
	return nil
}

func (p *Producer) Publish(ctx context.Context, exchange, routingKey string, message any) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.conn.ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp091.Persistent,
		Timestamp:    time.Now(),
	})

	if err != nil {
		return fmt.Errorf("failed to publish message to %s with key %s: %w", exchange, routingKey, err)
	}

	log.Printf("Published message to exchange=%s, routingKey=%s, size=%d bytes",
		exchange, routingKey, len(body))
	return nil

}
