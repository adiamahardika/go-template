package controllers

import (
	"monitoring-service/app/helpers"
	"monitoring-service/pkg/customerror"
	"net/http"

	"github.com/labstack/echo/v4"
)

type projectController controller

type ProjectControllersInterface interface {
	GetAll(ctx echo.Context) error
}

func (t *projectController) GetAll(ctx echo.Context) error {
	response, err := t.Options.UseCases.Project.GetAll()
	if err != nil {
		return helpers.StandardResponse(ctx, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}
	return helpers.StandardResponse(ctx, http.StatusOK, nil, response, nil)
}
