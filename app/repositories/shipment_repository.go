package repositories

import (
	"monitoring-service/app/models"
)

type shipmentRepository repository

type ShipmentRepositoryInterface interface {

	GetOrderByID(orderID int) (*models.Order, error)
	GetShipmentByOrderID(orderID int) (*models.Shipment, error)
}

func (r *shipmentRepository) GetOrderByID(orderID int) (*models.Order, error) {
	var order models.Order
	err := r.Options.Postgres.First(&order, orderID).Error
	return &order, err
}

func (r *shipmentRepository) GetShipmentByOrderID(orderID int) (*models.Shipment, error) {
	var shipment models.Shipment
	err := r.Options.Postgres.Preload("ShippingMethod").Where("order_id = ?", orderID).First(&shipment).Error
	return &shipment, err
}