package usecases

import (
	"monitoring-service/app/models"
	"monitoring-service/pkg/customerror"

	"gorm.io/gorm"
)

type ShippingPaymentUsecaseInterface interface {
	CreateShippingMethod(request models.ShippingMethodRequest) (*models.ShippingMethodResponse, error)
	GetShippingMethods(filter models.ShippingMethodFilter) ([]models.ShippingMethodResponse, models.Pagination, error)
	GetShippingMethodByID(id int) (*models.ShippingMethodResponse, error)
	UpdateShippingMethod(id int, request models.ShippingMethodRequest) (*models.ShippingMethodResponse, error)
	DeleteShippingMethod(id int) error

	CreatePaymentMethod(request models.PaymentMethodRequest) (*models.PaymentMethodResponse, error)
	GetPaymentMethods(filter models.PaymentMethodFilter) ([]models.PaymentMethodResponse, models.Pagination, error)
	GetPaymentMethodByID(id int) (*models.PaymentMethodResponse, error)
	UpdatePaymentMethod(id int, request models.PaymentMethodRequest) (*models.PaymentMethodResponse, error)
	DeletePaymentMethod(id int) error
}

type shippingPaymentUsecase struct {
	*usecase
}

func (u *shippingPaymentUsecase) CreateShippingMethod(request models.ShippingMethodRequest) (*models.ShippingMethodResponse, error) {
	method := models.ShippingMethod{
		Name:          request.Name,
		Cost:          request.Cost,
		EstimatedDays: request.EstimatedDays,
		Description:   request.Description,
	}

	created, err := u.options.Repository.ShippingPayment.CreateShippingMethod(method)
	if err != nil {
		return nil, customerror.NewInternalServiceError(err.Error())
	}

	return &models.ShippingMethodResponse{
		ID:            created.ID,
		Name:          created.Name,
		Cost:          created.Cost,
		EstimatedDays: created.EstimatedDays,
		Description:   created.Description,
		CreatedAt:     created.CreatedAt,
		UpdatedAt:     created.UpdatedAt,
	}, nil
}

func (u *shippingPaymentUsecase) GetShippingMethods(filter models.ShippingMethodFilter) ([]models.ShippingMethodResponse, models.Pagination, error) {
	methods, total, err := u.options.Repository.ShippingPayment.GetShippingMethods(filter)
	if err != nil {
		return nil, models.Pagination{}, customerror.NewInternalServiceError(err.Error())
	}

	responses := make([]models.ShippingMethodResponse, len(methods))
	for i, m := range methods {
		responses[i] = models.ShippingMethodResponse{
			ID:            m.ID,
			Name:          m.Name,
			Cost:          m.Cost,
			EstimatedDays: m.EstimatedDays,
			Description:   m.Description,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
		}
	}

	totalPages := 0
	if filter.Limit > 0 {
		totalPages = int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	}

	pagination := models.Pagination{
		Page:      filter.Page,
		PageSize:  filter.Limit,
		Total:     int(total),
		TotalPage: totalPages,
	}

	return responses, pagination, nil
}

func (u *shippingPaymentUsecase) GetShippingMethodByID(id int) (*models.ShippingMethodResponse, error) {
	method, err := u.options.Repository.ShippingPayment.GetShippingMethodByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerror.NewNotFoundError("shipping method not found")
		}
		return nil, customerror.NewInternalServiceError(err.Error())
	}

	return &models.ShippingMethodResponse{
		ID:            method.ID,
		Name:          method.Name,
		Cost:          method.Cost,
		EstimatedDays: method.EstimatedDays,
		Description:   method.Description,
		CreatedAt:     method.CreatedAt,
		UpdatedAt:     method.UpdatedAt,
	}, nil
}

