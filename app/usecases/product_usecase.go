package usecases

import (
	"monitoring-service/app/models"
)

type productUsecase usecase

type productUsecaseInterface interface {
	GetProductByID(id int) (*models.ProductResponse, error)
}

func (u *productUsecase) GetProductByID(id int) (*models.ProductResponse, error) {
	product, related, err := u.Options.Repository.Product.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	response := product.ToProductResponse(related)
	return &response, nil
}
