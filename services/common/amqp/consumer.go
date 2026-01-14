package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *AmqpClient
}

func NewConsumer(conn *AmqpClient) *Consumer {
	return &Consumer{conn: conn}
}

type MessageHandler func(body []byte) error

func (c *Consumer) DeclareQueue(queueName string, durable, autoDelete, exclusive bool) error {
	_, err := c.conn.ch.QueueDeclare(queueName, durable, autoDelete, exclusive, false, nil)

	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	log.Printf("queueu %s declared successfully", queueName)
	return nil
}

func (c *Consumer) BindQueue(queueName, exchangeName, routingKey string) error {
	err := c.conn.ch.QueueBind(queueName, routingKey, exchangeName, false, nil)

	if err != nil {
		return fmt.Errorf("failed to bind queue %s to exchange %s: %w", queueName, exchangeName, err)
	}

	log.Printf("Queue '%s' bound to exchange '%s' with routing key '%s'", queueName, exchangeName, routingKey)
	return nil
}

func (c *Consumer) Consume(ctx context.Context, queueName, tag string, handler MessageHandler) error {
	msgs, err := c.conn.ch.Consume(queueName, tag, false, false, false, false, nil)

	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Consumer '%s' started, waiting for messages from queue '%s'...", tag, queueName)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Consumer '%s' shutting down: %v", tag, ctx.Err())
				return

			case msg, ok := <-msgs:

				if !ok {
					log.Printf("Message channel closed for consumer '%s'", tag)
					return
				}

				c.handleMessage(msg, handler)
			}
		}
	}()

	return nil
}

func (c *Consumer) handleMessage(msg amqp091.Delivery, handler MessageHandler) {
	err := handler(msg.Body)

	if err != nil {
		log.Printf("Handler error: %v", err)

		nackErr := msg.Nack(
			false,
			true,
		)

		if nackErr != nil {
			log.Printf("Failed to NACK message: %v", nackErr)
		} else {
			log.Printf("Message NACK-ed and requeued for retry")
		}

		return
	}

	ackErr := msg.Ack(
		false,
	)

	if ackErr != nil {
		log.Printf("Failed to ACK message: %v", ackErr)
	} else {
		log.Printf("Message processed successfully and ACK-ed")
	}
}

func ConsumeJSON[T any](c *Consumer, ctx context.Context, queueName, tag string, handler func(*T) error) error {
	wrap := func(body []byte) error {
		var event T

		if err := json.Unmarshal(body, &event); err != nil {
			return fmt.Errorf("failed to unmarshal json: %w", err)
		}

		return handler(&event)
	}

	return c.Consume(ctx, queueName, tag, wrap)
}
