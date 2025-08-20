package usecases

import (
	"context"
	"errors"
	"monitoring-service/app/models"
	"monitoring-service/app/models/dto"
	"time"

	gorm "gorm.io/gorm"
)

type paymentUsecase usecase

type paymentUsecaseInterface interface {
	CreatePayment(ctx context.Context, req dto.CreatePaymentRequest, userID int) (*dto.PaymentResponse, error)
	UpdatePaymentStatus(ctx context.Context, paymentID int, req dto.UpdatePaymentStatusRequest) (*dto.PaymentResponse, error)
	GetUserPayments(ctx context.Context, userID, orderID int) ([]dto.PaymentResponse, error)
	GetAllPayments(ctx context.Context) ([]dto.PaymentResponse, error)
}

func (u *paymentUsecase) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest, userID int) (*dto.PaymentResponse, error) {
	isOwner, err := u.Options.Repository.Order.IsOrderOwner(userID, req.OrderID)
	if err != nil {
		return nil, errors.New("failed to verify order ownership")
	}
	if !isOwner {
		return nil, errors.New("you are not authorized to pay for this order")
	}

	hasPaid, err := u.Options.Repository.Payment.HasPaidPayment(ctx, req.OrderID)
	if err != nil {
		return nil, errors.New("failed to check existing payments")
	}
	if hasPaid {
		return nil, errors.New("this order has already been paid")
	}

	existingPending, _ := u.Options.Repository.Payment.GetPendingPayment(ctx, req.OrderID)
	if existingPending != nil {
		return nil, errors.New("a pending payment for this order already exists")
	}

	status := "pending"
	amount := req.Amount

	newPayment := &models.Payment{
		OrderID:         &req.OrderID,
		PaymentMethodID: &req.PaymentMethodID,
		Amount:          &amount,
		Status:          &status,
	}

	createdPayment, err := u.Options.Repository.Payment.CreatePayment(ctx, newPayment)
	if err != nil {
		return nil, errors.New("failed to record payment")
	}

	response := dto.ToPaymentResponse(createdPayment)
	return &response, nil
}


func (u *paymentUsecase) UpdatePaymentStatus(ctx context.Context, paymentID int, req dto.UpdatePaymentStatusRequest) (*dto.PaymentResponse, error) {
	payment, err := u.Options.Repository.Payment.GetPaymentByID(ctx, paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, errors.New("failed to retrieve payment")
	}
	if *payment.Status != "pending" {
		return nil, errors.New("only pending payments can be updated")
	}

	now := time.Now()
	payment.Status = &req.Status
	if req.Status == "paid" {
		payment.PaidAt = &now
	}

	if err := u.Options.Repository.Payment.UpdatePayment(ctx, payment); err != nil {
		return nil, errors.New("failed to update payment status")
	}

	if req.Status == "paid" && payment.OrderID != nil {
		_ = u.Options.Repository.Order.UpdateOrderStatus(*payment.OrderID, "paid")
	}

	response := dto.ToPaymentResponse(payment)
	return &response, nil
}


func (u *paymentUsecase) GetUserPayments(ctx context.Context, userID, orderID int) ([]dto.PaymentResponse, error) {
	payments, err := u.Options.Repository.Payment.GetUserPayments(ctx, userID, orderID)
	if err != nil {
		return nil, errors.New("failed to retrieve payments")
	}

	responses := dto.ToPaymentResponses(payments)
	
	return responses, nil
}

func (u *paymentUsecase) GetAllPayments(ctx context.Context) ([]dto.PaymentResponse, error) {
	payments, err := u.Options.Repository.Payment.GetAllPayments(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve all payments")
	}
	responses := dto.ToPaymentResponses(payments)

	return responses, nil
}