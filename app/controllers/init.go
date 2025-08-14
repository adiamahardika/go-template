package controllers

import (
	"monitoring-service/app/usecases"
	"monitoring-service/pkg/config"
)

type Main struct {
	User            UserControllerInterface
	ShippingPayment ShippingPaymentControllerInterface
}

type controller struct {
	Options Options
}

type Options struct {
	Config   *config.Config
	UseCases *usecases.Main
}

func Init(opts *Options) *Main {
	ctrl := &controller{*opts}

	m := &Main{
		User:            (*userController)(ctrl),
		ShippingPayment: newShippingPaymentController(ctrl),
	}

	return m
}
