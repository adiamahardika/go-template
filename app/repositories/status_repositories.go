package repositories

import (
    "monitoring-service/app/models"
)

type statusRepository repository

type StatusRepositoryInterface interface {
    GetAllStatus() ([]models.Status, error)
}

func (r *statusRepository) GetAllStatus() ([]models.Status, error) {
    var statuses []models.Status
    err := r.Options.Postgres.Debug().Table("status").Find(&statuses).Error
    if err != nil {
        return nil, err
    }
    return statuses, nil
}
