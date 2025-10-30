package events

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"order-service/internal/models"
	"order-service/internal/repository"
	"time"

	"github.com/streadway/amqp"
)

type OrderConsumer struct {
	amqpURL       string
	orderRepo     *repository.OrderRepository
	productSvcURL string
}

func NewOrderConsumer(amqpURL string, repo *repository.OrderRepository, productSvcURL string) *OrderConsumer {
	return &OrderConsumer{
		amqpURL:       amqpURL,
		orderRepo:     repo,
		productSvcURL: productSvcURL,
	}
}

func (oc *OrderConsumer) StartConsumer() {
	conn, err := amqp.Dial(oc.amqpURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v. Order consumer stopped.", err)
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

	q, err := ch.QueueDeclare("order-created-queue", true, false, false, false, nil)
	if err != nil {
		log.Printf("Queue declare error: %v", err)
		return
	}

	if err := ch.QueueBind(q.Name, "order.created", "events", false, nil); err != nil {
		log.Printf("Queue bind error: %v", err)
		return
	}

	log.Println("Order consumer started. Waiting for messages on queue:", q.Name)

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Printf("Consume error: %v", err)
		return
	}

	for msg := range msgs {
		var event struct {
			Event     string `json:"event"`
			Timestamp string `json:"timestamp"`
			Data      struct {
				ProductID int `json:"product_id"`
				Quantity  int `json:"quantity"`
			} `json:"data"`
		}

		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			msg.Nack(false, false)
			continue
		}

		log.Printf("Received create.order: %+v", event.Data)

		productURL := fmt.Sprintf("%s/products/%d", oc.productSvcURL, event.Data.ProductID)
		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(productURL)
		if err != nil {
			log.Printf("Failed to connect to product service: %v", err)
			msg.Nack(false, true)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			io.Copy(io.Discard, resp.Body)
			log.Printf("Product not found or invalid response for productID=%d", event.Data.ProductID)
			msg.Nack(false, true)
			continue
		}

		var product map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
			log.Printf("Failed to decode product response: %v", err)
			msg.Nack(false, true)
			continue
		}

		price, ok := product["price"].(float64)
		if !ok {
			log.Printf("⚠️ Invalid product price for productID=%d", event.Data.ProductID)
			msg.Nack(false, false)
			continue
		}

		order := models.Order{
			ProductID: event.Data.ProductID,
			Quantity:  event.Data.Quantity,
			Total:     float64(event.Data.Quantity) * price,
			Status:    "draft",
		}

		if err := oc.orderRepo.CreateOrder(&order); err != nil {
			log.Printf("Failed to create order in DB: %v", err)
			msg.Nack(false, true)
			continue
		}

		eventData := map[string]interface{}{
			"order_id":   order.ID,
			"product_id": order.ProductID,
			"quantity":   order.Quantity,
			"total":      order.Total,
			"status":     order.Status,
		}

		publisher := NewEventPublisher(oc.amqpURL)
		if err := publisher.EmitEvent("update.product.stock", eventData); err != nil {
			log.Printf("Failed to publish update.product.stock event: %v", err)
		} else {
			log.Printf("Published update.product.stock event for order_id=%d", order.ID)
		}

		log.Printf("Order created: ID=%d, ProductID=%d, Qty=%d, Total=%.2f",
			order.ID, order.ProductID, order.Quantity, order.Total)

		msg.Ack(false)
	}

	log.Println("Order consumer stopped.")
}
