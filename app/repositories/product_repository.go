package repositories

import (
	"monitoring-service/app/models"
)

type productRepository repository

type ProductRepositoryInterface interface {
	GetProductByID(productID int) (*models.Product, []models.ProductRelated, error)
}

func (r *productRepository) GetProductByID(productID int) (*models.Product, []models.ProductRelated, error) {
	var product models.Product
	var related []models.ProductRelated

	err := r.Options.Postgres.
		Joins("JOIN categories ON categories.id = products.category_id AND categories.deleted_at IS NULL").
		Preload("Category").
		Where("products.id = ? AND products.deleted_at IS NULL", productID).
		First(&product).Error

	if err != nil {
		return nil, nil, err
	}

	err = r.Options.Postgres.Table("products").
		Select("id, name, price, image_url").
		Where("category_id = ? AND id <> ? AND deleted_at IS NULL", product.CategoryID, product.ID).Limit(4).Find(&related).Error

	if err != nil {
		return nil, nil, err
	}

	return &product, related, nil
}
