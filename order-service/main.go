package main

import (
	"log"
	"os"

	"order-service/internal/events"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/service"

	"github.com/go-redis/redis/v8"
)

func main() {
	dbInstance, err := repository.OpenDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable is not set")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	rdb := redis.NewClient(opt)
	defer rdb.Close()

	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		log.Fatal("RABBITMQ_URL environment variable is not set")
	}

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		log.Fatal("PRODUCT_SERVICE_URL environment variable is not set")
	}

	orderRepo := repository.NewOrderRepository(dbInstance)
	eventPublisher := events.NewEventPublisher(amqpURL)

	orderService := service.NewOrderService(
		orderRepo,
		eventPublisher,
		rdb,
		productServiceURL,
	)

	events.StartAllConsumers(amqpURL, productServiceURL, orderRepo)

	r := handler.SetupRouter(orderService)

	log.Println("🚀 Order service is running on port 4000 ...")
	if err := r.Run(":4000"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}
