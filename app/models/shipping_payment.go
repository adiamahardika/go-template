package models

import (
	"time"
)

type ShippingMethodFilter struct {
	Name   string `json:"name" query:"name"`
	Active bool   `json:"active" query:"active"`
	Page   int    `json:"page" query:"page"`
	Limit  int    `json:"limit" query:"limit"`
}

type PaymentMethodRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type PaymentMethodResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type PaymentMethodFilter struct {
	Name   string `json:"name" query:"name"`
	Active bool   `json:"active" query:"active"`
	Page   int    `json:"page" query:"page"`
	Limit  int    `json:"limit" query:"limit"`
}
