package repositories

import (
	"errors"
	"monitoring-service/app/models"

	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	BeginTransaction() *gorm.DB
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error

	GetCartWithItems(userID int) (*models.Cart, error)
	GetProductWithLock(tx *gorm.DB, productID int) (*models.Product, error)
	UpdateProductStock(tx *gorm.DB, productID int, newStock int) error
	CreateOrder(tx *gorm.DB, order *models.Order) error
	CreateOrderItems(tx *gorm.DB, orderItems []models.OrderItem) error
	CreateShipment(tx *gorm.DB, shipment *models.Shipment) error
	DeleteCartItems(tx *gorm.DB, cartID int) error
	DeleteCart(tx *gorm.DB, cartID int) error
	GetShippingMethodByID(shippingMethodID int) (*models.ShippingMethod, error)
	GetCouponByCode(code string) (*models.Coupon, error)
	ValidateCoupon(couponID int) (bool, error)
}

type orderRepository struct {
	Options Options
}

func (r *orderRepository) BeginTransaction() *gorm.DB {
	return r.Options.Postgres.Begin()
}

func (r *orderRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *orderRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *orderRepository) GetCartWithItems(userID int) (*models.Cart, error) {
	var cart models.Cart
	err := r.Options.Postgres.
		Preload("Items").
		Preload("Items.Product").
		Preload("Coupon").
		Where("user_id = ?", userID).
		First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *orderRepository) GetProductWithLock(tx *gorm.DB, productID int) (*models.Product, error) {
	var product models.Product
	err := tx.
		Set("gorm:query_option", "FOR UPDATE").
		Where("id = ? AND deleted_at IS NULL", productID).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *orderRepository) UpdateProductStock(tx *gorm.DB, productID int, newStock int) error {
	return tx.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("stock", newStock).Error
}

func (r *orderRepository) CreateOrder(tx *gorm.DB, order *models.Order) error {
	// Jangan set ID manual, biarkan database yang generate
	if order.ID != 0 {
		order.ID = 0
	}
	return tx.Create(order).Error
}

func (r *orderRepository) CreateOrderItems(tx *gorm.DB, orderItems []models.OrderItem) error {
	if len(orderItems) == 0 {
		return nil
	}

	// Pastikan ID order items tidak di-set manual
	for i := range orderItems {
		orderItems[i].ID = 0
	}

	return tx.Create(&orderItems).Error
}

func (r *orderRepository) CreateShipment(tx *gorm.DB, shipment *models.Shipment) error {
	// Jangan set ID manual
	if shipment.ID != 0 {
		shipment.ID = 0
	}
	return tx.Create(shipment).Error
}

func (r *orderRepository) DeleteCartItems(tx *gorm.DB, cartID int) error {
	return tx.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

func (r *orderRepository) DeleteCart(tx *gorm.DB, cartID int) error {
	return tx.Delete(&models.Cart{}, cartID).Error
}

func (r *orderRepository) GetShippingMethodByID(shippingMethodID int) (*models.ShippingMethod, error) {
	var method models.ShippingMethod
	err := r.Options.Postgres.
		Where("id = ? AND deleted_at IS NULL", shippingMethodID).
		First(&method).Error
	if err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *orderRepository) GetCouponByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon
	err := r.Options.Postgres.
		Where("code = ? AND deleted_at IS NULL", code).
		First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func (r *orderRepository) ValidateCoupon(couponID int) (bool, error) {
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
