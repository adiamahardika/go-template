package controllers

import (
	"monitoring-service/app/helpers"
	"monitoring-service/app/models"
	"net/http"

	"strings"

	"github.com/ezartsh/inrequest"
	"github.com/ezartsh/validet"
	"github.com/labstack/echo/v4"
)

type authController controller

type AuthControllerInterface interface {
	Login(c echo.Context) error
	Register(c echo.Context) error
	IsAdminJWT(next echo.HandlerFunc) echo.HandlerFunc
}

func (ctrl *authController) Login(c echo.Context) error {
	var (
		reqBody models.LoginRequest
		resBody *models.AuthResponse
		err     error
	)

	req, err := inrequest.Json(c.Request())
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	mapReq := req.ToMap()
	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"email":    validet.String{Required: true, Email: true},
			"password": validet.String{Required: true},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, errorBags.Errors, nil, nil)
	}

	err = req.ToBind(&reqBody)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	resBody, err = ctrl.Options.UseCases.Auth.Login(c.Request().Context(), reqBody)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Login successful"}, resBody, nil)
}

func (ctrl *authController) Register(c echo.Context) error {
	var (
		reqBody models.RegisterRequest
		resBody *models.AuthResponse
		err     error
	)

	req, err := inrequest.Json(c.Request())
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	mapReq := req.ToMap()
	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"name":     validet.String{Required: true},
			"email":    validet.String{Required: true, Email: true},
			"password": validet.String{Required: true},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, errorBags.Errors, nil, nil)
	}

	err = req.ToBind(&reqBody)
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	resBody, err = ctrl.Options.UseCases.Auth.Register(c.Request().Context(), reqBody)
	if err != nil {
		// Check for specific errors and return appropriate status codes
		switch err.Error() {
		case "email already registered":
			return helpers.StandardResponse(c, http.StatusConflict, []string{err.Error()}, nil, nil)
		case "name is required", "email is required", "invalid email format",
			"password must be at least 8 characters long",
			"password must contain at least one letter",
			"password must contain at least one number":
			return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
		default:
			return helpers.StandardResponse(c, http.StatusInternalServerError, []string{err.Error()}, nil, nil)
		}
	}

	return helpers.StandardResponse(c, http.StatusCreated, []string{"Registration successful"}, resBody, nil)
}

func (ctrl *authController) IsAdminJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := c.Request().Header.Get("Authorization")
		if tokenStr == "" {
			return helpers.StandardResponse(c, http.StatusUnauthorized, []string{"Missing Authorization header"}, nil, nil)
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		if tokenStr == "" {
			return helpers.StandardResponse(c, http.StatusUnauthorized, []string{"Invalid Authorization header format"}, nil, nil)
		}

		err := ctrl.Options.UseCases.Auth.IsAdminJWT(tokenStr)
		if err != nil {
			return helpers.StandardResponse(c, http.StatusForbidden, []string{err.Error()}, nil, nil)
		}

		// ✅ Pass to next handler if check passes
		return next(c)
	}
}
