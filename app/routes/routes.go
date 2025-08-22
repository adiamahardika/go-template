package routes

import (
	"monitoring-service/app/controllers"
	"monitoring-service/pkg/config"
	"monitoring-service/pkg/middleware"

	"github.com/labstack/echo/v4"
)

func ConfigureRouter(e *echo.Echo, controller *controllers.Main, cfg *config.Config) {
	// API v1 group
	v1 := e.Group("/api/v1")

	// Public routes
	public := v1.Group("/public")
	public.POST("/register", controller.User.Register)

	// User routes (protected)
	authGroup := v1.Group("/auth")
	authGroup.POST("/login", controller.Auth.Login)
	authGroup.POST("/register", controller.Auth.Register)

	// User routes
	// Public routes (no authentication needed)
	v1.POST("/login", controller.User.Login)

	// User routes (protected)
	userGroup := v1.Group("/users")
	userGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	userGroup.GET("", controller.User.GetAllUsers)
	userGroup.GET("/:id", controller.User.GetUserByID)

	productGroup := v1.Group("/product")
	productGroup.GET("/:id", controller.Product.GetProductByID)

	// Order routes (protected)
	orderGroup := v1.Group("/orders")
	orderGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	orderGroup.GET("", controller.Order.GetOrderHistory)
	orderGroup.GET("/:id", controller.Order.GetOrderDetail)

	// Admin routes (protected and admin role required)
	adminGroup := v1.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	adminGroup.Use(middleware.RoleMiddleware("admin"))
	// Tambahkan route admin di sini nanti
}
