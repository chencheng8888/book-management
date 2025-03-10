package repo

import (
	"book-management/internal/pkg/common"
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/locker"
	"book-management/internal/repository/do"
	"book-management/internal/service"
	"book-management/pkg/logger"
	"context"
	"errors"
	"math"
	"time"

	"gorm.io/gorm"
)

type BookDao interface {
	// 新增书籍库存
	AddBookStock(ctx context.Context, id uint64, num uint, where *string) error
	// 注册并新增书籍库存
	RegisterAndAddBookStock(ctx context.Context, bookInfo do.BookInfo, addedNum uint, where string) error
	// 新增借阅记录
	AddBookBorrowRecord(ctx context.Context, bookID uint64, borrowerID string, expectedReturnTime time.Time, copyID *uint64) error
	// 根据ID获取书籍信息
	GetBookInfoByID(ctx context.Context, id ...uint64) ([]do.BookInfo, error)
	// 根据ID获取书籍库存
	GetBookStockByID(ctx context.Context, id ...uint64) ([]do.BookStock, error)
	// 检查书本是否存在
	CheckBookIfExist(ctx context.Context, name, author, publisher, category string) (uint64, bool)
	//模糊查询书籍ID
	FuzzyQueryBookID(ctx context.Context, pageSize int, page int, opts ...func(db *gorm.DB)) ([]uint64, error)
	// 获取书籍总数
	GetBookTotalNum(ctx context.Context, opts ...func(db *gorm.DB)) (int, error)
	// 获取书本借阅记录总数
	GetBookRecordTotalNum(ctx context.Context, opt ...func(db *gorm.DB)) (int, error)
	// 模糊查询借阅记录
	FuzzyQueryBookBorrowRecord(ctx context.Context, pageSize int, page int, opts ...func(db *gorm.DB)) ([]do.BookBorrow, error)
	//修改借阅状态
	UpdateBorrowStatus(ctx context.Context, bookID, copyID uint64, status string) error
	//查询书籍借阅记录的统计数据
	GetBookBorrowStatistics(ctx context.Context, startTime, endTime time.Time) (do.BorrowStatistics, error)
}

type UserDao interface {
	GetUserName(ctx context.Context, id ...string) (map[string]string, error)
}

type BookCache interface {
	DeleteBookStock(ctx context.Context, id uint64) error
	DeleteBookInfo(ctx context.Context, id uint64) error

	GetBookInfoByID(ctx context.Context, ids ...uint64) ([]do.BookInfo, []uint64)
	GetBookStockByID(ctx context.Context, ids ...uint64) ([]do.BookStock, []uint64)
	GetBookBorrowStatistics(ctx context.Context, pattern string) (do.BorrowStatistics, error)

	SaveBookInfo(ctx context.Context, infos ...do.BookInfo) error
	SaveBookStock(ctx context.Context, stocks ...do.BookStock) error

	SaveBookBorrowStatistics(ctx context.Context, pattern string, num do.BorrowStatistics) error
}

type BookRepo struct {
	bookDao   BookDao
	userDao   UserDao
	bookCache BookCache
	locker    *locker.Locker
}

func (b *BookRepo) GetBookStatisticsBorrow(ctx context.Context, pattern string, startTime, endTime time.Time) (map[string]int, error) {
	//尝试从缓存中获取
	statistics, err := b.bookCache.GetBookBorrowStatistics(ctx, pattern)
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

	statistics, err = b.bookDao.GetBookBorrowStatistics(ctx, startTime, endTime)
	if err == nil {
		err = b.bookCache.SaveBookBorrowStatistics(context.Background(), pattern, statistics)
		if err != nil {
			logger.LogPrinter.Warnf("cache: save book borrow statistics[pattern:%v statistics:%v] failed: %v", pattern, statistics, err)
		}
		return statistics.ToMap(), nil
	}
	//获取不到直接返回
	return nil, errors.New("get book borrow statistic failed")
}

func (b *BookRepo) UpdateBorrowStatus(ctx context.Context, bookID uint64, copyID uint64, status string) error {
	return b.bookDao.UpdateBorrowStatus(ctx, bookID, copyID, status)
}

func (b *BookRepo) QueryBookRecord(ctx context.Context, pageSize int, currentPage int, totalPage *int, opts ...func(db *gorm.DB)) ([]service.BookBorrowRecord, error) {
	if pageSize <= 0 || currentPage <= 0 {
		return nil, errcode.PageError
	}

	num, err := b.bookDao.GetBookRecordTotalNum(ctx, opts...)
	if err != nil {
		return nil, err
	}

	pageNum := int(math.Ceil(float64(num) / float64(pageSize)))

	totalPage = &pageNum

	if currentPage > pageNum {
		return nil, errcode.PageError
	}

	borrow, err := b.bookDao.FuzzyQueryBookBorrowRecord(ctx, pageSize, currentPage, opts...)
	if err != nil {
		return nil, err
	}

	var userID = make([]string, 0, len(borrow))
	for _, v := range borrow {
		userID = append(userID, v.BorrowerID)
	}

	mp, err := b.userDao.GetUserName(ctx, userID...)
	if err != nil {
		return nil, err
	}
	return batchToServiceBookRecord(borrow, mp), nil
}

