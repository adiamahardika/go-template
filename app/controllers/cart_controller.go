package controllers

import (
	"encoding/json"
	"monitoring-service/app/helpers"
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
}

type cartController struct {
	Options Options
}

func (ctrl *cartController) GetCart(c echo.Context) error {
    userID, err := appmw.CurrentUserID(c)
    if err != nil {
        return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
    }

    cart, items, err := ctrl.Options.UseCases.Cart.GetCart(c.Request().Context(), userID)
    if err != nil {
        return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
    }

    // Bentuk ulang items: product_id, name, price, quantity, subtotal
    viewItems := make([]map[string]interface{}, 0, len(items))
    var total float64

    for _, it := range items {
        pid := 0
        if it.ProductID != nil {
            pid = *it.ProductID
        }

        name := ""
        var price float64
        if it.Product != nil {
            name = it.Product.Name
            price = it.Product.Price
        }

        subtotal := float64(it.Quantity) * price
        total += subtotal

        viewItems = append(viewItems, map[string]interface{}{
            "product_id": pid,
            "name":       name,
            "price":      price,
            "quantity":   it.Quantity,
            "subtotal":   subtotal,
        })
    }

    // Data akhir sesuai spes
    data := map[string]interface{}{
        "id":         cart.ID,
        "created_at": cart.CreatedAt,
        "updated_at": cart.UpdatedAt,
        "items":      viewItems,
        "total":      total,
    }

    return helpers.StandardResponse(c, http.StatusOK, nil, data, nil)
}


func (ctrl *cartController) GetCartItems(c echo.Context) error {
	// Ambil user dari JWT (bukan query param)
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
	// Ambil user dari JWT (bukan query param)
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	// Body request: product_id & quantity
	var body struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
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

	data := map[string]any{
		"id":         item.ID,
		"cart_id":    item.CartID,
		"product_id": item.ProductID,
		"quantity":   item.Quantity,
	}
	if message != "" {
		data["message"] = message
	}

	return helpers.StandardResponse(c, http.StatusCreated, []string{"Item added to cart"}, data, nil)
}

func (ctrl *cartController) RemoveCartItem(c echo.Context) error {
	// Ambil user dari JWT (bukan query param)
	userID, err := appmw.CurrentUserID(c)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusUnauthorized, []string{err.Error()}, nil, nil)
	}

	// product_id tetap dari query (sesuai requirement)
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

	// build response
	cartID := 0
	if item.CartID != nil { cartID = *item.CartID }
	productID := 0
	if item.ProductID != nil { productID = *item.ProductID }

	data := map[string]any{
		"id":         item.ID,
		"cart_id":    cartID,
		"product_id": productID,
		"quantity":   item.Quantity,
	}
	return helpers.StandardResponse(c, http.StatusOK, []string{"Item updated"}, data, nil)
}
