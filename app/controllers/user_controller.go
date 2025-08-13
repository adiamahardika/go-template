package controllers

import (
	"errors"
	"monitoring-service/app/helpers"
	"monitoring-service/app/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/ezartsh/inrequest"
	"github.com/ezartsh/validet"
	"github.com/labstack/echo/v4"

	customerror "monitoring-service/pkg/customerror"
)

type userController controller

type UserControllerInterface interface {
	GetAllUsers(c echo.Context) error
	GetUserByID(c echo.Context) error
	Login(c echo.Context) error // Tambahkan method baru
}

func (ctrl *userController) GetAllUsers(c echo.Context) error {
	var (
		request    models.GetUsersRequest
		users      []models.UserResponse
		pagination models.Pagination
		err        error
	)

	queryReq := inrequest.Query(c.Request())
	mapReq := queryReq.ToMap()

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		err := customerror.NewBadRequestError(err.Error())
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), errorBags.Errors, nil, nil)
	}

	err = queryReq.ToBind(&request)
	if err != nil {
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}

	// Get users from usecase
	users, pagination, err = ctrl.Options.UseCases.User.GetAllUsers(request)
	if err != nil {
		return helpers.Response(c, http.StatusInternalServerError, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Users retrieved successfully"}, users, &pagination)
}

func (ctrl *userController) GetUserByID(c echo.Context) error {
	var (
		user *models.UserResponse
		err  error
		id   string
	)

	id = c.Param("id")

	mapReq := make(map[string]any)
	mapReq["id"] = id

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"id": validet.String{Required: true, Custom: func(v string, path validet.PathKey, lookup validet.Lookup) error {
				if _, err := strconv.Atoi(v); err != nil {
					return customerror.NewBadRequestError("Invalid user ID format")
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

	userID, _ := strconv.Atoi(id)

	// Get user from usecase
	user, err = ctrl.Options.UseCases.User.GetUserByID(userID)
	if err != nil {
		return helpers.Response(c, customerror.GetStatusCode(err), []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"User retrieved successfully"}, user, nil)
}

func (ctrl *userController) Login(c echo.Context) error {
	var request models.LoginRequest

	// Bind request
	if err := c.Bind(&request); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	// Konversi request ke map untuk validasi
	requestMap := map[string]any{
		"email":    request.Email,
		"password": request.Password,
	}

	// Validasi
	schema := validet.NewSchema(
		requestMap, // Gunakan map bukan struct langsung
		map[string]validet.Rule{
			"email": validet.String{
				Required: true,
				Custom: func(v string, path validet.PathKey, lookup validet.Lookup) error {
					// Tambahkan validasi email custom jika diperlukan
					if !strings.Contains(v, "@") {
						return errors.New("invalid email format")
					}
					return nil
				},
			},
			"password": validet.String{
				Required: true,
				Custom: func(v string, path validet.PathKey, lookup validet.Lookup) error {
					// Validasi panjang password minimal 6 karakter
					if len(v) < 6 {
						return errors.New("password must be at least 6 characters")
					}
					return nil
				},
			},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, errorBags.Errors, nil, nil)
	}

	// Panggil use case
	response, err := ctrl.Options.UseCases.User.Login(request)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid email or password" {
			statusCode = http.StatusUnauthorized
		}
		return helpers.StandardResponse(c, statusCode, []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Login successful"}, response, nil)
}
