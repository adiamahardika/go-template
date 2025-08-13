package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"monitoring-service/app/usecases"
)

type ProductController struct {
	ProductUsecase usecases.ProductUsecase
}

func NewProductController(productUsecase usecases.ProductUsecase) *ProductController {
	return &ProductController{ProductUsecase: productUsecase}
}

func (pc *ProductController) ListProducts(c echo.Context) error {
	// Query params
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("page_size")
	categoryIDStr := c.QueryParam("category_id")
	search := c.QueryParam("search")
	sort := c.QueryParam("sort")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize < 1 {
		pageSize = 10
	}
	categoryID, _ := strconv.Atoi(categoryIDStr)
	if categoryIDStr == "" {
		categoryID = 0
	}

	products, meta, err := pc.ProductUsecase.ListProducts(page, pageSize, categoryID, search, sort)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": products,
		"meta": meta,
	})
}
