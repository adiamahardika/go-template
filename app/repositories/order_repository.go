package repositories

import (
	"monitoring-service/app/models"

	"gorm.io/gorm"
)

type orderRepository repository

type OrderRepositoryInterface interface {
	GetOrderByID(orderID int) (*models.Order, error)
	GetOrderWithItems(orderID int) (*models.Order, error)
	UpdateOrderStatus(orderID int, status string) error
	GetOrderItems(orderID int) ([]models.OrderItem, error)
	ProcessCheckoutTransaction(orderID int, processFunc func(tx *gorm.DB) error) error
	CancelOrderTransaction(orderID int, processFunc func(tx *gorm.DB) error) error
}

func (r *orderRepository) GetOrderByID(orderID int) (*models.Order, error) {
	var order models.Order
	err := r.Options.Postgres.
		Where("id = ?", orderID).
		First(&order).Error
	return &order, err
}

func (r *orderRepository) GetOrderWithItems(orderID int) (*models.Order, error) {
	var order models.Order
	err := r.Options.Postgres.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("id = ?", orderID).
		First(&order).Error
	return &order, err
}

func (r *orderRepository) UpdateOrderStatus(orderID int, status string) error {
	return r.Options.Postgres.
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

func (r *orderRepository) GetOrderItems(orderID int) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.Options.Postgres.
		Where("order_id = ?", orderID).
		Find(&items).Error
	return items, err
}

func (r *orderRepository) ProcessCheckoutTransaction(orderID int, processFunc func(tx *gorm.DB) error) error {
	tx := r.Options.Postgres.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := processFunc(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *orderRepository) CancelOrderTransaction(orderID int, processFunc func(tx *gorm.DB) error) error {
	tx := r.Options.Postgres.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := processFunc(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}