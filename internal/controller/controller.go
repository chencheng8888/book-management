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
	NewActivityCtrl,
)

func NewCtrl(pingCtrl *PingController,
	authCtrl *AuthCtrl,
	bsCtrl *BookStockCtrl,
	bbCtrl *BookBorrowCtrl,
	userCtrl *UserCtrl,
	donateCtrl *BookDonateCtrl,
	activityCtrl *ActivityCtrl) []route.WebHandler {
	var webhandlers = []route.WebHandler{
		pingCtrl,
		authCtrl,
		bsCtrl,
		bbCtrl,
		userCtrl,
		donateCtrl,
		activityCtrl,
	}
	return webhandlers
}
