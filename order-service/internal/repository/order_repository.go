package repository

import (
	"os"
    "time"
	"fmt"
	"order-service/internal/models"
	"gorm.io/gorm"
    "gorm.io/driver/postgres"
)

type OrderRepository struct {
	db *gorm.DB
}

func OpenDB() (*gorm.DB, error) {
    dsn := os.Getenv("DATABASE_URL") 
    if dsn == "" {
        return nil, fmt.Errorf("DATABASE_URL environment variable not set")
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    err = db.AutoMigrate(&models.Order{})
    if err != nil {
        return nil, fmt.Errorf("failed to migrate database: %w", err)
    }

    sqlDB, err := db.DB()
    if err == nil {
        sqlDB.SetMaxIdleConns(10)
        sqlDB.SetConnMaxLifetime(time.Hour)
    }

    return db, nil
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
    if db == nil {
        fmt.Println("Warning: OrderRepository initialized with a nil DB connection.")
    }
    return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) GetOrdersByProduct(productID int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("product_id = ?", productID).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateOrderStatus(orderID int, status string) error {
	var order models.Order
	if err := r.db.First(&order, orderID).Error; err != nil {
		return err
	}
	order.Status = status
	return r.db.Save(&order).Error
}

