package usecases

import (
	"context"
	"errors"
	"fmt"
	"monitoring-service/app/models"
	
	"gorm.io/gorm"
)

type OrderUsecaseInterface interface {
	Checkout(ctx context.Context, userID int, shippingMethodID int, couponCode string) (*models.Order, error)
}

type orderUsecase struct {
	usecase
}

func (u *orderUsecase) Checkout(ctx context.Context, userID int, shippingMethodID int, couponCode string) (*models.Order, error) {
	tx := u.Options.Repository.Order.BeginTransaction()
	if tx.Error != nil {
		return nil, errors.New("failed to start transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Get cart with items
	cart, err := u.Options.Repository.Order.GetCartWithItems(userID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart not found")
		}
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// 2. Validate cart has items
	if len(cart.Items) == 0 {
		tx.Rollback()
		return nil, errors.New("cart is empty")
	}

	// 3. Validate and process each product
	var orderItems []models.OrderItem
	subtotal := 0.0

	for _, item := range cart.Items {
		if item.Product == nil {
			tx.Rollback()
			return nil, errors.New("product not found")
		}

		if item.ProductID == nil {
			tx.Rollback()
			return nil, errors.New("product ID is nil")
		}

		// Lock product for update
		product, err := u.Options.Repository.Order.GetProductWithLock(tx, *item.ProductID)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("product not available")
			}
			return nil, fmt.Errorf("failed to get product: %w", err)
		}

		// Validate stock
		if product.Stock < item.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient stock for product: %s. Available: %d, Requested: %d",
				product.Name, product.Stock, item.Quantity)
		}

		// Update stock
		newStock := product.Stock - item.Quantity
		if err := u.Options.Repository.Order.UpdateProductStock(tx, product.ID, newStock); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}

		// Create order item
		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
		}
		orderItems = append(orderItems, orderItem)
		subtotal += product.Price * float64(item.Quantity)
	}

	// 4. Validate and apply coupon if provided
	var coupon *models.Coupon
	var couponID *int
	discountAmount := 0.0

	if couponCode != "" {
		coupon, err = u.Options.Repository.Order.GetCouponByCode(couponCode)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("coupon not found")
			}
			return nil, fmt.Errorf("failed to get coupon: %w", err)
		}

		// Revalidate coupon
		isValid, err := u.Options.Repository.Order.ValidateCoupon(coupon.ID)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to validate coupon: %w", err)
		}
		if !isValid {
			tx.Rollback()
			return nil, errors.New("coupon is not valid or expired")
		}

		couponID = &coupon.ID
		discountAmount = u.calculateDiscount(subtotal, coupon)
	}

	// 5. Get shipping method
	shippingMethod, err := u.Options.Repository.Order.GetShippingMethodByID(shippingMethodID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shipping method not found")
		}
		return nil, fmt.Errorf("failed to get shipping method: %w", err)
	}

	// 6. Calculate total amount
	totalAmount := subtotal - discountAmount + shippingMethod.Cost
	if totalAmount < 0 {
		totalAmount = 0
	}

	// 7. Create order
	status := "pending"
	order := &models.Order{
		UserID:      &userID,
		CouponID:    couponID,
		TotalAmount: &totalAmount,
		Status:      &status,
	}

	if err := u.Options.Repository.Order.CreateOrder(tx, order); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// 8. Create order items with order ID
	for i := range orderItems {
		orderItems[i].OrderID = &order.ID
	}

	if err := u.Options.Repository.Order.CreateOrderItems(tx, orderItems); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create order items: %w", err)
	}

	// 9. Create shipment
	shipmentStatus := "pending"
	shipment := &models.Shipment{
		OrderID:          &order.ID,
		ShippingMethodID: &shippingMethod.ID,
		Status:           &shipmentStatus,
	}

	if err := u.Options.Repository.Order.CreateShipment(tx, shipment); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	// 10. Clear cart
	if err := u.Options.Repository.Order.DeleteCartItems(tx, cart.ID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to clear cart items: %w", err)
	}

	if err := u.Options.Repository.Order.DeleteCart(tx, cart.ID); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	// 11. Commit transaction
	if err := u.Options.Repository.Order.CommitTransaction(tx); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to complete checkout: %w", err)
	}

	return order, nil
}

func (u *orderUsecase) calculateDiscount(subtotal float64, coupon *models.Coupon) float64 {
	if coupon.DiscountPercent == nil {
		return 0
	}

	discount := subtotal * (*coupon.DiscountPercent / 100)

	if coupon.MaxDiscount != nil && discount > *coupon.MaxDiscount {
		return *coupon.MaxDiscount
	}

	return discount
}
