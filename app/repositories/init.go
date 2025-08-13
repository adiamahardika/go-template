package repositories

import (
	"monitoring-service/pkg/config"
	"gorm.io/gorm"
)

type Main struct {
	User      UserRepositoryInterface
	UserRole  UserRolesRepositoryInterface
	Product   ProductRepository
	Category  CategoryRepository
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

       var productRepo ProductRepository
       var categoryRepo CategoryRepository
       if opts.Postgres != nil {
	       // Directly use *gorm.DB
	       productRepo = NewProductRepository(opts.Postgres)
	       categoryRepo = NewCategoryRepository(opts.Postgres)
       }
       m := &Main{
	       User:     (*userRepository)(repo),
	       UserRole: (*userRolesRepository)(repo),
	       Product:  productRepo,
	       Category: categoryRepo,
       }
       return m
}
