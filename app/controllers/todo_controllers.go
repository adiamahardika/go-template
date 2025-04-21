package controllers

import (
	"monitoring-service/app/helpers"
	"monitoring-service/pkg/customerror"
	"net/http"

	"github.com/labstack/echo/v4"
)

type todoControllers controller

type TodoControllersInterface interface {
	GetAll(ctx echo.Context) error
}

func (t *todoControllers) GetAll(ctx echo.Context) error {
	response, err := t.Options.UseCases.Todo.GetAll()
	if err != nil {
		return helpers.StandardResponse(ctx, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}
	return helpers.StandardResponse(ctx, http.StatusOK, nil, response, nil)
}
