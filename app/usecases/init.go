package usecases

import (
	"monitoring-service/app/repositories"
	"monitoring-service/pkg/config"
)

type Main struct {
	User            UserUsecaseInterface
	ShippingPayment ShippingPaymentUsecaseInterface
}

type usecase struct {
	options *Options
}

type Options struct {
	Repository *repositories.Main
	Config     *config.Config
}

func Init(opts *Options) *Main {
	uc := &usecase{
		options: opts,
	}

	return &Main{
		User:            newUserUsecase(uc),
		ShippingPayment: newShippingPaymentUsecase(uc),
	}
}

// unexported constructor functions
func newUserUsecase(uc *usecase) UserUsecaseInterface {
	return &userUsecase{uc}
}

func newShippingPaymentUsecase(uc *usecase) ShippingPaymentUsecaseInterface {
	return &shippingPaymentUsecase{uc}
}
