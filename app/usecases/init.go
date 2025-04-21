package usecases

import (
	"monitoring-service/app/repositories"
	"monitoring-service/pkg/config"
)

type Main struct {
	Todo     TodoUsecaseInterface
	Category CategoryUsecaseInterface
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
		Todo:     (*todoUsecase)(ucs),
		Category: (*categoryUsecase)(ucs),
	}

	return m
}
