package controller

import (
	"book-management/internal/controller/ping"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(ping.NewPingController)