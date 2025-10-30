package events

import (
	"log"
	"order-service/internal/repository"
)

func StartAllConsumers(amqpURL, productSvcURL string, orderRepo *repository.OrderRepository) {
	go func() {
		orderConsumer := NewOrderConsumer(amqpURL, orderRepo, productSvcURL)
		orderConsumer.StartConsumer()
	}()

	go func() {
		productConsumer := NewProductConsumer(amqpURL, orderRepo)
		productConsumer.StartConsumer()
	}()

	log.Println("Consumers started: order.created + update.order.status")
}
