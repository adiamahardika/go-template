package controllers

import (
	"net/http"
	"strconv"
	"monitoring-service/app/usecases"
	"github.com/labstack/echo/v4"
)

// Helper untuk dereference pointer
func derefString(s *string) string {
	if s == nil { return "" }
	return *s
}
func derefFloat(f *float64) float64 {
	if f == nil { return 0 }
	return *f
}

type OrderListItemResponse struct {
	OrderID        int      `json:"order_id"`
	TotalAmount    float64  `json:"total_amount"`
	Status         string   `json:"status"`
	CreatedAt      string   `json:"created_at"`
	Coupon         *string  `json:"coupon,omitempty"`
	PaymentStatus  string   `json:"payment_status"`
	ShipmentStatus string   `json:"shipment_status"`
}

type OrderControllerInterface interface {
	GetOrderHistory(ctx echo.Context) error
	GetOrderDetail(ctx echo.Context) error
}

type orderController struct {
	Usecase usecases.OrderUsecaseInterface
}

func NewOrderController(uc usecases.OrderUsecaseInterface) OrderControllerInterface {
	return &orderController{Usecase: uc}
}

func (c *orderController) GetOrderHistory(ctx echo.Context) error {
	userIDRaw := ctx.Get("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id in context"})
	}
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	if page < 1 { page = 1 }
	pageSize, _ := strconv.Atoi(ctx.QueryParam("page_size"))
	if pageSize < 1 { pageSize = 10 }

	orders, total, err := c.Usecase.GetOrderHistory(ctx.Request().Context(), uint(userID), page, pageSize)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var resp []OrderListItemResponse
	for _, o := range orders {
		paymentStatus := ""
		if len(o.Payments) > 0 && o.Payments[0].Status != nil {
			paymentStatus = *o.Payments[0].Status
		}
		shipmentStatus := ""
		if len(o.Shipments) > 0 && o.Shipments[0].Status != nil {
			shipmentStatus = *o.Shipments[0].Status
		}

		var couponCode *string
		if o.Coupon != nil {
			code := o.Coupon.Code
			couponCode = &code
		}

		resp = append(resp, OrderListItemResponse{
			OrderID:        o.ID,
			TotalAmount:    derefFloat(o.TotalAmount),
			Status:         derefString(o.Status),
			CreatedAt:      o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Coupon:         couponCode,
			PaymentStatus:  paymentStatus,
			ShipmentStatus: shipmentStatus,
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": resp,
		"pagination": map[string]interface{}{
			"page": page,
			"page_size": pageSize,
			"total": total,
		},
	})
}

func (c *orderController) GetOrderDetail(ctx echo.Context) error {
	userIDRaw := ctx.Get("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid user id in context"})
	}
	orderID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}
	order, err := c.Usecase.GetOrderDetail(ctx.Request().Context(), uint(userID), uint(orderID))
	if err != nil {
		// don't reveal whether record exists; if not owner or not found -> 403
		return ctx.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	type OrderItemResp struct {
		ProductName string  `json:"product_name"`
		Quantity    int     `json:"quantity"`
		Price       float64 `json:"price"`
		Subtotal    float64 `json:"subtotal"`
	}
	type PaymentResp struct {
		PaymentMethod string  `json:"payment_method"`
		Amount        float64 `json:"amount"`
		Status        string  `json:"status"`
		PaidAt        string  `json:"paid_at,omitempty"`
	}
	type ShipmentResp struct {
		ShippingMethod string `json:"shipping_method"`
		TrackingNumber string `json:"tracking_number,omitempty"`
		Status         string `json:"status"`
		ShippedAt      string `json:"shipped_at,omitempty"`
		DeliveredAt    string `json:"delivered_at,omitempty"`
	}
	type OrderDetailResp struct {
		ID          int             `json:"id"`
		TotalAmount float64         `json:"total_amount"`
		Status      string          `json:"status"`
		CreatedAt   string          `json:"created_at"`
		UpdatedAt   string          `json:"updated_at"`
		Coupon      interface{}     `json:"coupon,omitempty"`
		OrderItems  []OrderItemResp `json:"order_items"`
		Payment     *PaymentResp    `json:"payment,omitempty"`
		Shipment    *ShipmentResp   `json:"shipment,omitempty"`
	}

	var items []OrderItemResp
	for _, it := range order.OrderItems {
		pname := ""
		if it.Product != nil {
			pname = it.Product.Name
		}
		subtotal := float64(it.Quantity) * it.Price
		items = append(items, OrderItemResp{
			ProductName: pname,
			Quantity:    it.Quantity,
			Price:       it.Price,
			Subtotal:    subtotal,
		})
	}

	var coupon interface{}
	if order.Coupon != nil {
		coupon = map[string]interface{}{
			"id": order.Coupon.ID,
			"code": order.Coupon.Code,
			"discount_percent": order.Coupon.DiscountPercent,
			"max_discount": order.Coupon.MaxDiscount,
			"expired_at": order.Coupon.ExpiredAt,
		}
	}

	var paymentResp *PaymentResp
	if len(order.Payments) > 0 {
		p := order.Payments[0]
		method := ""
		if p.PaymentMethod != nil {
			method = p.PaymentMethod.Name
		}
		paidAt := ""
		if p.PaidAt != nil { paidAt = p.PaidAt.Format("2006-01-02T15:04:05Z07:00") }
		amount := 0.0
		if p.Amount != nil { amount = *p.Amount }
		status := ""
		if p.Status != nil { status = *p.Status }
		paymentResp = &PaymentResp{PaymentMethod: method, Amount: amount, Status: status, PaidAt: paidAt}
	}

	var shipmentResp *ShipmentResp
	if len(order.Shipments) > 0 {
		s := order.Shipments[0]
		method := ""
		if s.ShippingMethod != nil { method = s.ShippingMethod.Name }
		shippedAt := ""
		if s.ShippedAt != nil { shippedAt = s.ShippedAt.Format("2006-01-02T15:04:05Z07:00") }
		deliveredAt := ""
		if s.DeliveredAt != nil { deliveredAt = s.DeliveredAt.Format("2006-01-02T15:04:05Z07:00") }
		tracking := ""
		if s.TrackingNumber != nil { tracking = *s.TrackingNumber }
		status := ""
		if s.Status != nil { status = *s.Status }
		shipmentResp = &ShipmentResp{ShippingMethod: method, TrackingNumber: tracking, Status: status, ShippedAt: shippedAt, DeliveredAt: deliveredAt}
	}

	resp := OrderDetailResp{
		ID: order.ID,
		TotalAmount: derefFloat(order.TotalAmount),
		Status: derefString(order.Status),
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Coupon: coupon,
		OrderItems: items,
		Payment: paymentResp,
		Shipment: shipmentResp,
	}

	return ctx.JSON(http.StatusOK, resp)
}
