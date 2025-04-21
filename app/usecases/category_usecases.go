package usecases

import "monitoring-service/app/models"

type categoryUsecase usecase

type CategoryUsecaseInterface interface {
	GetAll() ([]*models.Category, error)
}

func (u *categoryUsecase) GetAll() ([]*models.Category, error) {
	var categories []*models.Category
	var err error

	categories, err = u.Options.Repository.Category.GetAll()
	if err != nil {
		return nil, err
	}
	return categories, nil
}
