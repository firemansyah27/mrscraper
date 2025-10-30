package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const cacheExpiration = 24 * time.Hour

type OrderRepositoryInterface interface {
	GetOrdersByProduct(productID int) ([]models.Order, error)
}

type EventPublisherInterface interface {
	EmitEvent(event string, data interface{}) error
}

type OrderService struct {
	orderRepo      OrderRepositoryInterface
	eventPublisher EventPublisherInterface
	redisClient    *redis.Client
	productSvcURL  string
}

func NewOrderService(
	orderRepo OrderRepositoryInterface,
	eventPublisher EventPublisherInterface,
	redisClient *redis.Client,
	productSvcURL string,
) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
		redisClient:    redisClient,
		productSvcURL:  productSvcURL,
	}
}

func (s *OrderService) CreateOrder(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid payload",
			"detail": err.Error(),
		})
		return
	}

	eventData := map[string]interface{}{
		"product_id": req.ProductID,
		"quantity":   req.Quantity,
	}

	if err := s.eventPublisher.EmitEvent("order.created", eventData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to publish order.created event",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "order.created event published successfully",
		"event":   eventData,
	})
}

func (s *OrderService) GetOrdersByProduct(c *gin.Context) {
	productIDStr := c.Param("id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	orders, err := s.GetOrdersFromCache(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders from cache", "detail": err.Error()})
		return
	}

	if len(orders) == 0 {
		orders, err = s.orderRepo.GetOrdersByProduct(productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get orders from DB", "detail": err.Error()})
			return
		}

		if len(orders) > 0 {
			s.SetOrdersInCache(productID, orders)
		}
	}

	c.JSON(http.StatusOK, orders)
}

func (s *OrderService) GetOrdersFromCache(productID int) ([]models.Order, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("orders:product:%d", productID)

	cachedOrders, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var orders []models.Order
	if err := json.Unmarshal([]byte(cachedOrders), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) SetOrdersInCache(productID int, orders []models.Order) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("orders:product:%d", productID)

	ordersJSON, err := json.Marshal(orders)
	if err != nil {
		return err
	}

	err = s.redisClient.Set(ctx, cacheKey, ordersJSON, cacheExpiration).Err()
	if err != nil {
		return err
	}

	return nil
}
