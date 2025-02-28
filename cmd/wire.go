//+build wireinject

package main

import (
	"book-management/internal/controller"
	// "book-management/internal/middleware"
	"book-management/internal/route"

	"github.com/google/wire"
)

func InitializeApp() *App {
	wire.Build(controller.ProviderSet,
		// middleware.ProviderSet,
		route.ProviderSet,
		newApp,
	)
	return &App{}
}