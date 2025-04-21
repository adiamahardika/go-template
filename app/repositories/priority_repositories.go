package repositories

import (
	"monitoring-service/app/models"
)

type priorityRepository repository

type PriorityRepositoryInterface interface {
	GetAll() ([]models.Priority, error)
}

func (r *priorityRepository) GetAll() ([]models.Priority, error) {
	var statuses []models.Priority
	err := r.Options.Postgres.Debug().Table("priority").Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
