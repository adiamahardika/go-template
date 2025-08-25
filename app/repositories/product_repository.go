// file: app/repositories/product_repository.go
package repositories

import (
	"monitoring-service/app/models"

	"gorm.io/gorm"
)

type productRepository repository

type ProductRepositoryInterface interface {
	GetProductByID(productID int) (*models.Product, []models.ProductRelated, error)
	DecrementProductStock(productID int, quantity int) error
	IncrementProductStock(productID int, quantity int) error
	GetProductForUpdate(productID int) (*models.Product, error)
	GetProductByIDOnly(productID int) (*models.Product, error)
}

// ✅ Method yang sudah ada sebelumnya
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

// ✅ Method baru untuk mendapatkan product saja (tanpa related)
func (r *productRepository) GetProductByIDOnly(productID int) (*models.Product, error) {
	var product models.Product
	err := r.Options.Postgres.
		Where("id = ? AND deleted_at IS NULL", productID).
		First(&product).Error
	return &product, err
}

// ✅ Decrement stock dengan validation
func (r *productRepository) DecrementProductStock(productID int, quantity int) error {
	return r.Options.Postgres.
		Model(&models.Product{}).
		Where("id = ? AND stock >= ?", productID, quantity).
		Update("stock", gorm.Expr("stock - ?", quantity)).Error
}

// ✅ Increment stock
func (r *productRepository) IncrementProductStock(productID int, quantity int) error {
	return r.Options.Postgres.
		Model(&models.Product{}).
		Where("id = ?", productID).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

// ✅ Get product with lock for update
func (r *productRepository) GetProductForUpdate(productID int) (*models.Product, error) {
	var product models.Product
	err := r.Options.Postgres.
		Set("gorm:query_option", "FOR UPDATE").
		Where("id = ? AND deleted_at IS NULL", productID).
		First(&product).Error
	return &product, err
}