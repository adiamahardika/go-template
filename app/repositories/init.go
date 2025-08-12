package repositories

import (
	"monitoring-service/pkg/config"

	"gorm.io/gorm"
)

type Main struct {
	User     UserRepositoryInterface
	UserRole UserRolesRepositoryInterface
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
		User:     (*userRepository)(repo),
		UserRole: (*userRolesRepository)(repo),
	}

	return m
}
