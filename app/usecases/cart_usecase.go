package usecases

import (
	"context"
	"errors"
	"monitoring-service/app/models"
	"monitoring-service/app/models/dto"
	"time"

	"gorm.io/gorm"
)

type CartUsecaseInterface interface {
	GetCart(ctx context.Context, userID int) (models.Cart, []models.CartItem, error)
	GetCartItemsByUserID(ctx context.Context, userID int) ([]models.CartItem, error)
	AddCartItem(ctx context.Context, userID, productID, addQty int) (*models.CartItem, string, error)
	RemoveCartItem(ctx context.Context, userID, productID int) error
	UpdateCartItemQuantity(ctx context.Context, userID, cartItemID, newQty int) (*models.CartItem, error)
	ViewCart(ctx context.Context, userID int) (*dto.CartViewResponse, error)

	// New methods for coupon functionality
	ApplyCoupon(ctx context.Context, userID int, couponCode string) error
	RemoveCoupon(ctx context.Context, userID int) error
	CalculateCartTotal(ctx context.Context, userID int) (*dto.CartSummaryResponse, error)
}

type cartUsecase struct {
	*usecase
}

func (u *cartUsecase) AddCartItem(ctx context.Context, userID, productID, addQty int) (*models.CartItem, string, error) {
	if addQty <= 0 {
		return nil, "", errors.New("quantity must be >= 1")
	}

	// 1) Dapatkan/buat cart user
	cart, err := u.Options.Repository.Cart.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart, err = u.Options.Repository.Cart.CreateCart(userID)
			if err != nil {
				return nil, "", err
			}
		} else {
			return nil, "", err
		}
	}
	if cart.ID == 0 { // jaga-jaga
		cart, err = u.Options.Repository.Cart.CreateCart(userID)
		if err != nil {
			return nil, "", err
		}
	}

	// 2) Validasi produk
	product, err := u.Options.Repository.Cart.GetProductByID(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("product not found")
		}
		return nil, "", err
	}
	if product == nil {
		return nil, "", errors.New("product not found")
	}
	if product.DeletedAt != nil {
		return nil, "", errors.New("product is deleted")
	}
	if product.Stock <= 0 {
		return nil, "", errors.New("quantity exceeds stock")
	}

	// 3) Cek apakah item sudah ada
	existing, err := u.Options.Repository.Cart.GetCartItemByCartIDAndProductID(cart.ID, product.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", err
	}

	var message string

	if existing != nil && existing.ID != 0 {
		// Update quantity dengan batas stok
		newQty := existing.Quantity + addQty
		if newQty > product.Stock {
			newQty = product.Stock
			message = "quantity capped to product stock"
		}
		existing.Quantity = newQty

		if err := u.Options.Repository.Cart.UpdateCartItem(existing); err != nil {
			return nil, "", err
		}
		return existing, message, nil
	}

	// 4) Tambah item baru
	qty := addQty
	if qty > product.Stock {
		qty = product.Stock
		message = "quantity capped to product stock"
	}
	item := &models.CartItem{
		CartID:    &cart.ID,
		ProductID: &product.ID,
		Quantity:  qty,
	}
	if err := u.Options.Repository.Cart.AddCartItem(item); err != nil {
		return nil, "", err
	}
	return item, message, nil
}

// New methods for coupon functionality
func (u *cartUsecase) ApplyCoupon(ctx context.Context, userID int, couponCode string) error {
	// Get user's cart
	cart, err := u.Options.Repository.Cart.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create cart if not exists
			cart, err = u.Options.Repository.Cart.CreateCart(userID)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Get coupon by code
	coupon, err := u.Options.Repository.Cart.GetCouponByCode(couponCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("coupon not found")
		}
		return err
	}

	// Validate coupon
	valid, err := u.Options.Repository.Cart.ValidateCoupon(coupon.ID)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("coupon is not valid or expired")
	}

	// Apply coupon to cart
	return u.Options.Repository.Cart.ApplyCoupon(cart.ID, &coupon.ID)
}

func (u *cartUsecase) RemoveCoupon(ctx context.Context, userID int) error {
	cart, err := u.Options.Repository.Cart.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("cart not found")
		}
		return err
	}

	return u.Options.Repository.Cart.RemoveCoupon(cart.ID)
}

