package events

import (
	"encoding/json"
	"fmt"
	"log"
	"order-service/internal/repository"
	"github.com/streadway/amqp"
)

type ProductConsumer struct {
	amqpURL   string
	orderRepo *repository.OrderRepository
}

func NewProductConsumer(amqpURL string, repo *repository.OrderRepository) *ProductConsumer {
	return &ProductConsumer{
		amqpURL:   amqpURL,
		orderRepo: repo,
	}
}

func (pc *ProductConsumer) StartConsumer() {
	pc.consume("update-order-status-queue", "update.order.status")
}

func (pc *ProductConsumer) consume(queueName, routingKey string) {
	conn, err := amqp.Dial(pc.amqpURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare("events", "topic", true, false, false, false, nil); err != nil {
		log.Printf("Exchange declare error: %v", err)
		return
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Printf("Queue declare error: %v", err)
		return
	}

	if err := ch.QueueBind(q.Name, routingKey, "events", false, nil); err != nil {
		log.Printf("Queue bind error: %v", err)
		return
	}

	log.Printf("Listening on queue '%s' (routing key: '%s')", q.Name, routingKey)

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("Consume error: %v", err)
		return
	}

	for msg := range msgs {
		pc.handleMessage(msg)
	}

	log.Println("Consumer stopped.")
}

func (pc *ProductConsumer) handleMessage(m amqp.Delivery) {
	log.Printf("Message received (Tag: %d, Size: %d bytes)", m.DeliveryTag, len(m.Body))

	var event struct {
		Event     string `json:"event"`
		Timestamp string `json:"timestamp"`
		Data      struct {
			OrderID int    `json:"orderId"`
			Status  string `json:"status"`
		} `json:"data"`
	}

	if err := json.Unmarshal(m.Body, &event); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		m.Nack(false, false)
		return
	}

	orderID := event.Data.OrderID
	status := event.Data.Status
	log.Printf("Processing order ID: %d with new status: '%s'", orderID, status)

	if err := pc.updateOrderStatus(orderID, status); err != nil {
		log.Printf("Failed to update order %d: %v", orderID, err)
		m.Nack(false, true)
		return
	}

	m.Ack(false)
	log.Printf("Order %d updated to '%s'", orderID, status)
}

func (pc *ProductConsumer) updateOrderStatus(orderID int, status string) error {
	if err := pc.orderRepo.UpdateOrderStatus(orderID, status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
