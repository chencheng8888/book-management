package dao

import (
	"book-management/internal/repository/do"
	"book-management/pkg/logger"
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type BookDao struct {
	db *gorm.DB
}

func (b *BookDao) CheckBookIfExist(ctx context.Context, name, author, publisher, category string) (uint64, bool) {
	return checkBookIfExist(ctx, b.db, name, author, publisher, category)
}

func (b *BookDao) AddBookStock(ctx context.Context, id uint64, num uint, where *string) error {
	err := b.db.Transaction(func(tx *gorm.DB) error {
		if !checkBookStockIfExistByID(ctx, tx, id) {
			return errors.New("book stock is not exist")
		}
		if err := addStock(ctx, tx, id, num); err != nil {
			return err
		}
		if where != nil {
			if err := updateStockWhere(ctx, tx, id, *where); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.LogPrinter.Errorf("AddBookStock [id:%v num:%v where:%v] failed, err: %v", id, num, where, err)
	}
	return err
}

func (b *BookDao) RegisterAndAddBookStock(ctx context.Context, bookInfo do.BookInfo, addedNum uint, where string) error {
	createdTime := time.Now()
	err := b.db.Transaction(func(tx *gorm.DB) error {
		err := createBookInfo(ctx, tx, bookInfo)
		if err != nil {
			return err
		}

		if checkBookStockIfExistByID(ctx, tx, bookInfo.ID) {
			return errors.New("book stock unexpectedly exist")
		}

		err = createBookStock(ctx, tx, do.BookStock{
			BookID:    bookInfo.ID,
			Stock:     addedNum,
			Where:     where,
			CreatedAt: createdTime,
			UpdatedAt: createdTime,
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.LogPrinter.Errorf("RegisterAndAddBookStock [bookInfo:%v addedNum:%v where:%v] failed, err: %v", bookInfo, addedNum, where, err)
	}
	return err
}

func (b *BookDao) GetBookInfoByID(ctx context.Context, id ...uint64) ([]do.BookInfo, error) {
	if len(id) == 0 {
		return nil, errors.New("the length of the ID to be queried is 0")
	}
	return getBookInfoByID(ctx, b.db, id...)
}

func (b *BookDao) GetBookStockByID(ctx context.Context, id ...uint64) ([]do.BookStock, error) {
	if len(id) == 0 {
		return nil, errors.New("the length of the ID to be queried is 0")
	}
	return getBookStockByID(ctx, b.db, id...)
}

func (b *BookDao) FuzzyQueryBookID(ctx context.Context, pageSize int, page int, opts ...func(db *gorm.DB)) ([]uint64, error) {
	db := b.db.WithContext(ctx).Table(do.BookInfo{}.TableName())

	for _, opt := range opts {
		opt(db)
	}

	var ids []uint64

	err := db.Debug().Select("id").
		Where("id >= (?)", db.Select("id").Order("id ASC").Offset((page-1)*pageSize).Limit(1)).
		Order("id ASC").Limit(pageSize).Find(&ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (b *BookDao) GetBookTotalNum(ctx context.Context) (int, error) {
	return getBookTotalNum(ctx, b.db)
}

func checkBookIfExist(ctx context.Context, db *gorm.DB, name, author, publisher, category string) (uint64, bool) {
	var id uint64
	err := db.WithContext(ctx).Debug().Table(do.BookInfo{}.TableName()).
		Where("name = ? AND author = ? AND publisher = ? AND category = ?", name, author, publisher, category).Pluck("id", &id).Error
	if err != nil || id == 0 {
		return 0, false
	}
	return id, true
}

func checkBookStockIfExistByID(ctx context.Context, db *gorm.DB, id uint64) bool {
	var count int64
	err := db.WithContext(ctx).Debug().Table(do.BookStock{}.TableName()).
		Where("book_id = ?", id).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func addStock(ctx context.Context, db *gorm.DB, id uint64, num uint) error {
	return db.WithContext(ctx).Debug().Table(do.BookStock{}.TableName()).
		Where("book_id = ?", id).
		Update("stock", gorm.Expr("stock + ?", num)).Error
}

func updateStockWhere(ctx context.Context, db *gorm.DB, id uint64, where string) error {
	return db.WithContext(ctx).Debug().Table(do.BookStock{}.TableName()).
		Where("book_id = ?", id).Update("where", where).Error
}

func createBookInfo(ctx context.Context, db *gorm.DB, bookInfo do.BookInfo) error {
	return db.WithContext(ctx).Debug().Create(&bookInfo).Error
}

func createBookStock(ctx context.Context, db *gorm.DB, bookStock do.BookStock) error {
	return db.WithContext(ctx).Debug().Create(&bookStock).Error
}

func getBookInfoByID(ctx context.Context, db *gorm.DB, ids ...uint64) ([]do.BookInfo, error) {
	var infos []do.BookInfo
	err := db.WithContext(ctx).Debug().Table(do.BookInfo{}.TableName()).
		Where("id in (?)", ids).Find(&infos).Error
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func getBookStockByID(ctx context.Context, db *gorm.DB, ids ...uint64) ([]do.BookStock, error) {
	var stocks []do.BookStock
	err := db.WithContext(ctx).Debug().Table(do.BookStock{}.TableName()).
		Where("book_id in (?)", ids).Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func getBookTotalNum(ctx context.Context, db *gorm.DB) (int, error) {
	var count int64
	err := db.WithContext(ctx).Debug().Table(do.BookInfo{}.TableName()).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
