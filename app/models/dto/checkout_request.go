package dto

type CheckoutRequest struct {
	ShippingMethodID int    `json:"shipping_method_id"`
	CouponCode       string `json:"coupon_code"`
}
