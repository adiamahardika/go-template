package controllers

import (
	"encoding/json"
	"monitoring-service/app/helpers"
	"monitoring-service/app/models/dto"
	appmw "monitoring-service/middleware"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CartControllerInterface interface {
	GetCart(c echo.Context) error
	GetCartItems(c echo.Context) error
	AddCartItem(c echo.Context) error
	RemoveCartItem(c echo.Context) error
	UpdateCartItem(c echo.Context) error

	// New methods for coupon functionality
	ApplyCoupon(c echo.Context) error
	RemoveCoupon(c echo.Context) error
	GetCartSummary(c echo.Context) error
}

type cartController struct {
	Options Options
}

func (ctrl *cartController) GetCart(c echo.Context) error {
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	cartSummary, err := ctrl.Options.UseCases.Cart.CalculateCartTotal(c.Request().Context(), userID)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, nil, cartSummary, nil)
}

func (ctrl *cartController) GetCartItems(c echo.Context) error {
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	items, err := ctrl.Options.UseCases.Cart.GetCartItemsByUserID(c.Request().Context(), userID)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, nil, map[string]interface{}{
		"items": items,
	}, nil)
}

func (ctrl *cartController) AddCartItem(c echo.Context) error {
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	var body dto.AddCartItemRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"invalid json body"}, nil, nil)
	}
	if body.ProductID <= 0 {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"product_id is required"}, nil, nil)
	}
	if body.Quantity <= 0 {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"quantity must be >= 1"}, nil, nil)
	}

	item, message, err := ctrl.Options.UseCases.Cart.AddCartItem(c.Request().Context(), userID, body.ProductID, body.Quantity)
	if err != nil {
		switch err.Error() {
		case "product not found":
			return helpers.StandardResponse(c, http.StatusNotFound, []string{err.Error()}, nil, nil)
		case "product is deleted":
			return helpers.StandardResponse(c, http.StatusNotFound, []string{err.Error()}, nil, nil)
		case "quantity exceeds stock":
			return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
		default:
			return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
		}
	}

	data := dto.CartItemResponse{
		ID:        item.ID,
		CartID:    *item.CartID,
		ProductID: *item.ProductID,
		Quantity:  item.Quantity,
		Message:   message,
	}

	return helpers.StandardResponse(c, http.StatusCreated, []string{"Item added to cart"}, data, nil)
}

func (ctrl *cartController) RemoveCartItem(c echo.Context) error {
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	productIDStr := c.QueryParam("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"invalid product id"}, nil, nil)
	}

	if err := ctrl.Options.UseCases.Cart.RemoveCartItem(c.Request().Context(), userID, productID); err != nil {
		return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Item removed"}, nil, nil)
}

func (ctrl *cartController) UpdateCartItem(c echo.Context) error {
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	idStr := c.Param("id")
	itemID, err := strconv.Atoi(idStr)
	if err != nil || itemID <= 0 {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"invalid cart item id"}, nil, nil)
	}

	var body struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"invalid json body"}, nil, nil)
	}
	if body.Quantity < 1 {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"quantity must be >= 1"}, nil, nil)
	}

	item, err := ctrl.Options.UseCases.Cart.UpdateCartItemQuantity(c.Request().Context(), userID, itemID, body.Quantity)
	if err != nil {
		code := http.StatusInternalServerError
		switch err.Error() {
		case "cart not found", "cart item not found", "product not found", "product is deleted":
			code = http.StatusNotFound
		case "quantity exceeds stock", "invalid quantity":
			code = http.StatusBadRequest
		}
		return helpers.StandardResponse(c, code, []string{err.Error()}, nil, nil)
	}

	data := dto.CartItemResponse{
		ID:        item.ID,
		CartID:    *item.CartID,
		ProductID: *item.ProductID,
		Quantity:  item.Quantity,
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Item updated"}, data, nil)
}

// New methods for coupon functionality
func (ctrl *cartController) ApplyCoupon(c echo.Context) error {
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	var body dto.ApplyCouponRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"invalid json body"}, nil, nil)
	}
	if body.CouponCode == "" {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"coupon_code is required"}, nil, nil)
	}

	if err := ctrl.Options.UseCases.Cart.ApplyCoupon(c.Request().Context(), userID, body.CouponCode); err != nil {
		switch err.Error() {
		case "coupon not found":
			return helpers.StandardResponse(c, http.StatusNotFound, []string{err.Error()}, nil, nil)
		case "coupon is not valid or expired":
			return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
		default:
			return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
		}
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Coupon applied successfully"}, nil, nil)
}

func (ctrl *cartController) RemoveCoupon(c echo.Context) error {
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	if err := ctrl.Options.UseCases.Cart.RemoveCoupon(c.Request().Context(), userID); err != nil {
		if err.Error() == "cart not found" {
			return helpers.StandardResponse(c, http.StatusNotFound, []string{err.Error()}, nil, nil)
		}
		return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Coupon removed successfully"}, nil, nil)
}

func (ctrl *cartController) GetCartSummary(c echo.Context) error {
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	cartSummary, err := ctrl.Options.UseCases.Cart.CalculateCartTotal(c.Request().Context(), userID)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, nil, cartSummary, nil)
}
