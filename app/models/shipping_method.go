package models

import (
	"time"
)

type ShippingMethod struct {
	ID            int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string     `json:"name" gorm:"not null"`
	Cost          float64    `json:"cost" gorm:"type:decimal(10,2);not null"`
	EstimatedDays *int       `json:"estimated_days,omitempty"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Shipments []Shipment `json:"shipments,omitempty" gorm:"foreignKey:ShippingMethodID"`
}

func (ShippingMethod) TableName() string {
	return "shipping_methods"
}
