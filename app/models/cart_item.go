package models

import (
	"time"
)

type CartItem struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	CartID    *int      `json:"cart_id,omitempty"`
	ProductID *int      `json:"product_id,omitempty"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Cart    *Cart    `json:"cart,omitempty" gorm:"foreignKey:CartID;references:ID"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID;references:ID"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
