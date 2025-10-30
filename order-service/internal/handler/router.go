package handler

import (
	"net/http"
	"time"
	"order-service/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRouter(svc *service.OrderService) *gin.Engine {
	r := gin.Default()

	ordersGroup := r.Group("/orders")
	{
		ordersGroup.POST("", CreateOrderHandler(svc))
		ordersGroup.GET("/product/:id", GetOrdersByProductHandler(svc))
	}

	r.GET("/health", HealthCheckHandler)

	return r
}

func CreateOrderHandler(svc *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		svc.CreateOrder(c)
	}
}

func GetOrdersByProductHandler(svc *service.OrderService) gin.HandlerFunc {
	return func(c *gin.Context) {
		svc.GetOrdersByProduct(c)
	}
}

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "order-service",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