func (u *shippingPaymentUsecase) UpdateShippingMethod(id int, request models.ShippingMethodRequest) (*models.ShippingMethodResponse, error) {
	method := models.ShippingMethod{
		ID:            id,
		Name:          request.Name,
		Cost:          request.Cost,
		EstimatedDays: request.EstimatedDays,
		Description:   request.Description,
	}

	err := u.options.Repository.ShippingPayment.UpdateShippingMethod(method)
	if err != nil {
		return nil, customerror.NewInternalServiceError(err.Error())
	}

	updated, err := u.options.Repository.ShippingPayment.GetShippingMethodByID(id)
	if err != nil {
		return nil, customerror.NewInternalServiceError(err.Error())
	}

	return &models.ShippingMethodResponse{
		ID:            updated.ID,
		Name:          updated.Name,
		Cost:          updated.Cost,
		EstimatedDays: updated.EstimatedDays,
		Description:   updated.Description,
		CreatedAt:     updated.CreatedAt,
		UpdatedAt:     updated.UpdatedAt,
	}, nil
}

func (u *shippingPaymentUsecase) DeleteShippingMethod(id int) error {
	return u.options.Repository.ShippingPayment.DeleteShippingMethod(id)
}

func (u *shippingPaymentUsecase) CreatePaymentMethod(request models.PaymentMethodRequest) (*models.PaymentMethodResponse, error) {
	method := models.PaymentMethod{
		Name:        request.Name,
		Description: request.Description,
	}

	created, err := u.options.Repository.ShippingPayment.CreatePaymentMethod(method)
	if err != nil {
		return nil, customerror.NewInternalServiceError(err.Error())
	}

	return &models.PaymentMethodResponse{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}, nil
}

func (u *shippingPaymentUsecase) GetPaymentMethods(filter models.PaymentMethodFilter) ([]models.PaymentMethodResponse, models.Pagination, error) {
	methods, total, err := u.options.Repository.ShippingPayment.GetPaymentMethods(filter)
	if err != nil {
		return nil, models.Pagination{}, customerror.NewInternalServiceError(err.Error())
	}

	responses := make([]models.PaymentMethodResponse, len(methods))
	for i, m := range methods {
		responses[i] = models.PaymentMethodResponse{
			ID:          m.ID,
			Name:        m.Name,
			Description: m.Description,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
			DeletedAt:   m.DeletedAt,
		}
	}

	totalPages := 0
	if filter.Limit > 0 {
		totalPages = int((total + int64(filter.Limit) - 1) / int64(filter.Limit))
	}

	pagination := models.Pagination{
		Page:      filter.Page,
		PageSize:  filter.Limit,
		Total:     int(total),
		TotalPage: totalPages,
	}

	return responses, pagination, nil
}

func (u *shippingPaymentUsecase) GetPaymentMethodByID(id int) (*models.PaymentMethodResponse, error) {
	method, err := u.options.Repository.ShippingPayment.GetPaymentMethodByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerror.NewNotFoundError("payment method not found")
		}
		return nil, customerror.NewInternalServiceError(err.Error())
	}

	return &models.PaymentMethodResponse{
		ID:          method.ID,
		Name:        method.Name,
		Description: method.Description,
		CreatedAt:   method.CreatedAt,
		UpdatedAt:   method.UpdatedAt,
		DeletedAt:   method.DeletedAt,
	}, nil
}

func (u *shippingPaymentUsecase) UpdatePaymentMethod(id int, request models.PaymentMethodRequest) (*models.PaymentMethodResponse, error) {
	method := models.PaymentMethod{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
	}

	err := u.options.Repository.ShippingPayment.UpdatePaymentMethod(method)
	if err != nil {
		return nil, customerror.NewInternalServiceError(err.Error())
	}

	updated, err := u.options.Repository.ShippingPayment.GetPaymentMethodByID(id)
	if err != nil {
		return nil, customerror.NewInternalServiceError(err.Error())
	}

	return &models.PaymentMethodResponse{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}, nil
}

func (u *shippingPaymentUsecase) DeletePaymentMethod(id int) error {
	return u.options.Repository.ShippingPayment.DeletePaymentMethod(id)
}
