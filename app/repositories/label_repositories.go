package repositories

import "monitoring-service/app/models"

type labelRepository repository

type LabelRepositoryInterface interface {
	GetAll() ([]*models.Label, error)
}

func (r *labelRepository) GetAll() ([]*models.Label, error) {
	var labels []*models.Label
	err := r.Options.Postgres.Debug().Table("label").Find(&labels).Error
	if err != nil {
		return nil, err
	}
	return labels, nil
}
