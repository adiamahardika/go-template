package repositories

import (
	"monitoring-service/pkg/config"

	"gorm.io/gorm"
)

type repository struct {
	Options Options
}

type Options struct {
	Postgres *gorm.DB
	Config   *config.Config
}

type Main struct {
	User            UserRepositoryInterface
	ShippingPayment ShippingPaymentRepositoryInterface
}

func Init(options *Options) *Main {
	repo := &repository{
		Options: *options,
	}

	return &Main{
		User:            NewUserRepository(repo),
		ShippingPayment: NewShippingPaymentRepository(repo),
	}
}
