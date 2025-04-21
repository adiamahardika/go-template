package repositories

import "monitoring-service/app/models"

type categoryRepository repository

type CategoryRepositoryInterface interface {
    GetAll() ([]*models.Category, error)
}

func (r *categoryRepository) GetAll() ([]*models.Category, error) {
    var categories []*models.Category
    err := r.Options.Postgres.Debug().Table("category").Find(&categories).Error
    if err != nil {
        return nil, err
    }
    return categories, nil
}
