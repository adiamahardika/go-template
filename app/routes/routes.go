package routes

import (
	"monitoring-service/app/controllers"
	appmw "monitoring-service/middleware"
	"monitoring-service/pkg/config"

	"github.com/labstack/echo/v4"
)

func ConfigureRouter(e *echo.Echo, controller *controllers.Main, cfg *config.Config) {
	// API v1 group
	v1 := e.Group("/api/v1")

	// Auth (public)
	authGroup := v1.Group("/auth")
	authGroup.POST("/login", controller.Auth.Login)
	authGroup.POST("/register", controller.Auth.Register)

	// Users (contoh: public read-only, sesuaikan kebijakan)
	userGroup := v1.Group("/users")
	userGroup.GET("", controller.User.GetAllUsers)
	userGroup.GET("/:id", controller.User.GetUserByID)

	// Cart (shopper only)
	cartGroup := v1.Group("/cart", appmw.JWTRequireRoles(cfg, "shopper"))
	cartGroup.GET("", controller.Cart.GetCart)
	cartGroup.GET("/items", controller.Cart.GetCartItems)
	cartGroup.POST("/items", controller.Cart.AddCartItem)
	cartGroup.DELETE("/items", controller.Cart.RemoveCartItem)
	cartGroup.PUT("/items/:id", controller.Cart.UpdateCartItem)

	// New coupon routes
	cartGroup.POST("/coupon", controller.Cart.ApplyCoupon)
	cartGroup.DELETE("/coupon", controller.Cart.RemoveCoupon)
	cartGroup.GET("/summary", controller.Cart.GetCartSummary)
}