func (u *cartUsecase) CalculateCartTotal(ctx context.Context, userID int) (*dto.CartSummaryResponse, error) {
	cart, err := u.Options.Repository.Cart.GetCartWithCoupon(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return empty cart if not found
			return &dto.CartSummaryResponse{
				ID:        0,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Subtotal:  0,
				Discount:  0,
				Total:     0,
				Items:     []dto.CartItemView{},
			}, nil
		}
		return nil, err
	}

	// Calculate subtotal
	var subtotal float64 = 0
	items := make([]dto.CartItemView, 0)

	for _, item := range cart.Items {
		if item.Product != nil {
			itemSubtotal := float64(item.Quantity) * item.Product.Price
			subtotal += itemSubtotal

			items = append(items, dto.CartItemView{
				ProductID: item.Product.ID,
				Name:      item.Product.Name,
				Price:     item.Product.Price,
				Quantity:  item.Quantity,
				Subtotal:  itemSubtotal,
			})
		}
	}

	// Calculate discount
	var discount float64 = 0
	var couponCode *string = nil

	if cart.Coupon != nil && cart.Coupon.DiscountPercent != nil {
		discountValue := subtotal * (*cart.Coupon.DiscountPercent / 100)

		// Apply max discount if specified
		if cart.Coupon.MaxDiscount != nil && discountValue > *cart.Coupon.MaxDiscount {
			discount = *cart.Coupon.MaxDiscount
		} else {
			discount = discountValue
		}

		couponCode = &cart.Coupon.Code
	}

	total := subtotal - discount

	return &dto.CartSummaryResponse{
		ID:         cart.ID,
		CreatedAt:  cart.CreatedAt,
		UpdatedAt:  cart.UpdatedAt,
		CouponCode: couponCode,
		Items:      items,
		Subtotal:   subtotal,
		Discount:   discount,
		Total:      total,
	}, nil
}

// Update ViewCart to include coupon calculation
func (u *cartUsecase) ViewCart(ctx context.Context, userID int) (*dto.CartViewResponse, error) {
	cartSummary, err := u.CalculateCartTotal(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.CartViewResponse{
		ID:         cartSummary.ID,
		CreatedAt:  cartSummary.CreatedAt,
		UpdatedAt:  cartSummary.UpdatedAt,
		CouponCode: cartSummary.CouponCode,
		Discount:   cartSummary.Discount,
		Items:      cartSummary.Items,
		Subtotal:   cartSummary.Subtotal,
		Total:      cartSummary.Total,
	}, nil
}

// Existing methods remain unchanged...
func (u *cartUsecase) GetCart(ctx context.Context, userID int) (models.Cart, []models.CartItem, error) {
	cart, err := u.Options.Repository.Cart.GetCartByUserID(userID)
	if err != nil {
		return models.Cart{}, nil, err
	}

	items, err := u.Options.Repository.Cart.GetCartItemsByCartID(cart.ID)
	if err != nil {
		return cart, nil, err
	}

	return cart, items, nil
}

func (u *cartUsecase) GetCartItemsByUserID(ctx context.Context, userID int) ([]models.CartItem, error) {
	cart, err := u.Options.Repository.Cart.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}
	return u.Options.Repository.Cart.GetCartItemsByCartID(cart.ID)
}

func (u *cartUsecase) RemoveCartItem(ctx context.Context, userID, productID int) error {
	// Temukan cart milik user
	cart, err := u.Options.Repository.Cart.GetCartByUserID(userID)
	if err != nil {
		// jika not found, anggap sudah "terhapus" (idempotent)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if cart.ID == 0 {
		// tidak ada cart → nothing to delete
		return nil
	}

	// Hapus item berdasarkan cart_id & product_id
	return u.Options.Repository.Cart.DeleteCartItemByCartIDAndProductID(cart.ID, productID)
}

func (u *cartUsecase) UpdateCartItemQuantity(ctx context.Context, userID, cartItemID, newQty int) (*models.CartItem, error) {
	if newQty < 1 {
		return nil, errors.New("invalid quantity")
	}

	cart, err := u.Options.Repository.Cart.GetCartByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || cart.ID == 0 {
			return nil, errors.New("cart not found")
		}
		return nil, err
	}

	item, err := u.Options.Repository.Cart.GetCartItemByID(cartItemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart item not found")
		}
		return nil, err
	}
	if item == nil || item.ID == 0 || item.CartID == nil || *item.CartID != cart.ID {
		return nil, errors.New("cart item not found")
	}

	if item.ProductID == nil {
		return nil, errors.New("product not found")
	}
	product, err := u.Options.Repository.Cart.GetProductByID(*item.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	if product.DeletedAt != nil {
		return nil, errors.New("product is deleted")
	}

	// CAP ke stok terkini — tidak mengurangi stok produk
	qty := newQty
	if qty > product.Stock {
		qty = product.Stock
	}
	if qty < 1 {
		// kalau stok 0 sekarang, biarkan jadi 1? Atau 0?
		// Mengikuti aturan "must be ≥ 1", kita balikan error jika stok 0.
		return nil, errors.New("quantity exceeds stock")
	}

	item.Quantity = qty
	if err := u.Options.Repository.Cart.UpdateCartItem(item); err != nil {
		return nil, err
	}
	return item, nil
}
