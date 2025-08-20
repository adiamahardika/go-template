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