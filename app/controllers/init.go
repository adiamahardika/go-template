package controllers

import (
	"monitoring-service/app/usecases"
	"monitoring-service/pkg/config"
)

type Main struct {
	Todo TodoControllersInterface
	Label LabelControllersInterface
	
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
		Label: (*labelControllers)(ctrl),
	}

	return m
}
