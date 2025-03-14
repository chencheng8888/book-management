package repo

import (
	"book-management/internal/repository/do"
	"book-management/internal/service"
)

func toServiceBook(id uint64, bookInfo do.BookInfo, bookStock do.BookStock) service.Book {
	return service.Book{
		BookID: id,
		Info: service.BookInfo{
			Name:      bookInfo.Name,
			Author:    bookInfo.Author,
			Publisher: bookInfo.Publisher,
			Category:  bookInfo.Category,
		},
		Stock: service.BookStock{
			Stock:     bookStock.Stock,
			Status:    getStockStatus(bookStock),
			Where:     bookStock.Where,
			AddedTime: bookStock.UpdatedAt,
		},
	}
}

func batchToServiceBook(bookInfos []do.BookInfo, bookStocks []do.BookStock) []service.Book {
	infoMap := make(map[uint64]do.BookInfo, len(bookInfos))
	stockMap := make(map[uint64]do.BookStock, len(bookStocks))

	for _, info := range bookInfos {
		infoMap[info.ID] = info
	}
	for _, stock := range bookStocks {
		stockMap[stock.BookID] = stock
	}

	result := make([]service.Book, 0, len(bookInfos))

	for id, _ := range infoMap {
		result = append(result, toServiceBook(id, infoMap[id], stockMap[id]))
	}
	return result
}

func batchToServiceBookRecord(borrows []do.BookBorrow, user map[string]string) []service.BookBorrowRecord {
	var records = make([]service.BookBorrowRecord, 0, len(borrows))
	for _, borrow := range borrows {
		records = append(records, service.BookBorrowRecord{
			BookID:       borrow.BookID,
			BorrowerID:   borrow.BorrowerID,
			Borrower:     user[borrow.BorrowerID],
			CopyID:       borrow.CopyID,
			BorrowTime:   borrow.CreatedTime,
			ExpectedTime: borrow.ExpectedReturnTime,
			ReturnStatus: borrow.Status,
		})
	}
	return records
}
