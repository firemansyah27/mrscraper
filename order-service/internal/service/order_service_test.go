package service

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"order-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventPublisher struct {
	mock.Mock
}

func (m *MockEventPublisher) EmitEvent(event string, data interface{}) error {
	args := m.Called(event, data)
	return args.Error(0)
}

type DummyOrderRepo struct{}

func (r *DummyOrderRepo) GetOrdersByProduct(productID int) ([]models.Order, error) {
	return []models.Order{}, nil
}

func TestCreateOrder_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPublisher := new(MockEventPublisher)
	mockPublisher.On("EmitEvent", "order.created", mock.Anything).Return(nil)

	dummyRepo := &DummyOrderRepo{}
	dummyRedis := &redis.Client{}

	svc := &OrderService{
		orderRepo:      dummyRepo,
		eventPublisher: mockPublisher,
		redisClient:    dummyRedis,
		productSvcURL:  "http://dummy",
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	payload := `{"product_id":1,"quantity":5}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(payload))
	c.Request.Header.Set("Content-Type", "application/json")

	svc.CreateOrder(c)

	assert.Equal(t, http.StatusAccepted, w.Code, "HTTP status harus 202 Accepted")

	mockPublisher.AssertCalled(t, "EmitEvent", "order.created", mock.Anything)

	args := mockPublisher.Calls[0].Arguments
	data := args.Get(1).(map[string]interface{})
	assert.Equal(t, "1", fmt.Sprint(data["product_id"]))
	assert.Equal(t, "5", fmt.Sprint(data["quantity"]))
}
