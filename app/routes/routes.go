package routes

import (
	"monitoring-service/app/controllers"

	"github.com/labstack/echo/v4"
)

func ConfigureRouter(e *echo.Echo, controller *controllers.Main) {
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
	userGroup := v1.Group("/users")
	userGroup.GET("", controller.User.GetAllUsers)
	userGroup.GET("/:id", controller.User.GetUserByID)

	productGroup := v1.Group("/product")
	productGroup.GET("/:id", controller.Product.GetProductByID)
}
