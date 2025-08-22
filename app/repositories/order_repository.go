package repositories

import (
	"context"
	"monitoring-service/app/models"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	GetOrdersByUserID(ctx context.Context, userID uint, page, pageSize int) ([]models.Order, int64, error)
	GetOrderDetailByID(ctx context.Context, userID, orderID uint) (*models.Order, error)
}

type orderRepository struct {
	*repository
}

func NewOrderRepository(repo *repository) OrderRepositoryInterface {
	return &orderRepository{repo}
}

func (r *orderRepository) GetOrdersByUserID(ctx context.Context, userID uint, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	offset := (page - 1) * pageSize

	query := r.Options.Postgres.WithContext(ctx).
		Model(&models.Order{}).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	query.Count(&total)
	err := query.Preload("Coupon").Preload("Payments.PaymentMethod").Preload("Shipments.ShippingMethod").
		Limit(pageSize).Offset(offset).
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (r *orderRepository) GetOrderDetailByID(ctx context.Context, userID, orderID uint) (*models.Order, error) {
	var order models.Order
	err := r.Options.Postgres.WithContext(ctx).
		Model(&models.Order{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		// preload order items and their product using Unscoped() so we can read product name even if soft-deleted
		Preload("OrderItems.Product", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("Coupon").
		Preload("Payments.PaymentMethod").
		Preload("Shipments.ShippingMethod").
		First(&order).Error
	if err != nil {
		return nil, err
	}

	// explicitly load order items and preload product to guarantee snapshot product name is available
	var items []models.OrderItem
	if err := r.Options.Postgres.WithContext(ctx).
			Model(&models.OrderItem{}).
			Where("order_id = ?", order.ID).
			Preload("Product", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
			Find(&items).Error; err == nil {
		order.OrderItems = items
	}

	// If still empty, try to build order items from the user's cart (fallback)
	if order.OrderItems == nil || len(order.OrderItems) == 0 {
		// resolve user id
		var userID int
		if order.UserID != nil {
			userID = *order.UserID
		}
		// find cart for user
		var cart models.Cart
		if userID != 0 {
			if err := r.Options.Postgres.WithContext(ctx).
				Model(&models.Cart{}).
				Where("user_id = ?", userID).
				First(&cart).Error; err == nil {
				// load cart items with product
				var cartItems []models.CartItem
				if err := r.Options.Postgres.WithContext(ctx).
					Model(&models.CartItem{}).
					Where("cart_id = ?", cart.ID).
					Preload("Product").
					Find(&cartItems).Error; err == nil {
					// synthesize order items from cart items
					var synthesized []models.OrderItem
					for _, ci := range cartItems {
						if ci.Product == nil || ci.Product.ID == 0 {
							continue
						}
						pid := ci.Product.ID
						oid := order.ID
						synthesized = append(synthesized, models.OrderItem{
							OrderID:   &oid,
							ProductID: &pid,
							Quantity:  ci.Quantity,
							Price:     ci.Product.Price,
							Product:   ci.Product,
						})
					}
					order.OrderItems = synthesized
				}
			}
		}
	}

	return &order, nil
}
