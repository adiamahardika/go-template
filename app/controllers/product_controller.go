package controllers

import (
	"errors"
	"monitoring-service/app/helpers"
	"monitoring-service/app/models"
	"net/http"
	"strconv"

	"github.com/ezartsh/validet"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	customerror "monitoring-service/pkg/customerror"
)

type productController controller

type productControllerInterface interface {
	GetProductByID(c echo.Context) error
	CreateProduct(c echo.Context) error
	UpdateProduct(c echo.Context) error
	SoftDeleteProduct(c echo.Context) error
	GetAllProduct(c echo.Context) error
}

func (ctrl *productController) GetProductByID(c echo.Context) error {
	var (
		product *models.ProductResponse
		err     error
		id      string
	)

	id = c.Param("id")

	mapReq := make(map[string]any)
	mapReq["id"] = id

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"id": validet.String{Required: true, Custom: func(v string, path validet.PathKey, lookup validet.Lookup) error {
				if _, err := strconv.Atoi(v); err != nil {
					return customerror.NewBadRequestError("Invalid product ID format")
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

	productID, _ := strconv.Atoi(id)

	product, err = ctrl.Options.UseCases.Product.GetProductByID(productID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = customerror.NewNotFoundError("Product not found")
		}

		var nfErr customerror.NotFoundError
		if errors.As(err, &nfErr) {
			return helpers.Response(c, http.StatusNotFound, []string{"Product not found"})
		}

		return helpers.Response(c, customerror.GetStatusCode(err), []string{err.Error()})
	}

	if product.Stock <= 0 {
		product.Availability = "out_of_stock"
	}
	if product.Stock > 0 {
		product.Availability = "in_stock"
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Product retrieved successfully"}, product, nil)
}

func (ctrl *productController) CreateProduct(c echo.Context) error {
	var (
		req         models.Product
		description string
		categoryID  int
		imageURL    string
	)
	if err := c.Bind(&req); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid request"})
	}

	if req.CategoryID != nil {
		categoryID = *req.CategoryID
	}

	if req.Description != nil {
		description = *req.Description
	}

	if req.ImageURL != nil {
		imageURL = *req.ImageURL
	}

	mapReq := map[string]any{
		"name":        req.Name,
		"description": description,
		"price":       req.Price,
		"stock":       req.Stock,
		"category_id": categoryID,
		"image_url":   imageURL,
	}

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"name":        validet.String{Required: true, Min: 1, Max: 100},
			"description": validet.String{Required: false, Max: 500},
			"price":       validet.Numeric[float64]{Required: true, Min: 1},
			"stock":       validet.Numeric[int]{Required: true, Min: 1},
			"category_id": validet.Numeric[int]{Required: true, Min: 1},
			"image_url":   validet.String{Required: false, Max: 255},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, errorBags.Errors, nil, nil)
	}

	category, err := ctrl.Options.UseCases.Category.IsCategoryExist(categoryID)
	if category == nil {
		return helpers.Response(c, http.StatusNotFound, []string{"Category does not exist"})
	}
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
	}

	newProduct, err := ctrl.Options.UseCases.Product.CreateProduct(&req)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
	}
	return helpers.StandardResponse(c, http.StatusCreated, []string{"Product created successfully"}, newProduct, nil)
}

