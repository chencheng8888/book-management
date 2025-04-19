package service

import (
	"book-management/internal/controller"
	"book-management/internal/pkg/common"
	"book-management/internal/pkg/errcode"
	"book-management/pkg/logger"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"gorm.io/gorm"
)

type BookStockRepo interface {
	CheckBookInfoIfExist(ctx context.Context, name, author, publisher, category string) (uint64, bool)
	AddBookStock(ctx context.Context, bookID uint64, userID *uint64, num uint) error
	RegisterBookAndAddBookStock(ctx context.Context, bookID uint64, userID *uint64, book BookInfo, num uint) error
	FuzzyQueryBook(ctx context.Context, pageSize int, currentPage int, total *int, opts ...func(db *gorm.DB)) ([]Book, error)
}

type BookBorrowRepo interface {
	QueryBookRecord(ctx context.Context, pageSize int, currentPage int, total *int, opts ...func(db *gorm.DB)) ([]BookBorrowRecord, error)
	AddBookBorrowRecord(ctx context.Context, bookID uint64, borrowerID uint64, expectedReturnTime time.Time, copyID uint64) error
	UpdateBorrowStatus(ctx context.Context, bookID uint64, copyID uint64, status string) error
	GetBookStatisticsBorrow(ctx context.Context, pattern string, startTime time.Time, endTime time.Time) (map[string]int, error)
	GetAvailableCopyBook(ctx context.Context, bookID uint64, page, pageSize int) ([]uint64, error)
}

type BookDonateRepo interface {
	ListBookDonateRecordsReq(ctx context.Context, pageSize, currentPage int, total *int) ([]BookDonateRecord, error)
	GetBookDonateRanking(ctx context.Context, top int) ([]Rank, error)
}

type BookSvc struct {
	bookBorrowRepo BookBorrowRepo
	bookStockRepo  BookStockRepo
	bookDonateRepo BookDonateRepo
	ider           *MyIDer
}

func NewBookSvc(bookBorrowRepo BookBorrowRepo, bookStockRepo BookStockRepo, bookDonateRepo BookDonateRepo) *BookSvc {
	return &BookSvc{
		bookBorrowRepo: bookBorrowRepo,
		bookStockRepo:  bookStockRepo,
		bookDonateRepo: bookDonateRepo,
		ider:           NewMyIDer(),
	}
}

func (b *BookSvc) ListDonateRecordsReq(ctx context.Context, req controller.ListDonateRecordsReq, total *int) ([]controller.DonateRecords, error) {
	records, err := b.bookDonateRepo.ListBookDonateRecordsReq(ctx, req.PageSize, req.Page, total)
	if errors.Is(err, errcode.PageError) {
		return nil, errcode.PageError
	}
	if err != nil {
		return nil, errcode.SearchDataError
	}
	return batchToControllerBookDonateRecord(records), nil
}

func (b *BookSvc) GetDonationRanking(ctx context.Context, req controller.GetDonationRankingReq) ([]controller.Rank, error) {
	rankings, err := b.bookDonateRepo.GetBookDonateRanking(ctx, req.Top)
	if err != nil {
		return nil, errcode.SearchDataError
	}
	return batchToControllerRank(rankings), nil
}
func (b *BookSvc) GetStatisticBorrowRecords(ctx context.Context, req controller.QueryStatisticsBorrowRecordsReq) (map[string]int, error) {
	switch req.Pattern {
	case WeekPattern, MonthPattern, YearPattern:
	default:
		return nil, errcode.ParamError
	}
	startTime, endTime := getStartAndEndTime(req.Pattern)
	mp, err := b.bookBorrowRepo.GetBookStatisticsBorrow(ctx, req.Pattern, startTime, endTime)
	if err != nil {
		return nil, errcode.SearchDataError
	}
	return mp, nil
}

func (b *BookSvc) UpdateBorrowStatus(ctx context.Context, req controller.UpdateBorrowStatusReq) error {
	return b.bookBorrowRepo.UpdateBorrowStatus(ctx, req.BookID, req.CopyID, req.Status)
}

func (b *BookSvc) AddBorrowBookRecord(ctx context.Context, req controller.BorrowBookReq) (uint64, uint64, error) {
	expectedTime, err := convertStringToTime(req.ExpectedReturnTime)
	if err != nil {
		logger.LogPrinter.Errorf("parse time str[%v] failed:%v", req.ExpectedReturnTime, err)
		return 0, 0, err
	}

	if expectedTime.Before(time.Now()) {
		logger.LogPrinter.Infof("expected return time[%v] is  before timw.Now", expectedTime)
		return 0, 0, errors.New("expected return time is illegal")
	}
	copyID := req.CopyID
	err = b.bookBorrowRepo.AddBookBorrowRecord(ctx, req.BookID, req.BorrowerID, expectedTime, copyID)
	if errors.Is(err, errcode.InsufficientBookStock) {
		return req.BookID, copyID, err
	}
	if err != nil {
		return req.BookID, copyID, errcode.AddBookBorrowError
	}
	return req.BookID, copyID, nil
}

