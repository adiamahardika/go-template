package usecases

import (
	"monitoring-service/app/models"
)

type categoryUsecase usecase

type CategoryUsecaseInterface interface {
	GetCategoryByID(id int, include_related bool) (*models.CategoryResponses, error)
	//GetAllCategory(included_deleted bool) (*[]models.Category, error)
	GetAllCategory(page, pageSize int, q string, included_deleted bool) (*[]models.Category, int64, error)
	CreateCategory(category *models.Category) (*models.Category, error)
	UpdateCategory(id int, updates map[string]interface{}) (*models.Category, error)
	SoftDeleteCategory(id int) error
	IsCategoryExist(categoryID int) (*models.Category, error)
}

func (u *categoryUsecase) GetCategoryByID(id int, include_related bool) (*models.CategoryResponses, error) {
	category, related, err := u.Options.Repository.Category.GetCategoryByID(id, include_related)

	if err != nil {
		return nil, err
	}

	response := category.ToCategoryResponse(related)
	return &response, nil
}

/*
func (u *categoryUsecase) GetAllCategory(included_deleted bool) (*[]models.Category, error) {
	allCategory, err := u.Options.Repository.Category.GetAllCategory(included_deleted)

	if err != nil {
		return nil, err
	}

	return allCategory, nil
}

*/

func (u *categoryUsecase) GetAllCategory(page, pageSize int, q string, included_deleted bool) (*[]models.Category, int64, error) {
	allCategory, total, err := u.Options.Repository.Category.GetAllCategory(page, pageSize, q, included_deleted)

	if err != nil {
		return nil, 0, err
	}

	return allCategory, total, nil
}

func (u *categoryUsecase) CreateCategory(category *models.Category) (*models.Category, error) {
	newCategory, err := u.Options.Repository.Category.CreateCategory(category)

	if err != nil {
		return nil, err
	}

	return newCategory, nil
}

func (u *categoryUsecase) UpdateCategory(id int, updates map[string]interface{}) (*models.Category, error) {
	category, err := u.Options.Repository.Category.UpdateCategory(id, updates)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (u *categoryUsecase) SoftDeleteCategory(id int) error {
	return u.Options.Repository.Category.SoftDeleteCategory(id)
}

func (u *categoryUsecase) IsCategoryExist(categoryID int) (*models.Category, error) {
	category, err := u.Options.Repository.Category.IsCategoryExist(categoryID)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, nil
	}

	return category, nil
}
