// file: app/controllers/order_controller.go
package controllers

import (
	"monitoring-service/app/helpers"
	"monitoring-service/pkg/customerror"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type orderController controller

type OrderControllerInterface interface {
	CheckoutOrder(c echo.Context) error
	CancelOrder(c echo.Context) error
}

func (ctrl *orderController) CheckoutOrder(c echo.Context) error {
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid order ID"})
	}

	err = ctrl.Options.UseCases.Order.ProcessOrderCheckout(orderID)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Order checkout processed successfully"}, nil, nil)
}

func (ctrl *orderController) CancelOrder(c echo.Context) error {
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid order ID"})
	}

	err = ctrl.Options.UseCases.Order.CancelOrder(orderID)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Order cancelled successfully"}, nil, nil)
}