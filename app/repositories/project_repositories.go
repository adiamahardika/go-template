package repositories

import "monitoring-service/app/models"

type projectRepository repository

type ProjectRepositoryInterface interface {
	GetAll() ([]*models.Project, error)
}

func (r *projectRepository) GetAll() ([]*models.Project, error) {
	var projects []*models.Project
	err := r.Options.Postgres.Debug().Table("project").Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}
