package repositories

import "monitoring-service/app/models"

type shippingPaymentRepository struct {
	*repository
}

func NewShippingPaymentRepository(repo *repository) ShippingPaymentRepositoryInterface {
	return &shippingPaymentRepository{
		repository: repo,
	}
}

func (r *shippingPaymentRepository) CreateShippingMethod(method models.ShippingMethod) (*models.ShippingMethod, error) {
	if err := r.Options.Postgres.Create(&method).Error; err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *shippingPaymentRepository) GetShippingMethods(filter models.ShippingMethodFilter) ([]models.ShippingMethod, int64, error) {
	var methods []models.ShippingMethod
	var total int64

	query := r.Options.Postgres.Model(&models.ShippingMethod{})

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Active {
		query = query.Where("deleted_at IS NULL")
	} else {
		query = query.Unscoped()
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := query.Offset(offset).Limit(filter.Limit).Find(&methods).Error; err != nil {
		return nil, 0, err
	}

	return methods, total, nil
}

func (r *shippingPaymentRepository) GetShippingMethodByID(id int) (*models.ShippingMethod, error) {
	var method models.ShippingMethod
	if err := r.Options.Postgres.First(&method, id).Error; err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *shippingPaymentRepository) UpdateShippingMethod(method models.ShippingMethod) error {
	return r.Options.Postgres.Save(&method).Error
}

func (r *shippingPaymentRepository) DeleteShippingMethod(id int) error {
	return r.Options.Postgres.Delete(&models.ShippingMethod{}, id).Error
}

func (r *shippingPaymentRepository) CreatePaymentMethod(method models.PaymentMethod) (*models.PaymentMethod, error) {
	if err := r.Options.Postgres.Create(&method).Error; err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *shippingPaymentRepository) GetPaymentMethods(filter models.PaymentMethodFilter) ([]models.PaymentMethod, int64, error) {
	var methods []models.PaymentMethod
	var total int64

	query := r.Options.Postgres.Model(&models.PaymentMethod{})

	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}

	if filter.Active {
		query = query.Where("deleted_at IS NULL")
	} else {
		query = query.Unscoped()
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := query.Offset(offset).Limit(filter.Limit).Find(&methods).Error; err != nil {
		return nil, 0, err
	}

	return methods, total, nil
}

func (r *shippingPaymentRepository) GetPaymentMethodByID(id int) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	if err := r.Options.Postgres.First(&method, id).Error; err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *shippingPaymentRepository) UpdatePaymentMethod(method models.PaymentMethod) error {
	return r.Options.Postgres.Save(&method).Error
}

func (r *shippingPaymentRepository) DeletePaymentMethod(id int) error {
	return r.Options.Postgres.Delete(&models.PaymentMethod{}, id).Error
}
