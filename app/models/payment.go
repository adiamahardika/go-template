package models

import (
	"time"
)

type Payment struct {
	ID              int        `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID         *int       `json:"order_id,omitempty"`
	PaymentMethodID *int       `json:"payment_method_id,omitempty"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	Amount          *float64   `json:"amount,omitempty" gorm:"type:decimal(12,2)"`
	Status          *string    `json:"status,omitempty" gorm:"size:50"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Order         *Order         `json:"order,omitempty" gorm:"foreignKey:OrderID;references:ID"`
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty" gorm:"foreignKey:PaymentMethodID;references:ID"`
}

func (Payment) TableName() string {
	return "payments"
}
