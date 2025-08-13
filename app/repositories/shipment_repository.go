package repositories

import (
	"monitoring-service/app/models"
	"monitoring-service/app/models/dto"
	"time"
)

type shipmentRepository repository

type ShipmentRepositoryInterface interface {
	GetAllShipmentsWithFilters(filters dto.AdminGetShipmentsRequest) ([]dto.AdminShipmentListResponse, int64, error)
	GetShipmentByID(id int) (*dto.AdminShipmentDetailResponse, error)
	CreateShipment(shipment *models.Shipment) error
	UpdateShipment(id int, updates map[string]interface{}) error
	CheckDuplicateShipment(orderID int) (bool, error)
	GetOrderByID(orderID int) (*models.Order, error)
}

func (r *shipmentRepository) GetAllShipmentsWithFilters(filters dto.AdminGetShipmentsRequest) ([]dto.AdminShipmentListResponse, int64, error) {
	var shipments []dto.AdminShipmentListResponse
	var total int64

	query := r.Options.Postgres.Table("shipments s").
		Select(`
			s.id, s.order_id, s.shipping_method_id, s.tracking_number, s.status,
			s.shipped_at, s.delivered_at, s.created_at, o.total_amount,
			u.email as customer_email 
		`).
		Joins("LEFT JOIN orders o ON s.order_id = o.id").
		Joins("LEFT JOIN users u ON o.user_id = u.id").
		Joins("LEFT JOIN shipping_methods sm ON s.shipping_method_id = sm.id")

	if filters.Status != nil && *filters.Status != "" {
		query = query.Where("s.status = ?", *filters.Status)
	}
	if filters.OrderID != nil && *filters.OrderID > 0 {
		query = query.Where("s.order_id = ?", *filters.OrderID)
	}
	if filters.TrackingNumber != nil && *filters.TrackingNumber != "" {
		query = query.Where("s.tracking_number LIKE ?", "%"+*filters.TrackingNumber+"%")
	}
	if filters.ShippingMethodID != nil && *filters.ShippingMethodID > 0 {
		query = query.Where("s.shipping_method_id = ?", *filters.ShippingMethodID)
	}
	if filters.DateFrom != nil && *filters.DateFrom != "" {
		if date, err := time.Parse("2006-01-02", *filters.DateFrom); err == nil {
			query = query.Where("s.created_at >= ?", date)
		}
	}
	if filters.DateTo != nil && *filters.DateTo != "" {
		if date, err := time.Parse("2006-01-02", *filters.DateTo); err == nil {
			date = date.Add(24*time.Hour - time.Second)
			query = query.Where("s.created_at <= ?", date)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.PageSize
	if err := query.Limit(filters.PageSize).Offset(offset).Find(&shipments).Error; err != nil {
		return nil, 0, err
	}

	return shipments, total, nil
}

func (r *shipmentRepository) GetShipmentByID(id int) (*dto.AdminShipmentDetailResponse, error) {
	var shipment dto.AdminShipmentDetailResponse

	err := r.Options.Postgres.Table("shipments s").
		Select(`
			s.id, s.order_id, s.shipping_method_id, s.tracking_number, s.status,
			s.shipped_at, s.delivered_at, s.created_at, s.updated_at,
			o.id as order__id, o.status as order__status, o.total_amount as order__total_amount,
			o.created_at as order__created_at,sm.id as shipping_method__id, sm.name as shipping_method__name
		`).
		Joins("LEFT JOIN orders o ON s.order_id = o.id").
		Joins("LEFT JOIN shipping_methods sm ON s.shipping_method_id = sm.id").
		Where("s.id = ?", id).
		First(&shipment).Error

	if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *shipmentRepository) CreateShipment(shipment *models.Shipment) error {
	return r.Options.Postgres.Create(shipment).Error
}

func (r *shipmentRepository) UpdateShipment(id int, updates map[string]interface{}) error {
	return r.Options.Postgres.Model(&models.Shipment{}).Where("id = ?", id).Updates(updates).Error
}

func (r *shipmentRepository) CheckDuplicateShipment(orderID int) (bool, error) {
	var count int64
	err := r.Options.Postgres.Model(&models.Shipment{}).
		Where("order_id = ?", orderID).
		Count(&count).Error
	return count > 0, err
}

func (r *shipmentRepository) GetOrderByID(orderID int) (*models.Order, error) {
	var order models.Order
	err := r.Options.Postgres.First(&order, orderID).Error
	return &order, err
}