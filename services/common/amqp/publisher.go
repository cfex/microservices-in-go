package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Publsher struct {
	conn *Connection
}

func NewPublisher(conn *Connection) *Publsher {
	return &Publsher{conn: conn}
}

func (p *Publsher) DeclareExchange(exchangeName, exchangeType string) error {
	err := p.conn.ch.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)

	if err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", exchangeName, err)
	}

	log.Printf("Exchange '%s' declared successfully", exchangeName)
	return nil
}

func (p *Publsher) Publish(ctx context.Context, exchange, routingKey string, message any) error {
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
