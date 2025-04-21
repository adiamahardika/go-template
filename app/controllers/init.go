package controllers

import (
	"monitoring-service/app/usecases"
	"monitoring-service/pkg/config"
)

type Main struct {
	Todo     TodoControllersInterface
	Priority PriorityControllersInterface
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
		Todo:     (*todoControllers)(ctrl),
		Priority: (*priorityControllers)(ctrl),
	}

	return m
}
