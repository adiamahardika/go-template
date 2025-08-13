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
		product.Available = "out_of_stock"
	}
	if product.Stock > 0 {
		product.Available = "in_stock"
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Product retrieved successfully"}, product, nil)
}
