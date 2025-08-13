package controllers

import (
	"monitoring-service/app/helpers"
	"monitoring-service/app/models"
	"monitoring-service/app/models/dto"
	customerror "monitoring-service/pkg/customerror"
	"net/http"
	"strconv"

	"github.com/ezartsh/inrequest"
	"github.com/labstack/echo/v4"
)

type shipmentController controller

type ShipmentControllerInterface interface {
	GetAllShipments(c echo.Context) error
	GetShipmentByID(c echo.Context) error
	UpdateShipment(c echo.Context) error
	CreateShipment(c echo.Context) error
}

func (ctrl *shipmentController) GetAllShipments(c echo.Context) error {
	var (
		request    dto.AdminGetShipmentsRequest
		shipments  []dto.AdminShipmentListResponse
		pagination models.Pagination
		err        error
	)

	queryReq := inrequest.Query(c.Request())
	if err := queryReq.ToBind(&request); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	shipments, pagination, err = ctrl.Options.UseCases.Shipment.GetAllShipmentsWithFilters(request)
	if err != nil {
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Shipments retrieved successfully"}, shipments, &pagination)
}

func (ctrl *shipmentController) GetShipmentByID(c echo.Context) error {
	var (
		shipment *dto.AdminShipmentDetailResponse
		err      error
	)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"Invalid shipment ID format"}, nil, nil)
	}

	shipment, err = ctrl.Options.UseCases.Shipment.GetShipmentByID(id)
	if err != nil {
		if err.Error() == "record not found" {
			return helpers.StandardResponse(c, http.StatusNotFound, []string{"Shipment not found"}, nil, nil)
		}
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Shipment retrieved successfully"}, shipment, nil)
}

func (ctrl *shipmentController) UpdateShipment(c echo.Context) error {
	var (
		reqBody  dto.AdminUpdateShipmentRequest
		shipment *dto.AdminShipmentDetailResponse
		err      error
	)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"Invalid shipment ID format"}, nil, nil)
	}

	req, err := inrequest.Json(c.Request())
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	if err := req.ToBind(&reqBody); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	shipment, err = ctrl.Options.UseCases.Shipment.UpdateShipment(c.Request().Context(), id, reqBody)
	if err != nil {
		if err.Error() == "shipment not found" {
			return helpers.StandardResponse(c, http.StatusNotFound, []string{err.Error()}, nil, nil)
		}
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Shipment updated successfully"}, shipment, nil)
}

func (ctrl *shipmentController) CreateShipment(c echo.Context) error {
	var (
		reqBody  dto.AdminCreateShipmentRequest
		shipment *dto.AdminShipmentListResponse
		err      error
	)

	req, err := inrequest.Json(c.Request())
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	if err := req.ToBind(&reqBody); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	shipment, err = ctrl.Options.UseCases.Shipment.CreateShipment(c.Request().Context(), reqBody)
	if err != nil {
		errMsg := err.Error()
		switch errMsg {
		case "order not found":
			return helpers.StandardResponse(c, http.StatusNotFound, []string{errMsg}, nil, nil)
		case "shipment already exists for this order":
			return helpers.StandardResponse(c, http.StatusConflict, []string{errMsg}, nil, nil)
		case "order is not eligible for shipment creation (must be paid or processing)":
			return helpers.StandardResponse(c, http.StatusBadRequest, []string{errMsg}, nil, nil)
		default:
			return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{errMsg}, nil, nil)
		}
	}

	return helpers.StandardResponse(c, http.StatusCreated, []string{"Shipment created successfully"}, shipment, nil)
}