package service

import (
	"book-management/internal/controller"
	"book-management/internal/pkg/common"
	"book-management/internal/pkg/errcode"
	"book-management/pkg/logger"
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type BookRepo interface {
	CheckBookInfoIfExist(ctx context.Context, name, author, publisher, category string) (uint64, bool)

	AddBookStock(ctx context.Context, id uint64, num uint, where *string) error
	RegisterBookAndAddBookStock(ctx context.Context, id uint64, book BookInfo, num uint, where string) error

	SearchBookByID(ctx context.Context, id uint64) (Book, error)
	FuzzyQueryBook(ctx context.Context, pageSize int, currentPage int, totalPage *int, opts ...func(db *gorm.DB)) ([]Book, error)
}

type IDer interface {
	GenerateBookID(ctx context.Context) (uint64, error)
}

type BookSvc struct {
	bookRepo BookRepo
	ider     IDer
}

func (b *BookSvc) AddStock(ctx context.Context, req controller.AddStockReq) error {
	bookID, ok := b.bookRepo.CheckBookInfoIfExist(ctx, req.Name, req.Author, req.Publisher, req.Category)
	if ok {
		err := b.bookRepo.AddBookStock(ctx, bookID, req.QuantityAdded, req.Where)
		if err != nil {
			logger.LogPrinter.Errorf("add stock[id:%v addedNum:%v where: %v] failed: %v", bookID, req.QuantityAdded, req.Where, err)
			return errcode.AddBookStockError
		}
		return nil
	}

	//第一次加入库存，需给予where
	if req.Where == nil {
		return errcode.ParamError
	}

	bookID, err := b.ider.GenerateBookID(ctx)
	if err != nil {
		logger.LogPrinter.Errorf("generate book id failed: %v", err)
		return errcode.GenerateBookIDError
	}

	bookInfo := BookInfo{
		Name:      req.Name,
		Author:    req.Author,
		Publisher: req.Publisher,
		Category:  req.Category,
	}
	err = b.bookRepo.RegisterBookAndAddBookStock(ctx, bookID, bookInfo, req.QuantityAdded, *req.Where)

	if err != nil {
		logger.LogPrinter.Errorf("add stock[info:%v addedNum:%v where: %v] failed: %v", bookInfo, req.QuantityAdded, req.Where, err)
		return errcode.AddBookStockError
	}
	return nil
}

func (b *BookSvc) SearchBookStockByID(ctx context.Context, req controller.SearchStockByBookIDReq) (controller.Book, error) {
	book, err := b.bookRepo.SearchBookByID(ctx, req.BookID)
	if err != nil {
		return controller.Book{}, errcode.SearchBookError
	}
	return controller.Book{
		BookID:      book.BookID,
		Name:        book.Info.Name,
		Author:      book.Info.Author,
		Publisher:   book.Info.Publisher,
		Category:    book.Info.Category,
		Stock:       book.Stock.Stock,
		StockStatus: book.Stock.Status,
		StockWhere:  book.Stock.Where,
		CreatedAt:   convertTimeFormat(book.Stock.AddedTime),
	}, nil
}

func (b *BookSvc) FuzzyQueryBookStock(ctx context.Context, req controller.FuzzyQueryBookStockReq, totalPage *int) ([]controller.Book, error) {
	var Opts []func(db *gorm.DB)

	if req.Name != nil {
		Opts = append(Opts, func(db *gorm.DB) {
			db.Where(fmt.Sprintf("%s.name = ?", common.BookTableName), req.Name)
		})
	}
	if req.Author != nil {
		Opts = append(Opts, func(db *gorm.DB) {
			db.Where(fmt.Sprintf("%s.author = ?", common.BookTableName), req.Author)
		})
	}
	if req.Category != nil {
		Opts = append(Opts, func(db *gorm.DB) {
			db.Where(fmt.Sprintf("%s.category = ?", common.BookTableName), *req.Category)
		})
	}

	if req.AddStockTime != nil {
		t, err := convertStringToTime(*req.AddStockTime)
		if err != nil {
			return nil, errcode.ParamError
		}
		Opts = append(Opts, func(db *gorm.DB) {
			db.Where(fmt.Sprintf("%s.updated_at >= ? AND %s.updated_at < ?", common.BookStockTableName, common.BookStockTableName), t, t.Add(24*time.Hour))
		})
	}
	if req.AddStockWhere != nil {
		Opts = append(Opts, func(db *gorm.DB) {
			db.Where(fmt.Sprintf("%s.where = ?", common.BookStockTableName), *req.AddStockWhere)
		})
	}

	books, err := b.bookRepo.FuzzyQueryBook(ctx, req.PageSize, req.Page, totalPage, Opts...)
	if err != nil {
		return nil, errcode.SearchBookError
	}

	if req.Page > *totalPage {
		return nil, errcode.PageError
	}

	return batchToControllerBook(books), nil
}

func convertTimeFormat(t time.Time) string {
	return t.Format("2006-01-02")
}
func convertStringToTime(t string) (time.Time, error) {
	layout := "2006-01-02"

	// 加载上海时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Time{}, err
	}

	// 直接解析为上海时区的时间
	timeObj, err := time.ParseInLocation(layout, t, loc)
	if err != nil {
		return time.Time{}, err
	}

	return timeObj, nil
}

func toControllerBook(book Book) controller.Book {
	return controller.Book{
		BookID:      book.BookID,
		Name:        book.Info.Name,
		Author:      book.Info.Author,
		Publisher:   book.Info.Publisher,
		Category:    book.Info.Category,
		Stock:       book.Stock.Stock,
		StockStatus: book.Stock.Status,
		StockWhere:  book.Stock.Where,
		CreatedAt:   convertTimeFormat(book.Stock.AddedTime),
	}
}
func batchToControllerBook(books []Book) []controller.Book {
	var result []controller.Book
	for _, book := range books {
		result = append(result, toControllerBook(book))
	}
	return result
}
