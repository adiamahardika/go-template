package usecases

import "monitoring-service/app/models"

type labelUsecase usecase

type LabelUsecaseInterface interface {
	GetAll() ([]*models.Label, error)
}

func (u *labelUsecase) GetAll() ([]*models.Label, error) {
	var labels []*models.Label
	var err error
	
	labels, err = u.Options.Repository.Label.GetAll()
	if err != nil {
		return nil, err
	}
	return labels, nil
}
