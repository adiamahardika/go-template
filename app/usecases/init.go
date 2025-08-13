package usecases

import (
	"monitoring-service/app/repositories"
	"monitoring-service/pkg/config"
)

type Main struct {
	User    UserUsecaseInterface
	Auth    AuthUsecaseInterface
	Product productUsecaseInterface
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
		User:    (*userUsecase)(ucs),
		Auth:    (*authUsecase)(ucs),
		Product: (*productUsecase)(ucs),
	}

	return m
}
