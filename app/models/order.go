package models

import (
	"time"
)

type Order struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID      *int      `json:"user_id,omitempty"`
	CouponID    *int      `json:"coupon_id,omitempty"`
	TotalAmount *float64  `json:"total_amount,omitempty" gorm:"type:decimal(12,2)"`
	Status      *string   `json:"status,omitempty" gorm:"size:50"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	User       *User       `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Coupon     *Coupon     `json:"coupon,omitempty" gorm:"foreignKey:CouponID;references:ID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
	Payments   []Payment   `json:"payments,omitempty" gorm:"foreignKey:OrderID"`
	Shipments  []Shipment  `json:"shipments,omitempty" gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "orders"
}
