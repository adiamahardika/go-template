package usecases

import "monitoring-service/app/repositories"

type CategoryUsecase interface {
	ListCategories() ([]map[string]interface{}, error)
}

type categoryUsecase struct {
	CategoryRepo repositories.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repositories.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{CategoryRepo: categoryRepo}
}

func (u *categoryUsecase) ListCategories() ([]map[string]interface{}, error) {
	return u.CategoryRepo.ListCategories()
}
