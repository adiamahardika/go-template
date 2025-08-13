package repositories

import (
	"gorm.io/gorm"
)

type CategoryRepository interface {
       ListCategories() ([]map[string]interface{}, error)
}

type categoryRepository struct {
       DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
       return &categoryRepository{DB: db}
}

func (r *categoryRepository) ListCategories() ([]map[string]interface{}, error) {
       var results []struct {
	       ID   int
	       Name string
       }
       err := r.DB.Table("categories").Select("id, name").Where("deleted_at IS NULL").Scan(&results).Error
       if err != nil {
	       return nil, err
       }
       categories := make([]map[string]interface{}, 0, len(results))
       for _, row := range results {
	       categories = append(categories, map[string]interface{}{
		       "id":   row.ID,
		       "name": row.Name,
	       })
       }
       return categories, nil
}
