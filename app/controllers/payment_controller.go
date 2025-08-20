package controllers

import (
	"monitoring-service/app/helpers"
	"monitoring-service/app/models/dto"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type paymentController controller

type PaymentControllerInterface interface {
	CreatePayment(c echo.Context) error
	UpdatePaymentStatus(c echo.Context) error
	GetUserPayments(c echo.Context) error
	GetAllPayments(c echo.Context) error
}

func (ctrl *paymentController) CreatePayment(c echo.Context) error {
	var req dto.CreatePaymentRequest

	if err := c.Bind(&req); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"Invalid request body: " + err.Error()}, nil, nil)
	}

	userID, ok := c.Get("user_id").(int)
	if !ok {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{"Invalid user ID in token"}, nil, nil)
	}

	paymentResponse, err := ctrl.Options.UseCases.Payment.CreatePayment(c.Request().Context(), req, userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "you are not authorized to pay for this order" {
			statusCode = http.StatusForbidden
		}
		if err.Error() == "this order has already been paid" || err.Error() == "a pending payment for this order already exists" {
			statusCode = http.StatusConflict
		}
		return helpers.StandardResponse(c, statusCode, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusCreated, []string{"Payment initiated successfully"}, paymentResponse, nil)
}

func (ctrl *paymentController) UpdatePaymentStatus(c echo.Context) error {
	paymentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"Invalid payment ID format"}, nil, nil)
	}

	var req dto.UpdatePaymentStatusRequest
	if err := c.Bind(&req); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"Invalid request body"}, nil, nil)
	}
    if req.Status != "paid" && req.Status != "failed" {
        return helpers.StandardResponse(c, http.StatusBadRequest, []string{"Invalid status value. Must be 'paid' or 'failed'."}, nil, nil)
    }

	paymentResponse, err := ctrl.Options.UseCases.Payment.UpdatePaymentStatus(c.Request().Context(), paymentID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "payment not found" {
			statusCode = http.StatusNotFound 
		}
		if err.Error() == "only pending payments can be updated" {
			statusCode = http.StatusConflict 
		}
		return helpers.StandardResponse(c, statusCode, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Payment status updated successfully"}, paymentResponse, nil)
}

func (ctrl *paymentController) GetUserPayments(c echo.Context) error {
	userID, ok := c.Get("user_id").(int)
	if !ok {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{"Invalid user ID in token"}, nil, nil)
	}

	orderIDStr := c.QueryParam("order_id")
	orderID, _ := strconv.Atoi(orderIDStr) 

	paymentResponses, err := ctrl.Options.UseCases.Payment.GetUserPayments(c.Request().Context(), userID, orderID)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Payments retrieved successfully"}, paymentResponses, nil)
}

func (ctrl *paymentController) GetAllPayments(c echo.Context) error {
	paymentResponses, err := ctrl.Options.UseCases.Payment.GetAllPayments(c.Request().Context())
	if err != nil {
		return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"All payments retrieved successfully"}, paymentResponses, nil)
}