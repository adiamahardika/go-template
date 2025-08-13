package usecases

import (
	"monitoring-service/app/repositories"
)

type ProductUsecase interface {
	ListProducts(page, pageSize, categoryID int, search, sort string) ([]map[string]interface{}, map[string]interface{}, error)
}

type productUsecase struct {
	ProductRepo repositories.ProductRepository
}

func NewProductUsecase(productRepo repositories.ProductRepository) ProductUsecase {
	return &productUsecase{ProductRepo: productRepo}
}

func (u *productUsecase) ListProducts(page, pageSize, categoryID int, search, sort string) ([]map[string]interface{}, map[string]interface{}, error) {
	products, total, err := u.ProductRepo.ListProducts(page, pageSize, categoryID, search, sort)
	if err != nil {
		return nil, nil, err
	}
	meta := map[string]interface{}{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}
	return products, meta, nil
}
