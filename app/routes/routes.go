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

	// Public routes (no authentication needed)
	v1.POST("/login", controller.User.Login)

	// User routes (protected)
	userGroup := v1.Group("/users")
	userGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	userGroup.GET("", controller.User.GetAllUsers)
	userGroup.GET("/:id", controller.User.GetUserByID)

	// Admin routes (protected and admin role required)
	adminGroup := v1.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	adminGroup.Use(middleware.RoleMiddleware("admin"))
	// Tambahkan route admin di sini nanti
}
