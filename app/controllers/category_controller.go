package controllers

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"monitoring-service/app/usecases"
)

type CategoryController struct {
	CategoryUsecase usecases.CategoryUsecase
}

func NewCategoryController(categoryUsecase usecases.CategoryUsecase) *CategoryController {
	return &CategoryController{CategoryUsecase: categoryUsecase}
}

func (cc *CategoryController) ListCategories(c echo.Context) error {
	categories, err := cc.CategoryUsecase.ListCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": categories,
	})
}
