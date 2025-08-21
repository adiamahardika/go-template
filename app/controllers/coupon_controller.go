package controllers

import (
	"monitoring-service/app/helpers"
	"monitoring-service/app/models/dto"
	"monitoring-service/pkg/customerror"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type couponController controller

type CouponControllerInterface interface {
	CreateCoupon(c echo.Context) error
	GetCoupons(c echo.Context) error
	GetCouponByID(c echo.Context) error
	UpdateCoupon(c echo.Context) error
	DeleteCoupon(c echo.Context) error
}

func (ctrl *couponController) CreateCoupon(c echo.Context) error {
	var request dto.CouponRequest
	if err := c.Bind(&request); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	response, err := ctrl.Options.UseCases.Coupon.CreateCoupon(request)
	if err != nil {
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusCreated, []string{"Coupon created successfully"}, response, nil)
}

func (ctrl *couponController) GetCoupons(c echo.Context) error {
	var request dto.GetCouponsRequest
	if err := c.Bind(&request); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	responses, pagination, err := ctrl.Options.UseCases.Coupon.GetCoupons(request)
	if err != nil {
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Coupons retrieved successfully"}, responses, &pagination)
}

func (ctrl *couponController) GetCouponByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"Invalid coupon ID"}, nil, nil)
	}

	response, err := ctrl.Options.UseCases.Coupon.GetCouponByID(id)
	if err != nil {
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Coupon retrieved successfully"}, response, nil)
}

func (ctrl *couponController) UpdateCoupon(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"Invalid coupon ID"}, nil, nil)
	}

	var request dto.CouponRequest
	if err := c.Bind(&request); err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{err.Error()}, nil, nil)
	}

	response, err := ctrl.Options.UseCases.Coupon.UpdateCoupon(id, request)
	if err != nil {
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Coupon updated successfully"}, response, nil)
}

func (ctrl *couponController) DeleteCoupon(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return helpers.StandardResponse(c, http.StatusBadRequest, []string{"Invalid coupon ID"}, nil, nil)
	}

	err = ctrl.Options.UseCases.Coupon.DeleteCoupon(id)
	if err != nil {
		return helpers.StandardResponse(c, customerror.GetStatusCode(err), []string{err.Error()}, nil, nil)
	}

	return helpers.StandardResponse(c, http.StatusOK, []string{"Coupon deleted successfully"}, nil, nil)
}