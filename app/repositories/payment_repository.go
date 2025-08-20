package repositories

import (
	"context"
	"monitoring-service/app/models"
	
)

type paymentRepository repository

type PaymentRepositoryInterface interface {
	CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error)
	HasPaidPayment(ctx context.Context, orderID int) (bool, error)
	GetPendingPayment(ctx context.Context, orderID int) (*models.Payment, error)
	UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error
	GetUserPayments(ctx context.Context, userID, orderID int) ([]models.Payment, error)
	GetAllPayments(ctx context.Context) ([]models.Payment, error)
	UpdatePayment(ctx context.Context, payment *models.Payment) error
	GetPaymentByID(ctx context.Context, paymentID int) (*models.Payment, error)
}

func (r *paymentRepository) CreatePayment(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	if err := r.Options.Postgres.Create(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) HasPaidPayment(ctx context.Context, orderID int) (bool, error) {
	var count int64
	err := r.Options.Postgres.Model(&models.Payment{}).Where("order_id = ? AND status = ?", orderID, "paid").Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *paymentRepository) GetPendingPayment(ctx context.Context, orderID int) (*models.Payment, error) {
	var payment models.Payment
	err := r.Options.Postgres.Where("order_id = ? AND status = ?", orderID, "pending").First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePaymentStatus(ctx context.Context, paymentID int, status string) error {
	return r.Options.Postgres.Model(&models.Payment{}).Where("id = ?", paymentID).Update("status", status).Error
}

func (r *paymentRepository) GetUserPayments(ctx context.Context, userID, orderID int) ([]models.Payment, error) {
	var payments []models.Payment
	query := r.Options.Postgres.Joins("JOIN orders ON orders.id = payments.order_id").
		Where("orders.user_id = ?", userID)

	if orderID != 0 {
		query = query.Where("payments.order_id = ?", orderID)
	}

	err := query.Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) GetAllPayments(ctx context.Context) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.Options.Postgres.Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	return r.Options.Postgres.Save(payment).Error
}


func (r *paymentRepository) GetPaymentByID(ctx context.Context, paymentID int) (*models.Payment, error) {
	var payment models.Payment
	err := r.Options.Postgres.First(&payment, paymentID).Error
	return &payment, err
}