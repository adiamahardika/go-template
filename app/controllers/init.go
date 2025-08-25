package controllers

import (
	"monitoring-service/app/usecases"
	"monitoring-service/pkg/config"
)

type Main struct {
	User UserControllerInterface
	Auth AuthControllerInterface
	Cart CartControllerInterface // Tambahkan field ini
	Order OrderControllerInterface // Tambahkan field OrderControllerInterface
}

type controller struct {
	Options Options
}

type Options struct {
	Config   *config.Config
	UseCases *usecases.Main
}

func Init(opts Options) *Main {
	ctrl := &controller{opts}

	m := &Main{
		User: (*userController)(ctrl),
		Auth: (*authController)(ctrl),
		Cart: (*cartController)(ctrl), // Inisialisasi CartController
		Order: (*orderController)(ctrl), // Inisialisasi OrderController
	}

	return m
}
