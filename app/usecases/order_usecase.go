package usecases

import (
	"context"
	"monitoring-service/app/models"
	"monitoring-service/app/repositories"
)

type OrderUsecaseInterface interface {
	GetOrderHistory(ctx context.Context, userID uint, page, pageSize int) ([]models.Order, int64, error)
	GetOrderDetail(ctx context.Context, userID, orderID uint) (*models.Order, error)
}

type orderUsecase struct {
	repo repositories.OrderRepositoryInterface
}

func NewOrderUsecase(repo repositories.OrderRepositoryInterface) OrderUsecaseInterface {
	return &orderUsecase{repo}
}

func (u *orderUsecase) GetOrderHistory(ctx context.Context, userID uint, page, pageSize int) ([]models.Order, int64, error) {
	return u.repo.GetOrdersByUserID(ctx, userID, page, pageSize)
}

func (u *orderUsecase) GetOrderDetail(ctx context.Context, userID, orderID uint) (*models.Order, error) {
	return u.repo.GetOrderDetailByID(ctx, userID, orderID)
}
