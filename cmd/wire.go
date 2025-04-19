//go:build wireinject
// +build wireinject

package main

import (
	"book-management/configs"
	"book-management/internal/controller"
	"book-management/internal/ioc"
	"book-management/internal/middleware"
	"book-management/internal/repository"
	"book-management/internal/repository/cache"
	"book-management/internal/repository/dao"
	"book-management/internal/repository/repo"
	"book-management/internal/service"

	// "book-management/internal/middleware"
	"book-management/internal/route"

	"github.com/google/wire"
)

func InitializeApp(config configs.AppConfig) (*App, error) {
	wire.Build(
		ioc.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		controller.ProviderSet,
		middleware.ProviderSet,
		route.ProviderSet,
		newApp,
		wire.Bind(new(middleware.VerifyToken), new(*controller.AuthCtrl)),
		wire.Bind(new(controller.BookBorrowSvc), new(*service.BookSvc)),
		wire.Bind(new(controller.BookStockSvc), new(*service.BookSvc)),
		wire.Bind(new(controller.BookDonateSvc), new(*service.BookSvc)),
		wire.Bind(new(controller.UserSvc), new(*service.UserSvc)),
		wire.Bind(new(controller.ActivitySvc), new(*service.ActivitySvc)),
		wire.Bind(new(service.BookStockRepo), new(*repo.BookStockRepo)),
		wire.Bind(new(service.BookBorrowRepo), new(*repo.BookBorrowRepo)),
		wire.Bind(new(service.BookDonateRepo), new(*repo.BookDonateRepo)),
		wire.Bind(new(service.UserRepo), new(*repo.UserRepo)),
		wire.Bind(new(service.ActivityRepo), new(*repo.ActivityRepo)),
		wire.Bind(new(repo.BookStockDao), new(*dao.BookDao)),
		wire.Bind(new(repo.BookInfoDao), new(*dao.BookDao)),
		wire.Bind(new(repo.BookBorrowDao), new(*dao.BookDao)),
		wire.Bind(new(repo.BookDonateDao), new(*dao.BookDao)),
		wire.Bind(new(repo.UserDao), new(*dao.UserDao)),
		wire.Bind(new(repo.GetUserNamer), new(*dao.UserDao)),
		wire.Bind(new(repo.GetUserPhone), new(*dao.UserDao)),
		wire.Bind(new(repo.ActivityDao), new(*dao.ActivityDao)),
		wire.Bind(new(repo.BookBorrowCache), new(*cache.BookCache)),
		wire.Bind(new(repo.BookCache), new(*cache.BookCache)),
	)
	return &App{}, nil
}
