package models

import (
	"time"
)

type OrderItem struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID   *int      `json:"order_id,omitempty"`
	ProductID *int      `json:"product_id,omitempty"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"type:decimal(12,2);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Order   *Order   `json:"order,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
