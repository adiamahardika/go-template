package controllers

import (
	"monitoring-service/app/usecases"
	"monitoring-service/pkg/config"
)

type Main struct {
	User     UserControllerInterface
	Auth     AuthControllerInterface
	Product *ProductController
	Category *CategoryController
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

       // Wiring ProductController
       var productController *ProductController
       if opts.UseCases != nil && opts.UseCases.Product != nil {
	       productController = NewProductController(opts.UseCases.Product)
       }
       // Wiring CategoryController
       var categoryController *CategoryController
       if opts.UseCases != nil && opts.UseCases.Category != nil {
	       categoryController = NewCategoryController(opts.UseCases.Category)
       }
       m := &Main{
	       User:     (*userController)(ctrl),
	       Auth:     (*authController)(ctrl),
	       Product:  productController,
	       Category: categoryController,
       }
       return m
}
