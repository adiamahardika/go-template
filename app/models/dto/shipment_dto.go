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

type AdminGetShipmentsRequest struct {
	Page             int     `json:"page" query:"page" validate:"min=1"`
	PageSize         int     `json:"page_size" query:"page_size" validate:"min=1,max=100"`
	Status           *string `json:"status" query:"status"`
	OrderID          *int    `json:"order_id" query:"order_id"`
	TrackingNumber   *string `json:"tracking_number" query:"tracking_number"`
	ShippingMethodID *int    `json:"shipping_method_id" query:"shipping_method_id"`
	DateFrom         *string `json:"date_from" query:"date_from"` 
	DateTo           *string `json:"date_to" query:"date_to"`     
}

type AdminShipmentListResponse struct {
	ID               int        `json:"id"`
	OrderID          int        `json:"order_id"`
	CustomerEmail    string     `json:"customer_email"`
	TotalAmount      float64    `json:"total_amount"`
	ShippingMethodID int        `json:"shipping_method_id"`
	TrackingNumber   *string    `json:"tracking_number"`
	Status           string     `json:"status"`
	ShippedAt        *time.Time `json:"shipped_at"`
	DeliveredAt      *time.Time `json:"delivered_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

type AdminShipmentDetailResponse struct {
	ID               int        `json:"id"`
	OrderID          int        `json:"order_id"`
	ShippingMethodID int        `json:"shipping_method_id"`
	TrackingNumber   *string    `json:"tracking_number"`
	Status           string     `json:"status"`
	ShippedAt        *time.Time `json:"shipped_at"`
	DeliveredAt      *time.Time `json:"delivered_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	Order struct {
		ID          int       `json:"id"`
		Status      string    `json:"status"`
		TotalAmount float64   `json:"total_amount"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"order"`

	ShippingMethod struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"shipping_method"`
}

type AdminUpdateShipmentRequest struct {
	TrackingNumber *string    `json:"tracking_number,omitempty"`
	Status         *string    `json:"status,omitempty" validate:"omitempty,oneof=pending,processing,shipped,delivered,cancelled"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`
	DeliveredAt    *time.Time `json:"delivered_at,omitempty"`
}

type AdminCreateShipmentRequest struct {
	OrderID          int    `json:"order_id" validate:"required"`
	ShippingMethodID int    `json:"shipping_method_id" validate:"required"`
	TrackingNumber   string `json:"tracking_number" validate:"required"`
	Status           string `json:"status" validate:"required,oneof=pending,processing,shipped,delivered,cancelled"`
}