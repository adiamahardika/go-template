package routes


import (
	"monitoring-service/app/controllers"
	"monitoring-service/app/middleware"
	"monitoring-service/pkg/config"

	"github.com/labstack/echo/v4"
)
func ConfigureRouter(e *echo.Echo, controller *controllers.Main, cfg *config.Config) {
	// API v1 group
	v1 := e.Group("/api/v1")

	authGroup := v1.Group("/auth")
	authGroup.POST("/login", controller.Auth.Login)
	authGroup.POST("/register", controller.Auth.Register)

	// User routes
	userGroup := v1.Group("/users")
	userGroup.GET("", controller.User.GetAllUsers)
	userGroup.GET("/:id", controller.User.GetUserByID)

// Shipment routes
	shipmentGroup := e.Group("/api/v1/shipments")
	shipmentGroup.Use(middleware.JWTAuthMiddleware(cfg))
	shipmentGroup.GET("", controller.Shipment.GetAllShipments)
	shipmentGroup.GET("/:id", controller.Shipment.GetShipmentByID)
	shipmentGroup.PUT("/:id", controller.Shipment.UpdateShipment)
	shipmentGroup.POST("", controller.Shipment.CreateShipment)
}
