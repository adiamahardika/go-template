package usecases

import (
	"errors"
	"monitoring-service/app/models"

	"gorm.io/gorm"
)

type orderUsecase usecase

type OrderUsecaseInterface interface {
	ProcessOrderCheckout(orderID int) error
	CancelOrder(orderID int) error
}

func (u *orderUsecase) ProcessOrderCheckout(orderID int) error {
	// Get order with items untuk validasi awal
	order, err := u.Options.Repository.Order.GetOrderWithItems(orderID)
	if err != nil {
		return err
	}

	// Handle case where order status is nil
	if order.Status == nil {
		return errors.New("order status is not set")
	}

	// Check if order is eligible for stock reduction
	eligibleStatuses := []string{"pending", "paid", "processing"}
	if !contains(eligibleStatuses, order.Status) {
		return errors.New("order status not eligible for checkout")
	}

	// Gunakan transaction method dari repository
	return u.Options.Repository.Order.ProcessCheckoutTransaction(orderID, func(tx *gorm.DB) error {
		// Get order items within transaction
		var orderItems []models.OrderItem
		if err := tx.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
			return err
		}

		// Check stock and decrement
		for _, item := range orderItems {
			if item.ProductID == nil {
				return errors.New("product ID is nil in order item")
			}

			// Get product within transaction
			var product models.Product
			if err := tx.Where("id = ?", item.ProductID).First(&product).Error; err != nil {
				return err
			}

			if product.Stock < item.Quantity {
				return errors.New("insufficient stock for product: " + product.Name)
			}

			// Decrement stock
			if err := tx.Model(&models.Product{}).
				Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
				Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// Update order status
		status := "processing"
		return tx.Model(&models.Order{}).
			Where("id = ?", orderID).
			Update("status", &status).Error
	})
}

func (u *orderUsecase) CancelOrder(orderID int) error {
	// Get order with items untuk validasi awal
	order, err := u.Options.Repository.Order.GetOrderWithItems(orderID)
	if err != nil {
		return err
	}

	// Idempotency check
	if order.Status != nil && *order.Status == "cancelled" {
		return errors.New("order already cancelled")
	}

	// Handle case where order status is nil
	if order.Status == nil {
		return errors.New("order status is not set")
	}

	// Check if order was in eligible status for stock restoration
	eligibleStatuses := []string{"pending", "paid", "processing"}
	if !contains(eligibleStatuses, order.Status) {
		return errors.New("order status not eligible for cancellation with stock restoration")
	}

	// Gunakan transaction method dari repository
	return u.Options.Repository.Order.CancelOrderTransaction(orderID, func(tx *gorm.DB) error {
		// Get order items within transaction
		var orderItems []models.OrderItem
		if err := tx.Where("order_id = ?", orderID).Find(&orderItems).Error; err != nil {
			return err
		}

		// Restore stock
		for _, item := range orderItems {
			if item.ProductID == nil {
				return errors.New("product ID is nil in order item")
			}

			// Increment stock
			if err := tx.Model(&models.Product{}).
				Where("id = ?", item.ProductID).
				Update("stock", gorm.Expr("stock + ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// Update order status to cancelled
		cancelledStatus := "cancelled"
		return tx.Model(&models.Order{}).
			Where("id = ?", orderID).
			Update("status", &cancelledStatus).Error
	})
}

// Helper function
func contains(slice []string, item *string) bool {
	if item == nil {
		return false
	}
	for _, s := range slice {
		if s == *item {
			return true
		}
	}
	return false
}