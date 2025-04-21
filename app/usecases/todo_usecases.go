package usecases

import "monitoring-service/app/models"

type todoUsecase usecase

type TodoUsecaseInterface interface {
	GetAll() ([]*models.Todo, error)
}

func (u *todoUsecase) GetAll() ([]*models.Todo, error) {
	var todos []*models.Todo
	var err error

	todos, err = u.Options.Repository.Todo.GetAll()
	if err != nil {
		return nil, err
	}
	return todos, nil
}
