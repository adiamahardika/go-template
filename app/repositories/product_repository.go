package repositories

import (
	"errors"
	"monitoring-service/app/models"
	"time"
)

type productRepository repository

type ProductRepositoryInterface interface {
	GetProductByID(productID int) (*models.Product, []models.ProductRelated, error)
	CreateProduct(product *models.Product) (*models.Product, error)
	UpdateProduct(productID int, updates map[string]interface{}) (*models.Product, error)
	SoftDeleteProduct(productID int) error
	//GetAllProduct(included_deleted bool) (*[]models.Product, error)
	GetAllProduct(page, pageSize int, q string, categoryID int, included_deleted bool, sortBy string, sortOrder string) (*[]models.Product, int64, error)
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

func (r *productRepository) CreateProduct(product *models.Product) (*models.Product, error) {
	err := r.Options.Postgres.Create(product).Error

	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) GetAllProduct(page, pageSize int, q string, categoryID int, included_deleted bool, sortBy string, sortOrder string) (*[]models.Product, int64, error) {
	var allProduct []models.Product
	var total int64

	baseQuery := r.Options.Postgres.Model(&models.Product{})
	if !included_deleted {
		baseQuery = baseQuery.Where("deleted_at IS NULL")
	} else {
		baseQuery = baseQuery.Unscoped()
	}

	if q != "" {
		baseQuery = baseQuery.Where("name ILIKE ?", "%"+q+"%")
	}

	if categoryID > 0 {
		baseQuery = baseQuery.Where("category_id = ?", categoryID)
	}

	// Count total items
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated rows using a **new query** with same filters
	query := r.Options.Postgres.Model(&models.Product{}).Preload("Category")
	if !included_deleted {
		query = query.Where("deleted_at IS NULL")
	} else {
		query = query.Unscoped()
	}
	if q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if sortBy != "" {
		order := "asc"
		if sortOrder != "" {
			order = sortOrder
		}
		query = query.Order(sortBy + " " + order)
	}

	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&allProduct).Error; err != nil {
		return nil, 0, err
	}

	return &allProduct, total, nil
}

func (r *productRepository) UpdateProduct(productID int, updates map[string]interface{}) (*models.Product, error) {
	var product models.Product
	result := r.Options.Postgres.Model(&models.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("Product not found")
	}

	err := r.Options.Postgres.First(&product, productID).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) SoftDeleteProduct(productID int) error {
	result := r.Options.Postgres.Model(&models.Product{}).
		Where("id = ? AND deleted_at IS NULL", productID).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Product didn't found or already deleted")
	}
	return nil
}
