package usecases

import (
	"monitoring-service/app/repositories"
	"monitoring-service/pkg/config"
)

type Main struct {
	User     UserUsecaseInterface
	Auth     AuthUsecaseInterface
	Product  ProductUsecase
	Category CategoryUsecase
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

       var productUsecase ProductUsecase
       var categoryUsecase CategoryUsecase
       if opts.Repository != nil {
	       if opts.Repository.Product != nil {
		       productUsecase = NewProductUsecase(opts.Repository.Product)
	       }
	       if opts.Repository.Category != nil {
		       categoryUsecase = NewCategoryUsecase(opts.Repository.Category)
	       }
       }
       m := &Main{
	       User:     (*userUsecase)(ucs),
	       Auth:     (*authUsecase)(ucs),
	       Product:  productUsecase,
	       Category: categoryUsecase,
       }
       return m
}
