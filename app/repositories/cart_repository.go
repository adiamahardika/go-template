package repositories

import (
	"errors"
	"monitoring-service/app/models"

	"gorm.io/gorm"
)

type CartRepositoryInterface interface {
	GetCartByUserID(userID int) (models.Cart, error)
	GetCartItemsByCartID(cartID int) ([]models.CartItem, error)
	CreateCart(userID int) (models.Cart, error)
	GetProductByID(productID int) (*models.Product, error)
	GetCartItemByCartIDAndProductID(cartID, productID int) (*models.CartItem, error)
	AddCartItem(item *models.CartItem) error
	UpdateCartItem(item *models.CartItem) error
	DeleteCartItemByCartIDAndProductID(cartID, productID int) error
	GetCartItemByID(id int) (*models.CartItem, error)

	// New methods for coupon functionality
	ApplyCoupon(cartID int, couponID *int) error
	RemoveCoupon(cartID int) error
	GetCartWithCoupon(cartID int) (*models.Cart, error)
	GetCouponByCode(code string) (*models.Coupon, error)
	GetCouponByID(id int) (*models.Coupon, error)
	ValidateCoupon(couponID int) (bool, error)
}

type cartRepository struct {
	Options Options
}

func (r *cartRepository) GetCartByUserID(userID int) (models.Cart, error) {
	var cart models.Cart
	err := r.Options.Postgres.Where("user_id = ?", userID).First(&cart).Error
	return cart, err
}

func (r *cartRepository) GetCartItemsByCartID(cartID int) ([]models.CartItem, error) {
	var items []models.CartItem
	err := r.Options.Postgres.
		Preload("Product").
		Where("cart_id = ?", cartID).
		Find(&items).Error
	return items, err
}

func (r *cartRepository) CreateCart(userID int) (models.Cart, error) {
	cart := models.Cart{UserID: &userID}
	if err := r.Options.Postgres.Create(&cart).Error; err != nil {
		return models.Cart{}, err
	}
	return cart, nil
}

func (r *cartRepository) GetProductByID(productID int) (*models.Product, error) {
	var p models.Product
	err := r.Options.Postgres.
		Where("id = ? AND deleted_at IS NULL", productID).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *cartRepository) GetCartItemByCartIDAndProductID(cartID, productID int) (*models.CartItem, error) {
	var item models.CartItem
	err := r.Options.Postgres.
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *cartRepository) AddCartItem(item *models.CartItem) error {
	return r.Options.Postgres.Create(item).Error
}

func (r *cartRepository) UpdateCartItem(item *models.CartItem) error {
	return r.Options.Postgres.Save(item).Error
}

func (r *cartRepository) DeleteCartItemByCartIDAndProductID(cartID, productID int) error {
	return r.Options.Postgres.
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Delete(&models.CartItem{}).Error
}

func (r *cartRepository) GetCartItemByID(id int) (*models.CartItem, error) {
	var item models.CartItem
	err := r.Options.Postgres.
		Preload("Product").
		First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	return &item, err
}

// New methods for coupon functionality
func (r *cartRepository) ApplyCoupon(cartID int, couponID *int) error {
	return r.Options.Postgres.
		Model(&models.Cart{}).
		Where("id = ?", cartID).
		Update("coupon_id", couponID).Error
}

func (r *cartRepository) RemoveCoupon(cartID int) error {
	return r.Options.Postgres.
		Model(&models.Cart{}).
		Where("id = ?", cartID).
		Update("coupon_id", nil).Error
}

func (r *cartRepository) GetCartWithCoupon(cartID int) (*models.Cart, error) {
	var cart models.Cart
	err := r.Options.Postgres.
		Preload("Coupon").
		Preload("Items").
		Preload("Items.Product").
		First(&cart, cartID).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) GetCouponByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.Options.Postgres.
		Where("code = ? AND deleted_at IS NULL", code).
		First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *cartRepository) GetCouponByID(id int) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.Options.Postgres.
		Where("id = ? AND deleted_at IS NULL", id).
		First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *cartRepository) ValidateCoupon(couponID int) (bool, error) {
	var coupon models.Coupon
	err := r.Options.Postgres.
		Where("id = ? AND deleted_at IS NULL AND (expired_at IS NULL OR expired_at > NOW())", couponID).
		First(&coupon).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
