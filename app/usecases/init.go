package usecases

import (
	"monitoring-service/app/repositories"
	"monitoring-service/pkg/config"
)

type Main struct {
	Todo TodoUsecaseInterface
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
		Todo: (*todoUsecase)(ucs),
	}

	return m
}
