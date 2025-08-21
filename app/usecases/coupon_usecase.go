package usecases

import (
	"errors"
	"math"
	"monitoring-service/app/models"
	"monitoring-service/app/models/dto"
	"monitoring-service/pkg/customerror"

	"gorm.io/gorm"
)

type couponUsecase usecase

type CouponUsecaseInterface interface {
	CreateCoupon(req dto.CouponRequest) (*dto.CouponResponse, error)
	GetCoupons(req dto.GetCouponsRequest) ([]dto.CouponResponse, models.Pagination, error)
	GetCouponByID(id int) (*dto.CouponResponse, error)
	UpdateCoupon(id int, req dto.CouponRequest) (*dto.CouponResponse, error)
	DeleteCoupon(id int) error
}

func (u *couponUsecase) CreateCoupon(req dto.CouponRequest) (*dto.CouponResponse, error) {
	if req.DiscountPercent < 0 || req.DiscountPercent > 100 {
		return nil, customerror.NewBadRequestError("Discount percent must be between 0 and 100")
	}

	if req.MaxDiscount != nil && *req.MaxDiscount < 0 {
		return nil, customerror.NewBadRequestError("Max discount must be greater than or equal to 0")
	}

	existingCoupon, err := u.Options.Repository.Coupon.GetCouponByCode(req.Code)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existingCoupon != nil {
		return nil, customerror.NewConflictError("Coupon code already exists")
	}

	coupon := models.Coupon{
		Code:            req.Code,
		DiscountPercent: &req.DiscountPercent,
		MaxDiscount:     req.MaxDiscount,
		ExpiredAt:       req.ExpiredAt,
	}

	newCoupon, err := u.Options.Repository.Coupon.CreateCoupon(coupon)
	if err != nil {
		return nil, err
	}

	return &dto.CouponResponse{
		ID:              newCoupon.ID,
		Code:            newCoupon.Code,
		DiscountPercent: *newCoupon.DiscountPercent,
		MaxDiscount:     newCoupon.MaxDiscount,
		ExpiredAt:       newCoupon.ExpiredAt,
		CreatedAt:       newCoupon.CreatedAt,
		UpdatedAt:       newCoupon.UpdatedAt,
	}, nil
}

func (u *couponUsecase) GetCoupons(req dto.GetCouponsRequest) ([]dto.CouponResponse, models.Pagination, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	coupons, total, err := u.Options.Repository.Coupon.GetCoupons(req.Limit, offset, req.Code, req.Active)
	if err != nil {
		return nil, models.Pagination{}, err
	}

	var responses []dto.CouponResponse
	for _, coupon := range coupons {
		responses = append(responses, dto.CouponResponse{
			ID:              coupon.ID,
			Code:            coupon.Code,
			DiscountPercent: *coupon.DiscountPercent,
			MaxDiscount:     coupon.MaxDiscount,
			ExpiredAt:       coupon.ExpiredAt,
			CreatedAt:       coupon.CreatedAt,
			UpdatedAt:       coupon.UpdatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	pagination := models.Pagination{
		Page:      req.Page,
		PageSize:  req.Limit,
		Total:     int(total),
		TotalPage: totalPages,
	}

	return responses, pagination, nil
}

func (u *couponUsecase) GetCouponByID(id int) (*dto.CouponResponse, error) {
	coupon, err := u.Options.Repository.Coupon.GetCouponByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customerror.NewNotFoundError("Coupon not found")
		}
		return nil, err
	}

	return &dto.CouponResponse{
		ID:              coupon.ID,
		Code:            coupon.Code,
		DiscountPercent: *coupon.DiscountPercent,
		MaxDiscount:     coupon.MaxDiscount,
		ExpiredAt:       coupon.ExpiredAt,
		CreatedAt:       coupon.CreatedAt,
		UpdatedAt:       coupon.UpdatedAt,
	}, nil
}

func (u *couponUsecase) UpdateCoupon(id int, req dto.CouponRequest) (*dto.CouponResponse, error) {
	if req.DiscountPercent < 0 || req.DiscountPercent > 100 {
		return nil, customerror.NewBadRequestError("Discount percent must be between 0 and 100")
	}

	if req.MaxDiscount != nil && *req.MaxDiscount < 0 {
		return nil, customerror.NewBadRequestError("Max discount must be greater than or equal to 0")
	}

	coupon, err := u.Options.Repository.Coupon.GetCouponByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customerror.NewNotFoundError("Coupon not found")
		}
		return nil, err
	}

	if coupon.Code != req.Code {
		existingCoupon, err := u.Options.Repository.Coupon.GetCouponByCode(req.Code)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if existingCoupon != nil {
			return nil, customerror.NewConflictError("Coupon code already exists")
		}
	}

	coupon.Code = req.Code
	coupon.DiscountPercent = &req.DiscountPercent
	coupon.MaxDiscount = req.MaxDiscount
	coupon.ExpiredAt = req.ExpiredAt

	updatedCoupon, err := u.Options.Repository.Coupon.UpdateCoupon(*coupon)
	if err != nil {
		return nil, err
	}

	return &dto.CouponResponse{
		ID:              updatedCoupon.ID,
		Code:            updatedCoupon.Code,
		DiscountPercent: *updatedCoupon.DiscountPercent,
		MaxDiscount:     updatedCoupon.MaxDiscount,
		ExpiredAt:       updatedCoupon.ExpiredAt,
		CreatedAt:       updatedCoupon.CreatedAt,
		UpdatedAt:       updatedCoupon.UpdatedAt,
	}, nil
}

func (u *couponUsecase) DeleteCoupon(id int) error {
	_, err := u.Options.Repository.Coupon.GetCouponByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customerror.NewNotFoundError("Coupon not found")
		}
		return err
	}
	return u.Options.Repository.Coupon.DeleteCoupon(id)
}