package controllers

import (
	"monitoring-service/app/helpers"
	"monitoring-service/pkg/customerror"
	"net/http"

	"github.com/labstack/echo/v4"
)

type categoryControllers controller

type CategoryControllersInterface interface {
	GetAll(ctx echo.Context) error
}

func (c *categoryControllers) GetAll(ctx echo.Context) error {
	response, err := c.Options.UseCases.Category.GetAll()
	if err != nil {
		return helpers.StandardResponse(ctx, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}
	return helpers.StandardResponse(ctx, http.StatusOK, nil, response, nil)
}
