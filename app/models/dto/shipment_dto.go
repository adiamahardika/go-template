package dto

import (
	"monitoring-service/app/models"
	"time"
)

type ShipmentResponse struct {
	ID               int                     `json:"id"`
	OrderID          *int                    `json:"order_id,omitempty"`
	ShippingMethodID *int                    `json:"shipping_method_id,omitempty"`
	TrackingNumber   *string                 `json:"tracking_number,omitempty"`
	ShippedAt        *time.Time              `json:"shipped_at,omitempty"`
	DeliveredAt      *time.Time              `json:"delivered_at,omitempty"`
	Status           *string                 `json:"status,omitempty" gorm:"size:50"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	Order            *models.Order           `json:"order,omitempty"`
	ShippingMethod   *models.ShippingMethod  `json:"shipping_method,omitempty"`
}
