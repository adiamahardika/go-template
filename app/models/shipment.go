package models

import (
	"time"
)

type Shipment struct {
	ID               int        `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID          *int       `json:"order_id,omitempty"`
	ShippingMethodID *int       `json:"shipping_method_id,omitempty"`
	TrackingNumber   *string    `json:"tracking_number,omitempty"`
	ShippedAt        *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt      *time.Time `json:"delivered_at,omitempty"`
	Status           *string    `json:"status,omitempty" gorm:"size:50"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Order          *Order          `json:"order,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	ShippingMethod *ShippingMethod `json:"shipping_method,omitempty" gorm:"foreignKey:ShippingMethodID;references:ID"`
}

func (Shipment) TableName() string {
	return "shipments"
}
