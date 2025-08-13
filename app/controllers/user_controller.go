package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ezartsh/inrequest"
	"github.com/ezartsh/validet"
	"github.com/labstack/echo/v4"

	"monitoring-service/app/helpers"
	"monitoring-service/app/models"
	customerror "monitoring-service/pkg/customerror"
)

type userController controller

type UserControllerInterface interface {
	GetAllUsers(c echo.Context) error
	GetUserByID(c echo.Context) error
	Register(c echo.Context) error
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

// Custom email validator
func validateEmail(v string, path validet.PathKey, lookup validet.Lookup) error {
	if !strings.Contains(v, "@") || !strings.Contains(v, ".") {
		return errors.New("must be a valid email address")
	}
	return nil
}

// Register godoc
// @Summary Register a new shopper
// @Description Register a new user with shopper role
// @Tags public
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration data"
// @Success 201 {object} models.Response "Successfully registered"
// @Failure 400 {object} models.BasicResponse "Validation error"
// @Failure 409 {object} models.BasicResponse "Email already exists"
// @Failure 500 {object} models.BasicResponse "Internal server error"
// @Router /api/v1/public/register [post]
func (ctrl *userController) Register(c echo.Context) error {
	var request models.RegisterRequest

	// Bind request
	if err := c.Bind(&request); err != nil {
		return helpers.Response(c, http.StatusBadRequest, []string{err.Error()})
	}

	// Safe logging (no PII)
	if len(request.Email) >= 3 {
		log.Printf("Registration attempt for email starting with: %s", request.Email[:3]+"***")
	} else {
		log.Printf("Registration attempt with short email: %s", request.Email)
	}

	// Validate request
	mapReq := make(map[string]any)
	mapReq["name"] = request.Name
	mapReq["email"] = request.Email
	mapReq["password"] = request.Password

	schema := validet.NewSchema(
		mapReq,
		map[string]validet.Rule{
			"name":     validet.String{Required: true},
			"email":    validet.String{Required: true, Custom: validateEmail},
			"password": validet.String{Required: true, Min: 6},
		},
		validet.Options{},
	)

	errorBags, err := schema.Validate()
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	// Perbaikan di sini - cek apakah ada error validasi
	if len(errorBags.Errors) > 0 {
		var errorMessages []string
		for _, errs := range errorBags.Errors {
			for _, msg := range errs {
				errorMessages = append(errorMessages, msg)
			}
		}
		return helpers.StandardResponse(c, http.StatusBadRequest, errorMessages, nil, nil)
	}

	token, err := ctrl.Options.UseCases.User.Register(request)
	if err != nil {
		statusCode := customerror.GetStatusCode(err)
		return helpers.Response(c, statusCode, []string{err.Error()})
	}

	return helpers.StandardResponse(c, http.StatusCreated, []string{"User registered successfully"}, map[string]string{
		"token": token,
	}, nil)
}