func (b *BookRepo) AddBookBorrowRecord(ctx context.Context, bookID uint64, borrowerID string, expectedReturnTime time.Time, copyID *uint64) error {
	return b.bookDao.AddBookBorrowRecord(ctx, bookID, borrowerID, expectedReturnTime, copyID)
}

func (b *BookRepo) CheckBookInfoIfExist(ctx context.Context, name, author, publisher, category string) (uint64, bool) {
	return b.bookDao.CheckBookIfExist(ctx, name, author, publisher, category)
}

func (b *BookRepo) AddBookStock(ctx context.Context, id uint64, num uint, where *string) error {
	err := b.bookCache.DeleteBookStock(ctx, id)
	if err != nil {
		return err
	}

	//异步删除缓存
	defer func() {
		go time.AfterFunc(time.Second, func() {
			_ = b.bookCache.DeleteBookStock(ctx, id)
		})
	}()
	return b.bookDao.AddBookStock(ctx, id, num, where)
}

func (b *BookRepo) RegisterBookAndAddBookStock(ctx context.Context, id uint64, book service.BookInfo, num uint, where string) error {
	clean := func(ctx context.Context) error {
		if err := b.bookCache.DeleteBookInfo(ctx, id); err != nil {
			return err
		}
		if err := b.bookCache.DeleteBookStock(ctx, id); err != nil {
			return err
		}

		return nil
	}

	if err := clean(ctx); err != nil {
		return err
	}
	//异步删除缓存
	defer func() {
		go time.AfterFunc(time.Second, func() {
			_ = clean(context.Background())
		})
	}()

	currentTime := time.Now()

	return b.bookDao.RegisterAndAddBookStock(ctx, do.BookInfo{
		ID:        id,
		Name:      book.Name,
		Author:    book.Author,
		Publisher: book.Publisher,
		Category:  book.Category,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}, num, where)
}

func (b *BookRepo) SearchBookByID(ctx context.Context, id uint64) (service.Book, error) {
	bookInfos, bookStocks, err := b.getBookInID(ctx, id)
	if err != nil {
		return service.Book{}, err
	}

	return toServiceBook(id, bookInfos[0], bookStocks[0]), nil
}

func (b *BookRepo) FuzzyQueryBook(ctx context.Context, pageSize int, currentPage int, totalPage *int, opts ...func(db *gorm.DB)) ([]service.Book, error) {
	num, err := b.bookDao.GetBookTotalNum(ctx)
	if err != nil {
		return nil, err
	}

	pageNum := int(math.Ceil(float64(num) / float64(pageSize)))

	totalPage = &pageNum

	if currentPage > pageNum {
		return nil, nil
	}

	ids, err := b.bookDao.FuzzyQueryBookID(ctx, pageSize, currentPage, opts...)
	if err != nil {
		return nil, err
	}

	bookInfos, bookStocks, err := b.getBookInID(ctx, ids...)
	if err != nil {
		return nil, err
	}
	return batchToServiceBook(bookInfos, bookStocks), nil
}

func (b *BookRepo) getBookInID(ctx context.Context, ids ...uint64) ([]do.BookInfo, []do.BookStock, error) {
	bookInfos, leftID := b.bookCache.GetBookInfoByID(ctx, ids...)

	if len(leftID) > 0 {
		tmpBookInfos, err := b.bookDao.GetBookInfoByID(ctx, leftID...)
		if err != nil || len(tmpBookInfos) != len(leftID) {
			return nil, nil, err
		}
		//异步保存书籍信息
		defer func() {
			go b.bookCache.SaveBookInfo(context.Background(), tmpBookInfos...)
		}()
		bookInfos = append(bookInfos, tmpBookInfos...)
	}

	bookStocks, leftID := b.bookCache.GetBookStockByID(ctx, ids...)

	if len(leftID) > 0 {
		tmpBookStocks, err := b.bookDao.GetBookStockByID(ctx, leftID...)
		if err != nil || len(tmpBookStocks) != len(leftID) {
			return nil, nil, err
		}
		//异步保存书籍库存
		defer func() {
			go b.bookCache.SaveBookStock(context.Background(), tmpBookStocks...)
		}()
		bookStocks = append(bookStocks, tmpBookStocks...)
	}
	return bookInfos, bookStocks, nil
}

func getStockStatus(stock do.BookStock) string {
	if stock.IsAdequate() {
		return common.Adequate
	}
	if stock.IsEarlyWarning() {
		return common.EarlyWarning
	}
	return common.Shortage
}
