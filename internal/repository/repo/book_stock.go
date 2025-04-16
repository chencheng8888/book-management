package repo

import (
	"book-management/internal/pkg/common"
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/tool"
	"book-management/internal/repository/do"
	"book-management/internal/service"
	"context"
	"time"

	"gorm.io/gorm"
)

type BookInfoDao interface {
	// 根据ID获取书籍信息
	GetBookInfoByID(ctx context.Context, id ...uint64) ([]do.BookInfo, error)
}

type BookStockDao interface {
	// 新增书籍库存
	AddBookStock(ctx context.Context, bookID, userID uint64, num uint, donate bool) error
	// 注册并新增书籍库存
	RegisterAndAddBookStock(ctx context.Context, bookInfo do.BookInfo, userID uint64, addedNum uint, donate bool) error
	// 根据ID获取书籍库存
	GetBookStockByID(ctx context.Context, id ...uint64) ([]do.BookStock, error)
	// 检查书本是否存在
	CheckBookIfExist(ctx context.Context, name, author, publisher, category string) (uint64, bool)
	//模糊查询书籍ID
	FuzzyQueryBookID(ctx context.Context, pageSize int, page int, opts ...func(db *gorm.DB)) ([]uint64, error)
	// 获取书籍总数
	GetBookTotalNum(ctx context.Context, opts ...func(db *gorm.DB)) (int, error)
}

type BookCache interface {
	DeleteBookStock(ctx context.Context, id uint64) error
	DeleteBookInfo(ctx context.Context, id uint64) error

	GetBookInfoByID(ctx context.Context, ids ...uint64) ([]do.BookInfo, []uint64)
	GetBookStockByID(ctx context.Context, ids ...uint64) ([]do.BookStock, []uint64)

	SaveBookInfo(ctx context.Context, infos ...do.BookInfo) error
	SaveBookStock(ctx context.Context, stocks ...do.BookStock) error
}

type BookStockRepo struct {
	bookDao   BookStockDao
	infoDao   BookInfoDao
	bookCache BookCache
}

func NewBookStockRepo(bookDao BookStockDao, infoDao BookInfoDao, bookCache BookCache) *BookStockRepo {
	return &BookStockRepo{
		bookDao:   bookDao,
		infoDao:   infoDao,
		bookCache: bookCache,
	}
}

func (b *BookStockRepo) CheckBookInfoIfExist(ctx context.Context, name, author, publisher, category string) (uint64, bool) {
	return b.bookDao.CheckBookIfExist(ctx, name, author, publisher, category)
}

func (b *BookStockRepo) AddBookStock(ctx context.Context, bookID uint64, userID *uint64, num uint) error {
	err := b.bookCache.DeleteBookStock(ctx, bookID)
	if err != nil {
		return err
	}

	//异步删除缓存
	defer func() {
		go time.AfterFunc(time.Second, func() {
			_ = b.bookCache.DeleteBookStock(ctx, bookID)
		})
	}()
	var donate bool

	var user uint64
	if userID != nil {
		donate = true
		user = *userID
	} else {
		donate = false
	}

	return b.bookDao.AddBookStock(ctx, bookID, user, num, donate)
}

func (b *BookStockRepo) RegisterBookAndAddBookStock(ctx context.Context, bookID uint64, userID *uint64, book service.BookInfo, num uint) error {
	clean := func(ctx context.Context) error {
		if err := b.bookCache.DeleteBookInfo(ctx, bookID); err != nil {
			return err
		}
		if err := b.bookCache.DeleteBookStock(ctx, bookID); err != nil {
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

	var donate bool

	var user uint64
	if userID != nil {
		donate = true
		user = *userID
	} else {
		donate = false
	}

	return b.bookDao.RegisterAndAddBookStock(ctx, do.BookInfo{
		ID:        bookID,
		Name:      book.Name,
		Author:    book.Author,
		Publisher: book.Publisher,
		Category:  book.Category,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}, user, num, donate)
}

func (b *BookStockRepo) FuzzyQueryBook(ctx context.Context, pageSize int, currentPage int, totalNum *int, opts ...func(db *gorm.DB)) ([]service.Book, error) {
	num, err := b.bookDao.GetBookTotalNum(ctx, opts...)
	if err != nil {
		return nil, err
	}
	*totalNum = num
	if maxPage := tool.GetPage(num, pageSize); currentPage > maxPage {
		return nil, errcode.PageError
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

func (b *BookStockRepo) getBookInID(ctx context.Context, ids ...uint64) ([]do.BookInfo, []do.BookStock, error) {
	bookInfos, leftID := b.bookCache.GetBookInfoByID(ctx, ids...)

	if len(leftID) > 0 {
		tmpBookInfos, err := b.infoDao.GetBookInfoByID(ctx, leftID...)
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
