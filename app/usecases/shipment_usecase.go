package usecases

import (
	"errors"
	"monitoring-service/app/models/dto"
	
)

type shipmentUsecase usecase

type ShipmentUsecaseInterface interface {
	GetShipmentByOrderID(orderID int, userID int) (*dto.ShipmentResponse, error)
}


func (s *shipmentUsecase) GetShipmentByOrderID(orderID int, userID int) (*dto.ShipmentResponse, error) {
	order, err := s.Options.Repository.Shipment.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.New("shipment not found")
	}

	if *order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	shipment, err := s.Options.Repository.Shipment.GetShipmentByOrderID(orderID)
	if err != nil {
		return nil, errors.New("shipment not found")
	}

	var trackingNumber *string
	if shipment.TrackingNumber != nil {
		trackingNumber = shipment.TrackingNumber
	} else {
		noTracking := "Tracking number not yet available"
		trackingNumber = &noTracking
	}

	response := &dto.ShipmentResponse{
		ID:               shipment.ID,
		OrderID:          shipment.OrderID,
		ShippingMethodID: shipment.ShippingMethodID,
		TrackingNumber:   trackingNumber,
		ShippedAt:        shipment.ShippedAt,
		DeliveredAt:      shipment.DeliveredAt,
		Status:           shipment.Status,
		CreatedAt:        shipment.CreatedAt,
		UpdatedAt:        shipment.UpdatedAt,
		Order:            shipment.Order,
		ShippingMethod:   shipment.ShippingMethod,
	}

	return response, nil
}