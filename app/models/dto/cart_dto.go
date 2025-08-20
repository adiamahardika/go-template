package dto

import "time"

type AddCartItemRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

type CartItemResponse struct {
	ID        int    `json:"id"`
	CartID    int    `json:"cart_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Message   string `json:"message,omitempty"` // opsional: info jika quantity dikap karena stok
}

type CartItemView struct {
	ProductID int     `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Subtotal  float64 `json:"subtotal"`
}

type CartViewResponse struct {
	ID         int            `json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	CouponCode *string        `json:"coupon_code,omitempty"`
	Discount   float64        `json:"discount"`
	Items      []CartItemView `json:"items"`
	Subtotal   float64        `json:"subtotal"`
	Total      float64        `json:"total"`
}

type ApplyCouponRequest struct {
	CouponCode string `json:"coupon_code" validate:"required"`
}

type CartSummaryResponse struct {
	ID         int            `json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	CouponCode *string        `json:"coupon_code,omitempty"`
	Items      []CartItemView `json:"items"`
	Subtotal   float64        `json:"subtotal"`
	Discount   float64        `json:"discount"`
	Total      float64        `json:"total"`
}
