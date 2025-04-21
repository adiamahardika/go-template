package routes

import (
	"monitoring-service/app/controllers"

	"github.com/labstack/echo/v4"
)

func ConfigureRouter(e *echo.Echo, controller *controllers.Main) {
	v1 := e.Group("/v1")
	{
		todo := v1.Group("/todo")
		{
			todo.GET("", controller.Todo.GetAll)
		}

		priority := v1.Group("/priority")
		{
			priority.GET("", controller.Priority.GetAll)
		}
	}
}