func (ctrl *productController) GetAllProduct(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	categoryID, _ := strconv.Atoi(c.QueryParam("category_id"))

	q := c.QueryParam("q")
	includeDeleted := false
	if c.QueryParam("include_deleted") == "true" {
		includeDeleted = true
	}
	sortBy := c.QueryParam("sort_by")
	sortOrder := c.QueryParam("sort_order")

	mapParams := map[string]any{
		"page":            page,
		"page_size":       pageSize,
		"category_id":     categoryID,
		"q":               q,
		"include_deleted": includeDeleted,
		"sort_by":         sortBy,
		"sort_order":      sortOrder,
	}

	schema := validet.NewSchema(
		mapParams,
		map[string]validet.Rule{
			"page": validet.Numeric[int]{Required: false, Custom: func(v int, path validet.PathKey, lookup validet.Lookup) error {
				if v != 0 && v < 1 {
					return customerror.NewBadRequestError("page must be minimum of 1")
				}
				return nil
			}},
			"page_size": validet.Numeric[int]{Required: false, Custom: func(v int, path validet.PathKey, lookup validet.Lookup) error {
				if v != 0 && v < 1 {
					return customerror.NewBadRequestError("page_size must be minimum of 1 or maximum 100")
				}
				return nil
			}},
			"category_id": validet.Numeric[int]{Required: false, Custom: func(v int, path validet.PathKey, lookup validet.Lookup) error {
				if v != 0 && v < 1 {
					return customerror.NewBadRequestError("category_id must be minimum of 1")
				}
				return nil
			}},
			"q":               validet.String{Required: false, Max: 100},
			"include_deleted": validet.Boolean{Required: false},
			"sort_by": validet.String{Required: false, Max: 50, Custom: func(v string, path validet.PathKey, lookup validet.Lookup) error {
				if v != "" && v != "id" && v != "name" && v != "price" && v != "stock" && v != "category_id" && v != "created_at" && v != "updated_at" {
					return customerror.NewBadRequestError("Invalid sort by field")
				}
				return nil
			}},
			"sort_order": validet.String{Required: false, Max: 10, Custom: func(v string, path validet.PathKey, lookup validet.Lookup) error {
				if v != "" && v != "asc" && v != "desc" {
					return customerror.NewBadRequestError("Invalid sort order, must be 'asc' or 'desc")
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

	product, total, err := ctrl.Options.UseCases.Product.GetAllProduct(page, pageSize, q, categoryID, includeDeleted, sortBy, sortOrder)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
	}

	if len(*product) == 0 {
		return helpers.Response(c, http.StatusNotFound, []string{"No product found"})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Product fetched successfully"}, map[string]interface{}{
		"total":   total,
		"product": product,
	}, nil)
}

func (ctrl *productController) UpdateProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var req models.Product

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
	if req.Price != 0 {
		mapReq["price"] = req.Price
	}
	if req.Stock != 0 {
		mapReq["stock"] = req.Stock
	}
	if req.CategoryID != nil {
		category, err := ctrl.Options.UseCases.Category.IsCategoryExist(*req.CategoryID)
		if category == nil {
			return helpers.Response(c, http.StatusNotFound, []string{"Category does not exist"})
		}
		if err != nil {
			return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
		}
		mapReq["category_id"] = *req.CategoryID
	}

	if req.ImageURL != nil {
		mapReq["image_url"] = *req.ImageURL
	}

	if len(mapReq) == 0 {
		return helpers.Response(c, http.StatusBadRequest, []string{"No fields to update"})
	}
	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"name":        validet.String{Required: false, Max: 100},
			"description": validet.String{Required: false, Max: 500},
			"price": validet.Numeric[float64]{Required: false, Custom: func(v float64, path validet.PathKey, lookup validet.Lookup) error {
				if v != 0 && v < 1 {
					return customerror.NewBadRequestError("price must be minimum of 1")
				}
				return nil
			}},
			"stock": validet.Numeric[int]{Required: false, Custom: func(v int, path validet.PathKey, lookup validet.Lookup) error {
				if v != 0 && v < 1 {
					return customerror.NewBadRequestError("stock must be minimum of 1")
				}
				return nil
			}},
			"category_id": validet.Numeric[int]{Required: false, Custom: func(v int, path validet.PathKey, lookup validet.Lookup) error {
				if v != 0 && v < 1 {
					return customerror.NewBadRequestError("category_id must be minimum of 1")
				}
				return nil
			}},
			"image_url": validet.String{Required: false, Max: 255},
		},
		validet.Options{},
	)
	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, errorBags.Errors, nil, nil)
	}

	product, err := ctrl.Options.UseCases.Product.UpdateProduct(id, mapReq)

	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
	}
	return helpers.StandardResponse(c, http.StatusOK, []string{"Product updated successfully"}, product, nil)
}

func (ctrl *productController) SoftDeleteProduct(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{"Invalid or empty product ID"})
	}

	mapReq := make(map[string]any)
	if id != 0 && id > 0 {
		mapReq["id"] = id
	}

	if len(mapReq) == 0 {
		return helpers.Response(c, http.StatusBadRequest, []string{"Product ID is required"})
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

	if err := ctrl.Options.UseCases.Product.SoftDeleteProduct(id); err != nil {
		return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
	}
	return helpers.StandardResponse(c, http.StatusOK, []string{"Product deleted successfully"}, nil, nil)
}
