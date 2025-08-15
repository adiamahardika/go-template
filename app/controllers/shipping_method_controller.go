package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"monitoring-service/app/usecases"
)

type ShippingMethodControllerInterface interface {
	GetAll(c echo.Context) error
	GetByID(c echo.Context) error
	GetCartQuote(c echo.Context) error
}

type shippingMethodController struct {
	uc usecases.ShippingMethodUsecaseInterface
}

func NewShippingMethodController(uc usecases.ShippingMethodUsecaseInterface) ShippingMethodControllerInterface {
	return &shippingMethodController{uc}
}

func (ctrl *shippingMethodController) GetAll(c echo.Context) error {
	methods, err := ctrl.uc.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, methods)
}

func (ctrl *shippingMethodController) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	method, err := ctrl.uc.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, method)
}

func (ctrl *shippingMethodController) GetCartQuote(c echo.Context) error {
	userIDf := c.Get("user_id")
	if userIDf == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	userID, ok := userIDf.(int)
	if !ok {
		userIDFloat, ok := userIDf.(float64)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user_id in token"})
		}
		userID = int(userIDFloat)
	}
	shippingMethodID, err := strconv.Atoi(c.QueryParam("shipping_method_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid shipping_method_id"})
	}
	quote, err := ctrl.uc.GetCartQuote(userID, shippingMethodID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, quote)
}
