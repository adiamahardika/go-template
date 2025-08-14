package routes

import (
	"monitoring-service/app/controllers"
	"monitoring-service/pkg/middleware"

	"github.com/labstack/echo/v4"
)

func ConfigureRouter(e *echo.Echo, controller *controllers.Main) {
	// API v1 group
	v1 := e.Group("/api/v1")

	// Public routes
	public := v1.Group("/public")
	public.POST("/register", controller.User.Register)

	// User routes (protected)
	userGroup := v1.Group("/users")
	userGroup.GET("", controller.User.GetAllUsers)
	userGroup.GET("/:id", controller.User.GetUserByID)

	// Admin routes
	adminGroup := v1.Group("/admin")
	adminGroup.Use(middleware.AdminMiddleware())

	// Shipping methods
	shippingGroup := adminGroup.Group("/shipping-methods")
	shippingGroup.POST("", controller.ShippingPayment.CreateShippingMethod)
	shippingGroup.GET("", controller.ShippingPayment.GetShippingMethods)
	shippingGroup.GET("/:id", controller.ShippingPayment.GetShippingMethodByID)
	shippingGroup.PUT("/:id", controller.ShippingPayment.UpdateShippingMethod)
	shippingGroup.DELETE("/:id", controller.ShippingPayment.DeleteShippingMethod)

	// Payment methods
	paymentGroup := adminGroup.Group("/payment-methods")
	paymentGroup.POST("", controller.ShippingPayment.CreatePaymentMethod)
	paymentGroup.GET("", controller.ShippingPayment.GetPaymentMethods)
	paymentGroup.GET("/:id", controller.ShippingPayment.GetPaymentMethodByID)
	paymentGroup.PUT("/:id", controller.ShippingPayment.UpdatePaymentMethod)
	paymentGroup.DELETE("/:id", controller.ShippingPayment.DeletePaymentMethod)
}
