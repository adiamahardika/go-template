package usecases

import (
	"monitoring-service/app/repositories"
	"monitoring-service/pkg/config"

	"gorm.io/gorm"
)

type Main struct {
	User    UserUsecaseInterface
	Auth    AuthUsecaseInterface
	Product productUsecaseInterface
	Order   OrderUsecaseInterface
}

type usecase struct {
	Options Options
}

type Options struct {
	Repository *repositories.Main
	Config     *config.Config
	DB         *gorm.DB
}

func Init(opts Options) *Main {
	ucs := &usecase{opts}

	m := &Main{
		User:    (*userUsecase)(ucs),
		Auth:    (*authUsecase)(ucs),
		Product: (*productUsecase)(ucs),
		Order:   (*orderUsecase)(ucs),
	}

	return m
}
