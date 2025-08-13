package repositories

import (
	"gorm.io/gorm"
)

type ProductRepository interface {
	ListProducts(page, pageSize, categoryID int, search, sort string) ([]map[string]interface{}, int, error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{DB: db}
}

func (r *productRepository) ListProducts(page, pageSize, categoryID int, search, sort string) ([]map[string]interface{}, int, error) {
	offset := (page - 1) * pageSize
	orderBy := "p.created_at DESC"
	switch sort {
	case "price_asc":
		orderBy = "p.price ASC"
	case "price_desc":
		orderBy = "p.price DESC"
	case "name_asc":
		orderBy = "p.name ASC"
	case "name_desc":
		orderBy = "p.name DESC"
	case "created_desc":
		orderBy = "p.created_at DESC"
	}
	var products []struct {
		ID           int
		Name         string
		Price        float64
		Stock        int
		CategoryName string
		ImageURL     string
		Description  string
	}
	db := r.DB.Table("products as p").Select("p.id, p.name, p.price, p.stock, c.name as category_name, p.image_url, LEFT(p.description, 100) as description").Joins("JOIN categories c ON p.category_id = c.id").Where("p.deleted_at IS NULL AND c.deleted_at IS NULL")
	if categoryID > 0 {
		db = db.Where("p.category_id = ?", categoryID)
	}
	if search != "" {
		db = db.Where("LOWER(p.name) LIKE ?", "%"+search+"%")
	}
	var total int64
	db.Count(&total)
	err := db.Order(orderBy).Limit(pageSize).Offset(offset).Scan(&products).Error
	if err != nil {
		return nil, 0, err
	}
	result := make([]map[string]interface{}, 0, len(products))
	for _, p := range products {
		result = append(result, map[string]interface{}{
			"id":            p.ID,
			"name":          p.Name,
			"price":         p.Price,
			"stock":         p.Stock,
			"category_name": p.CategoryName,
			"image_url":     p.ImageURL,
			"description":   p.Description,
		})
	}
	return result, int(total), nil
}
