package dto

import (
	"monitoring-service/app/models" 
	"time"
)

type CreatePaymentRequest struct {
	OrderID         int     `json:"order_id" validate:"required"`
	PaymentMethodID int     `json:"payment_method_id" validate:"required"`
	Amount          float64 `json:"amount" validate:"required"`
}

type UpdatePaymentStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=paid failed"`
}

type PaymentResponse struct {
	ID              int        `json:"id"`
	OrderID         int        `json:"order_id"`
	PaymentMethodID int        `json:"payment_method_id"`
	Amount          float64    `json:"amount"`
	Status          string     `json:"status"`
	PaidAt          *time.Time `json:"paid_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func ToPaymentResponse(p *models.Payment) PaymentResponse {
	response := PaymentResponse{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		PaidAt:    p.PaidAt,
	}
	if p.OrderID != nil {
		response.OrderID = *p.OrderID
	}
	if p.PaymentMethodID != nil {
		response.PaymentMethodID = *p.PaymentMethodID
	}
	if p.Amount != nil {
		response.Amount = *p.Amount
	}
	if p.Status != nil {
		response.Status = *p.Status
	}
	return response
}
func ToPaymentResponses(payments []models.Payment) []PaymentResponse {
	responses := make([]PaymentResponse, len(payments))
	for i, p := range payments {
		responses[i] = ToPaymentResponse(&p)
	}
	return responses
}