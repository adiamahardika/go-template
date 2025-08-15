package repositories

import (
	"monitoring-service/app/models"
)

type ShippingMethodRepositoryInterface interface {
	FindAllActive() ([]models.ShippingMethod, error)
	FindByID(id int) (*models.ShippingMethod, error)
}

type shippingMethodRepository struct {
	Options Options
}

func (r *shippingMethodRepository) FindAllActive() ([]models.ShippingMethod, error) {
	var methods []models.ShippingMethod
	err := r.Options.Postgres.Where("deleted_at IS NULL").Order("id ASC").Find(&methods).Error
	return methods, err
}

func (r *shippingMethodRepository) FindByID(id int) (*models.ShippingMethod, error) {
	var method models.ShippingMethod
	err := r.Options.Postgres.Where("id = ? AND deleted_at IS NULL", id).First(&method).Error
	if err != nil {
		return nil, err
	}
	return &method, nil
}
