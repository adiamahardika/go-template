package usecases

import (
	"monitoring-service/app/models"
)

type productUsecase usecase

type productUsecaseInterface interface {
	GetProductByID(id int) (*models.ProductResponse, error)
	CreateProduct(product *models.Product) (*models.Product, error)
	UpdateProduct(id int, updates map[string]interface{}) (*models.Product, error)
	SoftDeleteProduct(id int) error
	//GetAllProduct(included_deleted bool) (*[]models.Product, error)
	GetAllProduct(page, pageSize int, q string, categoryID int, included_deleted bool, sortBy string, sortOrder string) (*[]models.Product, int64, error)
}

func (u *productUsecase) GetProductByID(id int) (*models.ProductResponse, error) {
	product, related, err := u.Options.Repository.Product.GetProductByID(id)
	if err != nil {
		return nil, err
	}

	response := product.ToProductResponse(related)
	return &response, nil
}

func (u *productUsecase) CreateProduct(product *models.Product) (*models.Product, error) {
	newProduct, err := u.Options.Repository.Product.CreateProduct(product)

	if err != nil {
		return nil, err
	}

	return newProduct, nil
}

func (u *productUsecase) GetAllProduct(page, pageSize int, q string, categoryID int, included_deleted bool, sortBy string, sortOrder string) (*[]models.Product, int64, error) {
	allProduct, total, err := u.Options.Repository.Product.GetAllProduct(page, pageSize, q, categoryID, included_deleted, sortBy, sortOrder)

	if err != nil {
		return nil, 0, err
	}

	return allProduct, total, nil
}

func (u *productUsecase) UpdateProduct(id int, updates map[string]interface{}) (*models.Product, error) {
	product, err := u.Options.Repository.Product.UpdateProduct(id, updates)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *productUsecase) SoftDeleteProduct(id int) error {
	return u.Options.Repository.Product.SoftDeleteProduct(id)
}
