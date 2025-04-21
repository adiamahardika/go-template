package controllers

import (
	"monitoring-service/app/helpers"
	"monitoring-service/pkg/customerror"
	"net/http"

	"github.com/labstack/echo/v4"
)

type priorityControllers controller

type PriorityControllersInterface interface {
	GetAll(ctx echo.Context) error
}

func (p *priorityControllers) GetAll(ctx echo.Context) error {
	response, err := p.Options.UseCases.Priority.GetAll()
	if err != nil {
		return helpers.StandardResponse(ctx, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}
	return helpers.StandardResponse(ctx, http.StatusOK, nil, response, nil)
}
