package repositories

import (
	"monitoring-service/app/models"
)

type orderRepository repository

type OrderRepositoryInterface interface {
	GetOrderByID(orderID int) (*models.Order, error)
	UpdateOrderStatus(orderID int, status string) error
	IsOrderOwner(userID, orderID int) (bool, error)
}

func (r *orderRepository) GetOrderByID(orderID int) (*models.Order, error) {
	var order models.Order
	err := r.Options.Postgres.First(&order, orderID).Error
	return &order, err
}

func (r *orderRepository) UpdateOrderStatus(orderID int, status string) error {
	return r.Options.Postgres.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *orderRepository) IsOrderOwner(userID, orderID int) (bool, error) {
	var count int64
	err := r.Options.Postgres.Model(&models.Order{}).Where("id = ? AND user_id = ?", orderID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}