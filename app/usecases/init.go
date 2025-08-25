package usecases

import (
	"monitoring-service/app/repositories"
	"monitoring-service/pkg/config"
)

type Main struct {
	User  UserUsecaseInterface
	Auth  AuthUsecaseInterface
	Cart  CartUsecaseInterface
	Order OrderUsecaseInterface
}

type usecase struct {
	Options Options
}

type Options struct {
	Repository *repositories.Main
	Config     *config.Config
}

func Init(opts Options) *Main {
	ucs := &usecase{opts}

	m := &Main{
		User:  (*userUsecase)(ucs),
		Auth:  (*authUsecase)(ucs),
		Cart:  &cartUsecase{usecase: ucs},
		Order: &orderUsecase{*ucs},
	}

	return m
}
