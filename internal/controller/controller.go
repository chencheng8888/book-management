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
	NewBookDonateCtrl,
	NewCtrl,
)

func NewCtrl(pingCtrl *PingController, authCtrl *AuthCtrl, bsCtrl *BookStockCtrl, bbCtrl *BookBorrowCtrl, userCtrl *UserCtrl, donateCtrl *BookDonateCtrl) []route.WebHandler {
	var webhandlers = []route.WebHandler{
		pingCtrl,
		authCtrl,
		bsCtrl,
		bbCtrl,
		userCtrl,
		donateCtrl,
	}
	return webhandlers
}
