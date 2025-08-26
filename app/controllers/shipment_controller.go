package controllers

import (
	"monitoring-service/app/helpers"
	"monitoring-service/app/models/dto"
	customerror "monitoring-service/pkg/customerror"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type shipmentController controller

type ShipmentControllerInterface interface {
	GetShipmentByOrderID(c echo.Context) error
}

func (ctrl *shipmentController) GetShipmentByOrderID(c echo.Context) error {
	var (
		shipment *dto.ShipmentResponse
		err      error
	)

	orderID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"Invalid order ID format"}, nil, nil)
	}

	userID := c.Get("user_id").(int)

	shipment, err = ctrl.Options.UseCases.Shipment.GetShipmentByOrderID(orderID, userID)
	if err != nil {
		if err.Error() == "shipment not found" {
			return helpers.StandardResponse(c, http.StatusNotFound, []string{"Shipment not found for this order"}, nil, nil)
		}
		if err.Error() == "unauthorized" {
			return helpers.StandardResponse(c, http.StatusForbidden, []string{"You are not authorized to view this shipment"}, nil, nil)
		}
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Shipment retrieved successfully"}, shipment, nil)
}