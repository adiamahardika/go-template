package usecases

import "monitoring-service/app/models"

type priorityUsecase usecase

type PriorityUsecaseInterface interface {
	GetAll() ([]models.Priority, error)
}

func (u *priorityUsecase) GetAll() ([]models.Priority, error) {
	var priorities []models.Priority
	var err error

	priorities, err = u.Options.Repository.Priority.GetAll()
	if err != nil {
		return nil, err
	}
	return priorities, nil
}
