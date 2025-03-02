package repository

import (
	"book-management/internal/pkg/common"
	"book-management/internal/repository/do"
	"book-management/internal/service"
	"context"
	"gorm.io/gorm"
	"math"
	"time"
)

type BookDao interface {
	CheckBookIfExist(ctx context.Context, name, author, publisher, category string) (uint64, bool)
	AddBookStock(ctx context.Context, id uint64, num uint, where *string) error
	RegisterAndAddBookStock(ctx context.Context, bookInfo do.BookInfo, addedNum uint, where *string) error

	GetBookInfoByID(ctx context.Context, id ...uint64) ([]do.BookInfo, error)
	GetBookStockByID(ctx context.Context, id ...uint64) ([]do.BookStock, error)

	FuzzyQueryBookID(ctx context.Context, pageSize int, page int, opts ...func(db *gorm.DB)) ([]uint64, error)
	GetBookTotalNum(ctx context.Context) (int, error)
}

type BookCache interface {
	DeleteBookStock(ctx context.Context, id uint64) error
	DeleteBookInfo(ctx context.Context, id uint64) error
	DeleteBookTotalNum(ctx context.Context) error
	GetBookInfoByID(ctx context.Context, ids ...uint64) ([]do.BookInfo, []uint64)
	GetBookStockByID(ctx context.Context, ids ...uint64) ([]do.BookStock, []uint64)
	GetBookTotalNum(ctx context.Context) (int, error)
	SaveBookInfo(ctx context.Context, infos ...do.BookInfo) error
	SaveBookStock(ctx context.Context, stocks ...do.BookStock) error
	SaveBookTotalNum(ctx context.Context, num int) error
}

type BookRepo struct {
	bookDao   BookDao
	bookCache BookCache
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

func (b *BookRepo) RegisterBookAndAddBookStock(ctx context.Context, id uint64, book service.BookInfo, num uint, where *string) error {
	clean := func(ctx context.Context) error {
		if err := b.bookCache.DeleteBookInfo(ctx, id); err != nil {
			return err
		}
		if err := b.bookCache.DeleteBookStock(ctx, id); err != nil {
			return err
		}
		if err := b.bookCache.DeleteBookTotalNum(ctx); err != nil {
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
	var num int
	var err error
	var pageNum int
	num, err = b.bookCache.GetBookTotalNum(ctx)
	if err != nil {
		num, err = b.bookDao.GetBookTotalNum(ctx)
		if err != nil {
			return nil, err
		}

		pageNum = int(math.Ceil(float64(num) / float64(pageSize)))

		totalPage = &pageNum

		//异步保存书籍总数
		defer func() {
			go b.bookCache.SaveBookTotalNum(ctx, num)
		}()
	}

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
