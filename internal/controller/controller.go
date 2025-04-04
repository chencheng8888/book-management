package controller

import (
	"book-management/internal/route"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewPingController,
	NewAuthCtrl,
	NewBookStockCtrl,
	NewBookBorrowCtrl,
	NewUserCtrl,
	NewCtrl,
)

func NewCtrl(pingCtrl *PingController, authCtrl *AuthCtrl, bsCtrl *BookStockCtrl, bbCtrl *BookBorrowCtrl, userCtrl *UserCtrl) []route.WebHandler {
	var webhandlers = []route.WebHandler{
		pingCtrl,
		authCtrl,
		bsCtrl,
		bbCtrl,
		userCtrl,
	}
	return webhandlers
}
