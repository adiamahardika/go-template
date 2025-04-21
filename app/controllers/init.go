package controllers

import (
	"monitoring-service/app/usecases"
	"monitoring-service/pkg/config"
)

type Main struct {
	Todo TodoControllersInterface
	Status StatusControllersInterface
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
		Todo: (*todoControllers)(ctrl),
		Status: (*statusControllers)(ctrl),
	}

	return m
}
