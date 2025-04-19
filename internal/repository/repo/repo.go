package repo

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewBookStockRepo,
	NewBookBorrowRepo,
	NewBookDonateRepo,
	NewActivityRepo,
	NewUserRepo)
