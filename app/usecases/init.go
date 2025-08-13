package usecases

import (
	"monitoring-service/app/repositories"
	"monitoring-service/pkg/config"
)

type Main struct {
	User UserUsecaseInterface
	Auth AuthUsecaseInterface
	Shipment ShipmentUsecaseInterface
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
		User: (*userUsecase)(ucs),
		Auth: (*authUsecase)(ucs),
		Shipment: (*shipmentUsecase)(ucs),
	}

	return m
}
