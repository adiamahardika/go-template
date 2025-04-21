package usecases

import (
	"monitoring-service/app/models"
)

type statusUsecase usecase

type StatusUsecaseInterface interface {
    GetAllStatus() ([]models.Status, error)

}

func (u *statusUsecase) GetAllStatus() ([]models.Status, error) {
	var statuses []models.Status
	var err error

	statuses, err = u.Options.Repository.Status.GetAllStatus()
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
