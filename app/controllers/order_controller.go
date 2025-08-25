package controllers

import (
	"encoding/json"
	"monitoring-service/app/helpers"
	"monitoring-service/app/models/dto"
	appmw "monitoring-service/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

type OrderControllerInterface interface {
	Checkout(c echo.Context) error
}

type orderController struct {
	Options Options
}

func (ctrl *orderController) Checkout(c echo.Context) error {
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	var body dto.CheckoutRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"invalid json body"}, nil, nil)
	}

	// Validate shipping method ID
	if body.ShippingMethodID <= 0 {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"shipping_method_id is required"}, nil, nil)
	}

	order, err := ctrl.Options.UseCases.Order.Checkout(c.Request().Context(), userID, body.ShippingMethodID, body.CouponCode)
	if err != nil {
		// Log error for debugging
		c.Logger().Errorf("Checkout error: %v", err)

		switch err.Error() {
		case "cart not found", "cart is empty", "product not available", "shipping method not found", "coupon not found":
			return helpers.StandardResponse(c, http.StatusNotFound, []string{err.Error()}, nil, nil)
		case "insufficient stock for product:":
			return helpers.StandardResponse(c, http.StatusConflict, []string{err.Error()}, nil, nil)
		case "coupon is not valid or expired":
			return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
		default:
			if len(err.Error()) > 100 {
				// Jangan expose internal error details ke client
				return helpers.StandardResponse(c, http.StatusInternalServerError, []string{"internal server error"}, nil, nil)
			}
			return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
		}
	}

	return helpers.StandardResponse(c, http.StatusCreated, []string{"Order created successfully"}, order, nil)
}
