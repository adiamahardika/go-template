package controllers

import (
	"net/http"
	"strconv"

	"monitoring-service/app/helpers"
	"monitoring-service/app/models"
	"monitoring-service/pkg/customerror"

	"github.com/ezartsh/validet"
	"github.com/labstack/echo/v4"
)

type ShippingPaymentControllerInterface interface {
	CreateShippingMethod(c echo.Context) error
	GetShippingMethods(c echo.Context) error
	GetShippingMethodByID(c echo.Context) error
	UpdateShippingMethod(c echo.Context) error
	DeleteShippingMethod(c echo.Context) error

	CreatePaymentMethod(c echo.Context) error
	GetPaymentMethods(c echo.Context) error
	GetPaymentMethodByID(c echo.Context) error
	UpdatePaymentMethod(c echo.Context) error
	DeletePaymentMethod(c echo.Context) error
}

type shippingPaymentController struct {
	*controller
}

func newShippingPaymentController(ctrl *controller) *shippingPaymentController {
	return &shippingPaymentController{controller: ctrl}
}

func (ctrl *shippingPaymentController) CreateShippingMethod(c echo.Context) error {
	var request models.ShippingMethodRequest
	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{err.Error()})
	}

	// Validate request
	mapReq := make(map[string]any)
	mapReq["name"] = request.Name
	mapReq["cost"] = request.Cost
	mapReq["estimated_days"] = request.EstimatedDays

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"name":           validet.String{Required: true},
			"cost":           validet.Numeric[float64]{Required: true, Min: 0},
			"estimated_days": validet.Numeric[float64]{Required: true, Min: 0},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	if len(errorBags.Errors) > 0 {
		var errorMessages []string
		for _, errs := range errorBags.Errors {
			for _, msg := range errs {
				errorMessages = append(errorMessages, msg)
			}
		}
		return helpers.StandardResponse(c, http.StatusBadRequest, errorMessages, nil, nil)
	}

	method, err := ctrl.Options.UseCases.ShippingPayment.CreateShippingMethod(request)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusCreated, []string{"Shipping method created successfully"}, method, nil)
}

func (ctrl *shippingPaymentController) GetShippingMethods(c echo.Context) error {
	var filter models.ShippingMethodFilter
	if err := c.Bind(&filter); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{err.Error()})
	}

	methods, pagination, err := ctrl.Options.UseCases.ShippingPayment.GetShippingMethods(filter)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Shipping methods retrieved successfully"}, methods, &pagination)
}

func (ctrl *shippingPaymentController) GetShippingMethodByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid ID format"})
	}

	method, err := ctrl.Options.UseCases.ShippingPayment.GetShippingMethodByID(id)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Shipping method retrieved successfully"}, method, nil)
}

func (ctrl *shippingPaymentController) UpdateShippingMethod(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid ID format"})
	}

	var request models.ShippingMethodRequest
	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{err.Error()})
	}

	// Validate request
	mapReq := make(map[string]any)
	mapReq["name"] = request.Name
	mapReq["cost"] = request.Cost
	mapReq["estimated_days"] = request.EstimatedDays

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"name":           validet.String{Required: true},
			"cost":           validet.Numeric[float64]{Required: true, Min: 0},
			"estimated_days": validet.Numeric[float64]{Required: true, Min: 0},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	if len(errorBags.Errors) > 0 {
		var errorMessages []string
		for _, errs := range errorBags.Errors {
			for _, msg := range errs {
				errorMessages = append(errorMessages, msg)
			}
		}
		return helpers.StandardResponse(c, http.StatusBadRequest, errorMessages, nil, nil)
	}

	method, err := ctrl.Options.UseCases.ShippingPayment.UpdateShippingMethod(id, request)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Shipping method updated successfully"}, method, nil)
}

func (ctrl *shippingPaymentController) DeleteShippingMethod(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid ID format"})
	}

	err = ctrl.Options.UseCases.ShippingPayment.DeleteShippingMethod(id)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.Response(c, http.StatusOK, []string{"Shipping method deleted successfully"})
}

func (ctrl *shippingPaymentController) CreatePaymentMethod(c echo.Context) error {
	var request models.PaymentMethodRequest
	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{err.Error()})
	}

	// Validate request
	mapReq := make(map[string]any)
	mapReq["name"] = request.Name

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"name": validet.String{Required: true},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	if len(errorBags.Errors) > 0 {
		var errorMessages []string
		for _, errs := range errorBags.Errors {
			for _, msg := range errs {
				errorMessages = append(errorMessages, msg)
			}
		}
		return helpers.StandardResponse(c, http.StatusBadRequest, errorMessages, nil, nil)
	}

	method, err := ctrl.Options.UseCases.ShippingPayment.CreatePaymentMethod(request)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusCreated, []string{"Payment method created successfully"}, method, nil)
}

func (ctrl *shippingPaymentController) GetPaymentMethods(c echo.Context) error {
	var filter models.PaymentMethodFilter
	if err := c.Bind(&filter); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{err.Error()})
	}

	methods, pagination, err := ctrl.Options.UseCases.ShippingPayment.GetPaymentMethods(filter)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Payment methods retrieved successfully"}, methods, &pagination)
}

func (ctrl *shippingPaymentController) GetPaymentMethodByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid ID format"})
	}

	method, err := ctrl.Options.UseCases.ShippingPayment.GetPaymentMethodByID(id)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Payment method retrieved successfully"}, method, nil)
}

func (ctrl *shippingPaymentController) UpdatePaymentMethod(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid ID format"})
	}

	var request models.PaymentMethodRequest
	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{err.Error()})
	}

	// Validate request
	mapReq := make(map[string]any)
	mapReq["name"] = request.Name

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"name": validet.String{Required: true},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	if len(errorBags.Errors) > 0 {
		var errorMessages []string
		for _, errs := range errorBags.Errors {
			for _, msg := range errs {
				errorMessages = append(errorMessages, msg)
			}
		}
		return helpers.StandardResponse(c, http.StatusBadRequest, errorMessages, nil, nil)
	}

	method, err := ctrl.Options.UseCases.ShippingPayment.UpdatePaymentMethod(id, request)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Payment method updated successfully"}, method, nil)
}

func (ctrl *shippingPaymentController) DeletePaymentMethod(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid ID format"})
	}

	err = ctrl.Options.UseCases.ShippingPayment.DeletePaymentMethod(id)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.Response(c, http.StatusOK, []string{"Payment method deleted successfully"})
}
