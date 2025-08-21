package repositories

import (
	"monitoring-service/app/models"

	"gorm.io/gorm"
)

type couponRepository repository

type CouponRepositoryInterface interface {
	CreateCoupon(coupon models.Coupon) (*models.Coupon, error)
	GetCoupons(limit, offset int, code string, active bool) ([]models.Coupon, int64, error)
	GetCouponByID(id int) (*models.Coupon, error)
	UpdateCoupon(coupon models.Coupon) (*models.Coupon, error)
	DeleteCoupon(id int) error
	GetCouponByCode(code string) (*models.Coupon, error)
}

func (r *couponRepository) CreateCoupon(coupon models.Coupon) (*models.Coupon, error) {
	err := r.Options.Postgres.Create(&coupon).Error
	return &coupon, err
}

func (r *couponRepository) GetCoupons(limit, offset int, code string, active bool) ([]models.Coupon, int64, error) {
	var coupons []models.Coupon
	var total int64

	query := r.Options.Postgres.Model(&models.Coupon{})

	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}

	if active {
		query = query.Where("expired_at > NOW() OR expired_at IS NULL")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&coupons).Error; err != nil {
		return nil, 0, err
	}	

	return coupons, total, nil
}

func (r *couponRepository) GetCouponByID(id int) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.Options.Postgres.First(&coupon, id).Error
	return &coupon, err
}

func (r *couponRepository) UpdateCoupon(coupon models.Coupon) (*models.Coupon, error) {
	err := r.Options.Postgres.Save(&coupon).Error
	return &coupon, err
}

func (r *couponRepository) DeleteCoupon(id int) error {
	return r.Options.Postgres.Delete(&models.Coupon{}, id).Error
}

// func (r *couponRepository) DeleteCoupon(id int) error {
// 	return r.Options.Postgres.Model(&models.Coupon{}).Where("id = ?", id).Update("deleted_at", time.Now()).Error
// }

func (r *couponRepository) GetCouponByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.Options.Postgres.Where("code = ?", code).First(&coupon).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &coupon, nil
}