func (b *BookSvc) GetAvailableCopyBook(ctx context.Context, req controller.GetAvailableCopyBookReq) ([]uint64, error) {
	return b.bookBorrowRepo.GetAvailableCopyBook(ctx, req.BookID, req.Page, req.PageSize)
}

func (b *BookSvc) QueryBookBorrowRecord(ctx context.Context, req controller.QueryBookBorrowRecordReq, totalNum *int) ([]controller.BookBorrowRecord, error) {

	var opt func(db *gorm.DB)

	if req.QueryStatus != nil {
		switch *req.QueryStatus {
		case common.BookStatusWaitingReturn, common.BookStatusOverdue, common.BookStatusReturned:
			opt = func(db *gorm.DB) {
				db.Where(fmt.Sprintf("%s.status = ?", common.BookBorrowTableName), *req.QueryStatus)
			}
		default:
			return nil, errcode.ParamError
		}
	}
	var (
		records []BookBorrowRecord
		err     error
	)
	if opt == nil {
		records, err = b.bookBorrowRepo.QueryBookRecord(ctx, req.PageSize, req.Page, totalNum)
	} else {
		records, err = b.bookBorrowRepo.QueryBookRecord(ctx, req.PageSize, req.Page, totalNum, opt)
	}
	if errors.Is(err, errcode.PageError) {
		return nil, errcode.PageError
	}
	if err != nil {
		return nil, errcode.SearchDataError
	}
	return batchToControllerBookBorrowRecord(records), nil
}

func (b *BookSvc) AddStock(ctx context.Context, req controller.AddStockReq, ID *uint64) error {

	bookID, ok := b.bookStockRepo.CheckBookInfoIfExist(ctx, req.Name, req.Author, req.Publisher, req.Category)
	if ok {
		err := b.bookStockRepo.AddBookStock(ctx, bookID, req.UserID, req.QuantityAdded)
		if err != nil {
			logger.LogPrinter.Errorf("add stock[id:%v addedNum:%v] failed: %v", bookID, req.QuantityAdded, err)
			return errcode.AddBookStockError
		}
		*ID = bookID
		return nil
	}

	bookID, err := b.ider.GenerateID(ctx)
	if err != nil {
		logger.LogPrinter.Errorf("generate book id failed: %v", err)
		return errcode.GenerateIDError
	}
	*ID = bookID
	bookInfo := BookInfo{
		Name:      req.Name,
		Author:    req.Author,
		Publisher: req.Publisher,
		Category:  req.Category,
	}
	err = b.bookStockRepo.RegisterBookAndAddBookStock(ctx, bookID, req.UserID, bookInfo, req.QuantityAdded)

	if err != nil {
		logger.LogPrinter.Errorf("add stock[info:%v addedNum:%v] failed: %v", bookInfo, req.QuantityAdded, err)
		return errcode.AddBookStockError
	}
	return nil
}

func (b *BookSvc) FuzzyQueryBookStock(ctx context.Context, req controller.FuzzyQueryBookStockReq, totalNum *int) ([]controller.Book, error) {
	var Opts []func(db *gorm.DB)

	if req.BookID != nil {
		Opts = append(Opts, func(db *gorm.DB) {
			db.Where(fmt.Sprintf("%s.id = ?", common.BookTableName), req.BookID)
		})
		//如果ID不为空，直接查询
		goto DB
	}

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

DB:
	books, err := b.bookStockRepo.FuzzyQueryBook(ctx, req.PageSize, req.Page, totalNum, Opts...)
	if errors.Is(err, errcode.PageError) {
		return nil, errcode.PageError
	}
	if err != nil {
		return nil, errcode.SearchDataError
	}
	return batchToControllerBook(books), nil
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

type MyIDer struct {
	randPool *rand.Rand
	mu       sync.Mutex
}

func NewMyIDer() *MyIDer {
	return &MyIDer{
		randPool: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (i *MyIDer) GenerateID(ctx context.Context) (uint64, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// 使用纳秒时间戳的中间部分，避免高位和低位的不稳定性
	ts := uint64(time.Now().UnixNano())

	// 取中间15位 (从第5位到第19位)
	// 这样避免最高位总是1和最低位的随机性不足
	id := (ts / 1e4) % 1e15

	// 确保在10-15位范围内
	switch {
	case id < 1e9: // 小于10位
		// 组合时间戳后5位和随机10位
		id = (ts%1e5)*1e10 + i.randPool.Uint64()%1e10
	case id > 1e15-1: // 超过15位
		id = uint64(id) % 1e15
	}

	// 最终确保在10-15位
	if id < 1e9 {
		id += 1e9
	}

	return id, nil
}
