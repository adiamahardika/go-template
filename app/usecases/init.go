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
	Payment paymentUsecaseInterface
}

type usecase struct {
	Options Options
}

type Options struct {
	Repository *repositories.Main
	Config     *config.Config
	Postgres   *gorm.DB
}

func Init(opts Options) *Main {
	ucs := &usecase{opts}

	m := &Main{
		User:    (*userUsecase)(ucs),
		Auth:    (*authUsecase)(ucs),
		Product: (*productUsecase)(ucs),
		Payment: (*paymentUsecase)(ucs),
	}

	return m
}
