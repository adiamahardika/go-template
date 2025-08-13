package controllers

import (
	"monitoring-service/app/usecases"
	"monitoring-service/pkg/config"
)

type Main struct {
	User    UserControllerInterface
	Auth    AuthControllerInterface
	Product productControllerInterface
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
		User:    (*userController)(ctrl),
		Auth:    (*authController)(ctrl),
		Product: (*productController)(ctrl),
	}

	return m
}
