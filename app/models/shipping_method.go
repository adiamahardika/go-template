package models

import "time"

type ShippingMethod struct {
	ID            int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string     `json:"name" gorm:"not null"`
	Cost          float64    `json:"cost" gorm:"type:decimal(10,2);not null"`
	EstimatedDays int        `json:"estimated_days"`
	Description   *string    `json:"description,omitempty"`
	CreatedAt     time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type ShippingMethodRequest struct {
	Name          string  `json:"name" validate:"required"`
	Cost          float64 `json:"cost" validate:"required,min=0"`
	EstimatedDays int     `json:"estimated_days" validate:"required,min=0"`
	Description   *string `json:"description,omitempty"`
}

type ShippingMethodResponse struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Cost          float64   `json:"cost"`
	EstimatedDays int       `json:"estimated_days"`
	Description   *string   `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
