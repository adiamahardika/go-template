package controllers

import (
	"errors"
	"log"
	"monitoring-service/app/helpers"
	"monitoring-service/app/models"
	"net/http"
	"strconv"

	"github.com/ezartsh/validet"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	customerror "monitoring-service/pkg/customerror"
)

type categoryController controller

type CategoryControllerInterface interface {
	GetCategoryByID(c echo.Context) error
	GetAllCategory(c echo.Context) error
	CreateCategory(c echo.Context) error
	UpdateCategory(c echo.Context) error
	SoftDeleteCategory(c echo.Context) error
}

func (ctrl *categoryController) GetCategoryByID(c echo.Context) error {
	var (
		req      *models.CategoryAdmin
		category *models.CategoryResponses
		err      error
		id       string
	)

	id = c.Param("id")

	mapReq := make(map[string]any)
	mapReq["id"] = id

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"id": validet.String{Required: true, Custom: func(v string, path validet.PathKey, lookup validet.Lookup) error {
				if _, err := strconv.Atoi(v); err != nil {
					return customerror.NewBadRequestError("Invalid category ID format")
				}
				return nil
			}},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		err := customerror.NewBadRequestError(err.Error())
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), errorBags.Errors, nil, nil)
	}

	categoryID, _ := strconv.Atoi(id)

	if err := c.Bind(&req); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid request"})
	}

	log.Printf("Fetching category with Include related: %t", req.IncludeRelated)
	category, err = ctrl.Options.UseCases.Category.GetCategoryByID(categoryID, req.IncludeRelated)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = customerror.NewNotFoundError("Category not found")
		}

		var nfErr customerror.NotFoundError
		if errors.As(err, &nfErr) {
			return helpers.Response(c, http.StatusNotFound, []string{"Category not found"})
		}

		return helpers.Response(c, customerror.GetStatusCode(err), []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Caregory retrieved successfully"}, category, nil)
}

func (ctrl *categoryController) GetAllCategory(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}

	q := c.QueryParam("q")
	includeDeleted := false
	if c.QueryParam("include_deleted") == "true" {
		includeDeleted = true
	}

	mapParams := map[string]any{
		"page":            page,
		"page_size":       pageSize,
		"q":               q,
		"include_deleted": includeDeleted,
	}
	schema := validet.NewSchema(
		mapParams,
		map[string]validet.Rule{
			"page": validet.Numeric[int]{Required: false, Custom: func(v int, path validet.PathKey, lookup validet.Lookup) error {
				if v != 0 && v < 1 {
					return customerror.NewBadRequestError("Page must be minimum of 1")
				}
				return nil
			}},
			"page_size": validet.Numeric[int]{Required: false, Custom: func(v int, path validet.PathKey, lookup validet.Lookup) error {
				if v != 0 && v < 1 {
					return customerror.NewBadRequestError("page_size must be minimum of 1 or maximum 100")
				}
				return nil
			}},
			"q":               validet.String{Required: false, Max: 100},
			"include_deleted": validet.Boolean{Required: false},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, errorBags.Errors, nil, nil)
	}

	category, total, err := ctrl.Options.UseCases.Category.GetAllCategory(page, pageSize, q, includeDeleted)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
	}

	if len(*category) == 0 {
		return helpers.Response(c, http.StatusNotFound, []string{"No categories found"})
	}
	return helpers.StandardResponse(c, http.StatusOK, []string{"Categories fetched successfully"}, map[string]interface{}{
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
		"category": category,
	}, nil)
}

func (ctrl *categoryController) CreateCategory(c echo.Context) error {
	var (
		req models.Category
		//		name string
		description string
	)
	if err := c.Bind(&req); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid request"})
	}

	if req.Description != nil {
		description = *req.Description
	}

	mapReq := map[string]any{
		"name":        req.Name,
		"description": description,
	}

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"name":        validet.String{Required: true, Min: 1, Max: 100},
			"description": validet.String{Required: false, Max: 500},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, errorBags.Errors, nil, nil)
	}

	newCategory, err := ctrl.Options.UseCases.Category.CreateCategory(&req)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
	}
	return helpers.StandardResponse(c, http.StatusCreated, []string{"Category created successfully"}, newCategory, nil)
}

func (ctrl *categoryController) UpdateCategory(c echo.Context) error {
	var req models.Category
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid category ID"})
	}

	if id != 0 && id < 1 {
		return helpers.Response(c, http.StatusBadRequest, []string{"Category ID must be greater than 0"})
	}

	if err := c.Bind(&req); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid update request"})
	}

	mapReq := map[string]any{}
	if req.Name != "" {
		mapReq["name"] = req.Name
	}
	if req.Description != nil {
		mapReq["description"] = *req.Description
	}

	if len(mapReq) == 0 {
		return helpers.Response(c, http.StatusBadRequest, []string{"No fields to update"})
	}

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"name":        validet.String{Required: false, Max: 100},
			"description": validet.String{Required: false, Max: 500},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()

	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, errorBags.Errors, nil, nil)
	}

	category, err := ctrl.Options.UseCases.Category.UpdateCategory(id, mapReq)

	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
	}
	return helpers.StandardResponse(c, http.StatusOK, []string{"Category updated successfully"}, category, nil)
}

func (ctrl *categoryController) SoftDeleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid or empty category ID"})
	}

	mapReq := make(map[string]any)
	if id != 0 && id > 0 {
		mapReq["id"] = id
	}

	if len(mapReq) == 0 {
		return helpers.Response(c, http.StatusBadRequest, []string{"Category ID is required"})
	}

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"id": validet.Numeric[int]{Required: true, Custom: func(v int,
				path validet.PathKey, lookup validet.Lookup) error {
				if v != 0 && v < 1 {
					return customerror.NewBadRequestError("Product ID must be greater than 0")
				}
				return nil
			}},
		},
		validet.Options{},
	)
	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, errorBags.Errors, nil, nil)
	}

	if err := ctrl.Options.UseCases.Category.SoftDeleteCategory(id); err != nil {
		return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
	}
	return helpers.StandardResponse(c, http.StatusOK, []string{"Category deleted successfully"}, nil, nil)
}
