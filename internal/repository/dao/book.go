package dao

import (
	"book-management/internal/pkg/common"
	"book-management/internal/repository/do"
	"book-management/pkg/logger"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type BookDao struct {
	db *gorm.DB
}

func (b *BookDao) AddBookBorrowRecord(ctx context.Context, bookID uint64, borrowerID string, expectedReturnTime time.Time, copyID *uint64) error {
	var copy_id uint64
	err := b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Table(common.BookCopyTableName).Select("copy_id").Where("book_id = ? AND status = ?", bookID, true).Scan(&copy_id).Error
		if err != nil {
			return err
		}
		copyID = &copy_id
		err = updateCopyStatus(ctx, tx, bookID, copy_id, false)
		if err != nil {
			return err
		}
		err = tx.Debug().Table(common.BookBorrowTableName).Create(&do.BookBorrow{
			BookID:             bookID,
			CopyID:             copy_id,
			BorrowerID:         borrowerID,
			ExpectedReturnTime: expectedReturnTime,
			CreatedTime:        time.Now(),
			ReturnTime:         nil,
			Status:             common.BookStatusWaitingReturn,
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.LogPrinter.Errorf("DB:AddBookBorrowRecord [bookID:%v borrowerID:%v expectedReturnTime:%v copyID:%v] failed, err: %v", bookID, borrowerID, expectedReturnTime, copyID, err)
	}
	return err
}

func (b *BookDao) GetBookRecordTotalNum(ctx context.Context, opts ...func(db *gorm.DB)) (int, error) {
	var num int64
	db := b.db.WithContext(ctx).Table(common.BookBorrowTableName)
	for _, opt := range opts {
		opt(db)
	}
	err := db.Debug().Count(&num).Error
	if err != nil {
		return 0, err
	}
	return int(num), nil
}

func (b *BookDao) FuzzyQueryBookBorrowRecord(ctx context.Context, pageSize int, page int, opts ...func(db *gorm.DB)) ([]do.BookBorrow, error) {
	db := b.db.WithContext(ctx).Table(common.BookBorrowTableName)
	for _, opt := range opts {
		opt(db)
	}
	var records []do.BookBorrow
	err := db.Debug().Limit(pageSize).Offset((page - 1) * pageSize).Order("book_id ASC,copy_id ASC").Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
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
		logger.LogPrinter.Errorf("DB:AddBookStock [id:%v num:%v where:%v] failed, err: %v", id, num, where, err)
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
		logger.LogPrinter.Errorf("DB:RegisterAndAddBookStock [bookInfo:%v addedNum:%v where:%v] failed, err: %v", bookInfo, addedNum, where, err)
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

	// 基础查询：联表查询
	db := b.db.WithContext(ctx).Table(common.BookTableName).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.id = %s.book_id", common.BookStockTableName, common.BookTableName, common.BookStockTableName)).
		Select("DISTINCT id") // 确保ID唯一

	for _, opt := range opts {
		opt(db)
	}

	var ids []uint64

	err := db.Debug().
		Where("id >= (?)", db.Order("id ASC").Offset((page-1)*pageSize).Limit(1)).
		Order("id ASC").Limit(pageSize).Find(&ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (b *BookDao) GetBookTotalNum(ctx context.Context, opts ...func(db *gorm.DB)) (int, error) {

	// 基础查询：联表查询
	db := b.db.WithContext(ctx).Table(common.BookTableName).
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.id = %s.book_id", common.BookStockTableName, common.BookTableName, common.BookStockTableName)).
		Select("DISTINCT id") // 确保ID唯一

	for _, opt := range opts {
		opt(db)
	}

	var cnt int64

	err := db.Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return int(cnt), nil
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
	err := db.WithContext(ctx).Debug().Table(common.BookTableName).
		Where("id in (?)", ids).Find(&infos).Error
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func getBookStockByID(ctx context.Context, db *gorm.DB, ids ...uint64) ([]do.BookStock, error) {
	var stocks []do.BookStock
	err := db.WithContext(ctx).Debug().Table(common.BookStockTableName).
		Where("book_id in (?)", ids).Find(&stocks).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func updateCopyStatus(ctx context.Context, db *gorm.DB, bookID uint64, copyID uint64, status bool) error {
	return db.WithContext(ctx).Debug().Table(common.BookCopyTableName).
		Where("book_id =? AND copy_id =?", bookID, copyID).
		Update("status", status).Error
}
