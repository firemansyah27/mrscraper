package events

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type EventPublisher struct {
	amqpURL      string
	exchangeName string
}

func NewEventPublisher(amqpURL string) *EventPublisher {
	return &EventPublisher{
		amqpURL:      amqpURL,
		exchangeName: "events",
	}
}

func (p *EventPublisher) EmitEvent(routingKey string, data interface{}) error {
	conn, err := amqp.Dial(p.amqpURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()


	err = ch.ExchangeDeclare(
		p.exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	payload := map[string]interface{}{
		"event":     routingKey,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"data":      data,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event payload: %w", err)
	}

	err = ch.Publish(
		p.exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf(" Published event [%s]: %s", routingKey, string(body))
	return nil
}
