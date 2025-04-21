package repositories

import (
	"monitoring-service/pkg/config"

	"gorm.io/gorm"
)

type Main struct {
	Todo     TodoRepositoryInterface
	Priority PriorityRepositoryInterface
}

type repository struct {
	Options Options
}

type Options struct {
	Postgres *gorm.DB
	Config   *config.Config
}

func Init(opts Options) *Main {
	repo := &repository{opts}

	m := &Main{
		Todo:     (*todoRepository)(repo),
		Priority: (*priorityRepository)(repo),
	}

	return m
}
