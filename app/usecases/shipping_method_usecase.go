package usecases

import (
	"monitoring-service/app/models"
	"monitoring-service/app/repositories"
)

type ShippingMethodUsecaseInterface interface {
	GetAll() ([]models.ShippingMethod, error)
	GetByID(id int) (*models.ShippingMethod, error)
	GetCartQuote(userID int, shippingMethodID int) (*CartQuote, error)
}

type shippingMethodUsecase struct {
	repo     repositories.ShippingMethodRepositoryInterface
	userRepo repositories.UserRepositoryInterface
}

type CartQuote struct {
	ItemsTotal   float64 `json:"items_total"`
	ShippingCost float64 `json:"shipping_cost"`
	Discount     float64 `json:"discount"`
	GrandTotal   float64 `json:"grand_total"`
}

func NewShippingMethodUsecase(repo repositories.ShippingMethodRepositoryInterface, userRepo repositories.UserRepositoryInterface) ShippingMethodUsecaseInterface {
	return &shippingMethodUsecase{repo, userRepo}
}

func (u *shippingMethodUsecase) GetAll() ([]models.ShippingMethod, error) {
	return u.repo.FindAllActive()
}

func (u *shippingMethodUsecase) GetByID(id int) (*models.ShippingMethod, error) {
	return u.repo.FindByID(id)
}

func (u *shippingMethodUsecase) GetCartQuote(userID int, shippingMethodID int) (*CartQuote, error) {
	cart, err := u.userRepo.GetActiveCartByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil || len(cart.CartItems) == 0 {
		return &CartQuote{
			ItemsTotal:   0,
			ShippingCost: 0,
			Discount:     0,
			GrandTotal:   0,
		}, nil
	}

	var itemsTotal float64
	for _, item := range cart.CartItems {
		if item.Product != nil {
			itemsTotal += item.Product.Price * float64(item.Quantity)
		}
	}

	method, err := u.repo.FindByID(shippingMethodID)
	if err != nil {
		return nil, err
	}

	shippingCost := method.Cost
	discount := 0.0
	if cart.Coupon != nil {
		var percent float64
		if cart.Coupon.DiscountPercent != nil {
			percent = *cart.Coupon.DiscountPercent
		}
		var maxDisc float64
		if cart.Coupon.MaxDiscount != nil {
			maxDisc = *cart.Coupon.MaxDiscount
		}

		discount = itemsTotal * (percent / 100.0)
		if maxDisc > 0 && discount > maxDisc {
			discount = maxDisc
		}
		if discount < 0 {
			discount = 0
		}
	}
	grandTotal := itemsTotal + shippingCost - discount

	return &CartQuote{
		ItemsTotal:   itemsTotal,
		ShippingCost: shippingCost,
		Discount:     discount,
		GrandTotal:   grandTotal,
	}, nil
}
