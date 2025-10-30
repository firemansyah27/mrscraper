package models

import "time"

type Order struct {
    ID        int       `json:"id,omitempty" gorm:"primaryKey"`
    ProductID int       `json:"productId" binding:"required"`
    Quantity  int       `json:"quantity" binding:"required"`
    Total     float64   `json:"totalPrice"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"createdAt,omitempty" gorm:"autoCreateTime"`
}
