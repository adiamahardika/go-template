package models

// This file exports all models for easy importing
// Usage: import "your-project/app/models"

// User related models
type (
	UserModel     = User
	RoleModel     = Role
	UserRoleModel = UserRole
)

// Product related models
type (
	CategoryModel = Category
	ProductModel  = Product
)

// Cart related models
type (
	CartModel     = Cart
	CartItemModel = CartItem
)

// Order related models
type (
	OrderModel     = Order
	OrderItemModel = OrderItem
	CouponModel    = Coupon
)

// Payment related models
type (
	PaymentModel       = Payment
	PaymentMethodModel = PaymentMethod
)

// Shipping related models
type (
	ShipmentModel       = Shipment
	ShippingMethodModel = ShippingMethod
)

// Response models
type (
	ResponseModel             = Response
	BasicResponseModel        = BasicResponse
	ResponseWithPaginateModel = ResponseWithPaginate
	PaginationModel           = Pagination
)
