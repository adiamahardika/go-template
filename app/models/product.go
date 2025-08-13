package models

import (
	"time"
)

type Product struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"not null"`
	Description *string    `json:"description,omitempty"`
	Price       float64    `json:"price" gorm:"type:decimal(12,2);not null"`
	Stock       int        `json:"stock" gorm:"default:0"`
	CategoryID  *int       `json:"category_id,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Category   *Category   `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
	CartItems  []CartItem  `json:"cart_items,omitempty" gorm:"foreignKey:ProductID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string {
	return "products"
}
