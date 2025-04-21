package repositories

import "monitoring-service/app/models"

type todoRepository repository

type TodoRepositoryInterface interface {
	GetAll() ([]*models.Todo, error)
}

func (r *todoRepository) GetAll() ([]*models.Todo, error) {
	var todos []*models.Todo
	err := r.Options.Postgres.Debug().Table("todo").Find(&todos).Error
	if err != nil {
		return nil, err
	}
	return todos, nil
}
