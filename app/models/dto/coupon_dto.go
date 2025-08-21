package dto

import "time"

// CouponRequest is the DTO for creating and updating coupons
type CouponRequest struct {
	Code            string     `json:"code" binding:"required"`
	DiscountPercent float64    `json:"discount_percent" binding:"required,gte=0,lte=100"`
	MaxDiscount     *float64   `json:"max_discount,omitempty" binding:"omitempty,gte=0"`
	ExpiredAt       *time.Time `json:"expired_at,omitempty"`
}

// CouponResponse is the DTO for coupon responses
type CouponResponse struct {
	ID              int        `json:"id"`
	Code            string     `json:"code"`
	DiscountPercent float64    `json:"discount_percent"`
	MaxDiscount     *float64   `json:"max_discount,omitempty"`
	ExpiredAt       *time.Time `json:"expired_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// GetCouponsRequest represents the request for getting coupons with filters and pagination
type GetCouponsRequest struct {
	Page   int    `json:"page" query:"page"`
	Limit  int    `json:"limit" query:"limit"`
	Code   string `json:"code" query:"code"`
	Active bool   `json:"active" query:"active"`
}