package repo

import (
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/locker"
	"book-management/internal/pkg/tool"
	"book-management/internal/repository/do"
	"book-management/internal/service"
	"book-management/pkg/logger"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type BookBorrowDao interface {
	// 新增借阅记录
	AddBookBorrowRecord(ctx context.Context, bookID uint64, borrowerID uint64, expectedReturnTime time.Time, copyID *uint64) error
	// 获取书本借阅记录总数
	GetBookRecordTotalNum(ctx context.Context, opt ...func(db *gorm.DB)) (int, error)
	// 模糊查询借阅记录
	FuzzyQueryBookBorrowRecord(ctx context.Context, pageSize int, page int, opts ...func(db *gorm.DB)) ([]do.BookBorrow, error)
	//修改借阅状态
	UpdateBorrowStatus(ctx context.Context, bookID, copyID uint64, status string) error
	//查询书籍借阅记录的统计数据
	GetBookBorrowStatistics(ctx context.Context, startTime, endTime time.Time) (do.BorrowStatistics, error)
}

type BookBorrowCache interface {
	GetBookBorrowStatistics(ctx context.Context, pattern string) (do.BorrowStatistics, error)
	SaveBookBorrowStatistics(ctx context.Context, pattern string, num do.BorrowStatistics) error
}
type GetUserNamer interface {
	GetUserName(ctx context.Context, id ...uint64) (map[uint64]string, error)
}
type BookBorrowRepo struct {
	bookBorrowDao   BookBorrowDao
	bookBorrowCache BookBorrowCache
	userDao         GetUserNamer
	locker          *locker.Locker
}

func NewBookBorrowRepo(bookBorrowDao BookBorrowDao, bookBorrowCache BookBorrowCache, userDao GetUserNamer) *BookBorrowRepo {
	return &BookBorrowRepo{
		bookBorrowDao:   bookBorrowDao,
		bookBorrowCache: bookBorrowCache,
		userDao:         userDao,
		locker:          locker.NewLocker(),
	}
}

func (b *BookBorrowRepo) GetBookStatisticsBorrow(ctx context.Context, pattern string, startTime, endTime time.Time) (map[string]int, error) {
	//尝试从缓存中获取
	statistics, err := b.bookBorrowCache.GetBookBorrowStatistics(ctx, pattern)
	if err == nil {
		return statistics.ToMap(), nil
	}

	//这个操作，应该是个耗时操作
	if b.locker.IsLock() {
		logger.LogPrinter.Warnf("when book_repo is getting book borrow statistic, it has encountered a lock")
		return nil, errors.New("locker: book_repo locker is lock")
	}

	//加锁
	b.locker.Lock()
	defer b.locker.Unlock()

	statistics, err = b.bookBorrowDao.GetBookBorrowStatistics(ctx, startTime, endTime)
	if err == nil {
		err = b.bookBorrowCache.SaveBookBorrowStatistics(context.Background(), pattern, statistics)
		if err != nil {
			logger.LogPrinter.Warnf("cache: save book borrow statistics[pattern:%v statistics:%v] failed: %v", pattern, statistics, err)
		}
		return statistics.ToMap(), nil
	}
	//获取不到直接返回
	return nil, errors.New("get book borrow statistic failed")
}

func (b *BookBorrowRepo) UpdateBorrowStatus(ctx context.Context, bookID uint64, copyID uint64, status string) error {
	return b.bookBorrowDao.UpdateBorrowStatus(ctx, bookID, copyID, status)
}

func (b *BookBorrowRepo) QueryBookRecord(ctx context.Context, pageSize int, currentPage int, totalPage *int, opts ...func(db *gorm.DB)) ([]service.BookBorrowRecord, error) {

	if totalPage == nil || pageSize <= 0 || currentPage <= 0 {
		return nil, errcode.PageError
	}

	num, err := b.bookBorrowDao.GetBookRecordTotalNum(ctx, opts...)
	if err != nil {
		return nil, err
	}

	*totalPage = num
	if maxPage := tool.GetPage(num, pageSize); currentPage > maxPage {
		return nil, errcode.PageError
	}

	borrow, err := b.bookBorrowDao.FuzzyQueryBookBorrowRecord(ctx, pageSize, currentPage, opts...)
	if err != nil {
		return nil, err
	}

	var userID = make([]uint64, 0, len(borrow))
	for _, v := range borrow {
		userID = append(userID, v.BorrowerID)
	}

	mp, err := b.userDao.GetUserName(ctx, userID...)
	if err != nil {
		return nil, err
	}
	return batchToServiceBookRecord(borrow, mp), nil
}

func (b *BookBorrowRepo) AddBookBorrowRecord(ctx context.Context, bookID uint64, borrowerID uint64, expectedReturnTime time.Time, copyID *uint64) error {
	return b.bookBorrowDao.AddBookBorrowRecord(ctx, bookID, borrowerID, expectedReturnTime, copyID)
}
