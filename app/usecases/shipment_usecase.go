package usecases

import (
	"context"
	"errors"
	"monitoring-service/app/models"
	"monitoring-service/app/models/dto"
	"time"
)

type shipmentUsecase usecase

type ShipmentUsecaseInterface interface {
	GetAllShipmentsWithFilters(filters dto.AdminGetShipmentsRequest) ([]dto.AdminShipmentListResponse, models.Pagination, error)
	GetShipmentByID(id int) (*dto.AdminShipmentDetailResponse, error)
	UpdateShipment(ctx context.Context, id int, req dto.AdminUpdateShipmentRequest) (*dto.AdminShipmentDetailResponse, error)
	CreateShipment(ctx context.Context, req dto.AdminCreateShipmentRequest) (*dto.AdminShipmentListResponse, error)
}

func (s *shipmentUsecase) GetAllShipmentsWithFilters(filters dto.AdminGetShipmentsRequest) ([]dto.AdminShipmentListResponse, models.Pagination, error) {
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PageSize <= 0 {
		filters.PageSize = 10
	}

	shipments, total, err := s.Options.Repository.Shipment.GetAllShipmentsWithFilters(filters)
	if err != nil {
		return nil, models.Pagination{}, err
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	pagination := models.Pagination{
		Page:      filters.Page,
		PageSize:  filters.PageSize,
		Total:     int(total),
		TotalPage: totalPages,
	}

	return shipments, pagination, nil
}

func (s *shipmentUsecase) GetShipmentByID(id int) (*dto.AdminShipmentDetailResponse, error) {
	return s.Options.Repository.Shipment.GetShipmentByID(id)
}

func (s *shipmentUsecase) UpdateShipment(ctx context.Context, id int, req dto.AdminUpdateShipmentRequest) (*dto.AdminShipmentDetailResponse, error) {
	_, err := s.Options.Repository.Shipment.GetShipmentByID(id)
	if err != nil {
		return nil, errors.New("shipment not found")
	}

	updates := make(map[string]interface{})

	if req.TrackingNumber != nil {
		updates["tracking_number"] = *req.TrackingNumber
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		if *req.Status == "shipped" && req.ShippedAt == nil {
			now := time.Now()
			updates["shipped_at"] = now
		}
		if *req.Status == "delivered" && req.DeliveredAt == nil {
			now := time.Now()
			updates["delivered_at"] = now
		}
	}
	if req.ShippedAt != nil {
		updates["shipped_at"] = *req.ShippedAt
	}
	if req.DeliveredAt != nil {
		updates["delivered_at"] = *req.DeliveredAt
	}

	if len(updates) > 0 {
		if err := s.Options.Repository.Shipment.UpdateShipment(id, updates); err != nil {
			return nil, err
		}
	}

	return s.Options.Repository.Shipment.GetShipmentByID(id)
}

func (s *shipmentUsecase) CreateShipment(ctx context.Context, req dto.AdminCreateShipmentRequest) (*dto.AdminShipmentListResponse, error) {
	order, err := s.Options.Repository.Shipment.GetOrderByID(req.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.Status == nil || (*order.Status != "paid" && *order.Status != "processing") {
		return nil, errors.New("order is not eligible for shipment creation (must be paid or processing)")
	}

	exists, err := s.Options.Repository.Shipment.CheckDuplicateShipment(req.OrderID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("shipment already exists for this order")
	}

	shipment := &models.Shipment{
		OrderID:          &req.OrderID,
		ShippingMethodID: &req.ShippingMethodID,
		TrackingNumber:   &req.TrackingNumber,
		Status:           &req.Status,
	}

	if req.Status == "shipped" {
		now := time.Now()
		shipment.ShippedAt = &now
	}

	if err := s.Options.Repository.Shipment.CreateShipment(shipment); err != nil {
		return nil, err
	}

	response := &dto.AdminShipmentListResponse{
		ID:               shipment.ID,
		OrderID:          *shipment.OrderID,
		ShippingMethodID: *shipment.ShippingMethodID,
		TrackingNumber:   shipment.TrackingNumber,
		Status:           *shipment.Status,
		ShippedAt:        shipment.ShippedAt,
		CreatedAt:        shipment.CreatedAt,
	}
	if order.TotalAmount != nil {
		response.TotalAmount = *order.TotalAmount
	}

	return response, nil
}