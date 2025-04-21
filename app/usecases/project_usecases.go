package usecases

import "monitoring-service/app/models"

type projectUsecase usecase

type ProjectUsecaseInterface interface {
	GetAll() ([]*models.Project, error)
}

func (u *projectUsecase) GetAll() ([]*models.Project, error) {
	var projects []*models.Project
	var err error

	projects, err = u.Options.Repository.Project.GetAll()
	if err != nil {
		return nil, err
	}
	return projects, nil
}
