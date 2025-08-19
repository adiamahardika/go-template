package repositories

import (
	"errors"
	"monitoring-service/app/models"
	"time"
)

type categoryRepository repository

type CategoryRepositoryInterface interface {
	GetCategoryByID(categoryID int, include_related bool) (*models.Category, []models.Product, error)
	//GetAllCategory(included_deleted bool) (*[]models.Category, error)
	GetAllCategory(page, pageSize int, q string, included_deleted bool) (*[]models.Category, int64, error)
	CreateCategory(category *models.Category) (*models.Category, error)
	UpdateCategory(categoryID int, updates map[string]interface{}) (*models.Category, error)
	SoftDeleteCategory(categoryID int) error
	IsCategoryExist(categoryID int) (*models.Category, error)
}

func (r *categoryRepository) GetCategoryByID(categoryID int, include_related bool) (*models.Category, []models.Product, error) {
	var category models.Category
	var related []models.Product

	if include_related {
		err := r.Options.Postgres.
			Preload("Products", "deleted_at IS NULL").
			First(&category, "id = ? AND deleted_at IS NULL", categoryID).Error
		if err != nil {
			return nil, nil, err
		}
		related = category.Products
	} else {
		err := r.Options.Postgres.
			First(&category, "id = ? AND deleted_at IS NULL", categoryID).Error
		if err != nil {
			return nil, nil, err
		}
		related = nil
	}

	return &category, related, nil
}

func (r *categoryRepository) GetAllCategory(page, pageSize int, q string, included_deleted bool) (*[]models.Category, int64, error) {
	var allCategory []models.Category
	var total int64

	baseQuery := r.Options.Postgres.Model(&models.Category{})
	if !included_deleted {
		baseQuery = baseQuery.Where("deleted_at IS NULL")
	} else {
		baseQuery = baseQuery.Unscoped()
	}

	if q != "" {
		baseQuery = baseQuery.Where("name ILIKE ?", "%"+q+"%")
	}

	// Count total items
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated rows using a **new query** with same filters
	query := r.Options.Postgres.Model(&models.Category{})
	if !included_deleted {
		query = query.Where("deleted_at IS NULL")
	} else {
		query = query.Unscoped()
	}
	if q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
	}

	offset := (page - 1) * pageSize
	if err := query.Limit(pageSize).Offset(offset).Find(&allCategory).Error; err != nil {
		return nil, 0, err
	}

	return &allCategory, total, nil
}

func (r *categoryRepository) CreateCategory(category *models.Category) (*models.Category, error) {
	err := r.Options.Postgres.Create(category).Error

	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) UpdateCategory(categoryID int, updates map[string]interface{}) (*models.Category, error) {
	var category models.Category
	result := r.Options.Postgres.Model(&models.Category{}).
		Where("id = ? AND deleted_at IS NULL", categoryID).
		Updates(updates)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("Category not found")
	}

	err := r.Options.Postgres.First(&category, categoryID).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) SoftDeleteCategory(categoryID int) error {
	var count int64

	err := r.Options.Postgres.Model(&models.Product{}).
		Where("category_id = ? AND deleted_at IS NULL", categoryID).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("cannot delete category: there are active products referencing it")
	}

	result := r.Options.Postgres.Model(&models.Category{}).
		Where("id = ? AND deleted_at IS NULL", categoryID).
		Update("deleted_at", time.Now())

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("Category didn't found or already deleted")
	}
	return nil
}

func (r *categoryRepository) IsCategoryExist(categoryID int) (*models.Category, error) {
	var category models.Category

	err := r.Options.Postgres.
		First(&category, "id = ? AND deleted_at IS NULL", categoryID).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}
