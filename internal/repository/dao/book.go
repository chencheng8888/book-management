package dao

import (
	"book-management/internal/pkg/common"
	"book-management/internal/repository/do"
	"book-management/internal/repository/repo"
	"book-management/pkg/logger"
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BookDao struct {
	db          *gorm.DB
	fixedPoints int
}

func (b *BookDao) GetAvailableBookCopy(ctx context.Context, bookID uint64, page, pageSize int) ([]uint64, error) {
	var res []uint64
	err := b.db.WithContext(ctx).Table(common.BookCopyTableName).
		Where("book_id = ? and status = ?", bookID, true).Order("copy_id").Offset((page-1)*pageSize).
		Limit(pageSize+1).Pluck("copy_id", &res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewBookDao(db *gorm.DB) *BookDao {
	return &BookDao{
		db:          db,
		fixedPoints: 100,
	}
}
func (b *BookDao) GetBookDonateRecordsNum(ctx context.Context) (int, error) {
	db := b.db.WithContext(ctx).Table(common.DonateTableName)
	var cnt int64
	err := db.Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return int(cnt), nil
}

func (b *BookDao) GetBookDonateInfos(ctx context.Context, pageSize int, page int) ([]do.DonateInfo, error) {
	db := b.db.WithContext(ctx).Table(common.DonateTableName)
	var donateInfos []do.DonateInfo
	err := db.Debug().Limit(pageSize).Offset((page - 1) * pageSize).Order("created_at DESC").Find(&donateInfos).Error
	if err != nil {
		return nil, err
	}
	return donateInfos, nil
}

func (b *BookDao) GetDonateRanking(ctx context.Context, top int) ([]repo.DonateRanking, error) {
	db := b.db.WithContext(ctx).Table(common.DonateTableName)
	var donateRanking []repo.DonateRanking
	err := db.Debug().Select(fmt.Sprintf("%s.user_id,%s.name as user_name,sum(%s.num) as donate_num,count(*) as donate_times,MAX(%s.updated_at) as updated_at",
		common.DonateTableName, common.UserTableName, common.DonateTableName, common.DonateTableName)).
		Joins(fmt.Sprintf("left join %s on %s.user_id = %s.id", common.UserTableName, common.DonateTableName, common.UserTableName)).
		Group(fmt.Sprintf("%s.user_id, %s.name",
			common.DonateTableName, common.UserTableName)).
		Order("donate_num DESC").Limit(top).Find(&donateRanking).Error

	if err != nil {
		return nil, err
	}
	return donateRanking, nil
}
func (b *BookDao) GetBookBorrowStatistics(ctx context.Context, startTime, endTime time.Time) (do.BorrowStatistics, error) {
	db := b.db.WithContext(ctx).Table(common.BookBorrowTableName)
	type tmp struct {
		Category string
		Num      int
	}
	var res []tmp
	err := db.Debug().Select("category, COUNT(*) num").Joins(fmt.Sprintf("JOIN %s ON %s.book_id = %s.id", common.BookTableName, common.BookBorrowTableName, common.BookTableName)).
		Group("category").Find(&res).Error
	if err != nil {
		return do.BorrowStatistics{}, err
	}

	var ans do.BorrowStatistics
	for _, v := range res {
		if v.Category == common.ChildrenStory {
			ans.ChildrenStoryNum = v.Num
		} else if v.Category == common.ScienceKnowledge {
			ans.ScienceKnowledgeNum = v.Num
		} else if v.Category == common.ArtEnlightenment {
			ans.ArtEnlightenmentNum = v.Num
		}
	}
	return ans, nil
}

func (b *BookDao) UpdateBorrowStatus(ctx context.Context, bookID, copyID uint64, status string) error {
	err := b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var oldStatus string
		err := tx.WithContext(ctx).Debug().Table(common.BookBorrowTableName).Pluck("status", &oldStatus).Error
		if err != nil {
			return err
		}
		if oldStatus == "" {
			return errors.New("data is not exist")
		}
		//如果状态相同就直接返回
		if oldStatus == status {
			return nil
		}

		//然后更新状态信息
		err = updateBorrowRecordStatus(ctx, tx, bookID, copyID, status)
		if err != nil {
			return err
		}

		IfInStock := false

		var returnTime *time.Time

		if status == common.BookStatusReturned {
			IfInStock = true
			returnTime = new(time.Time)
			*returnTime = time.Now()
		}
		err = updateCopyStatus(ctx, tx, bookID, copyID, IfInStock)
		if err != nil {
			return err
		}

		return updateBorrowRecordReturnTime(ctx, tx, bookID, copyID, returnTime)
	})
	if err != nil {
		logger.LogPrinter.Errorf("DB: update borrow status[book_id:%v copy_id:%v status:%v] failed: %v", bookID, copyID, status, err)
	}
	return err
}

func (b *BookDao) AddBookBorrowRecord(ctx context.Context, bookID uint64, borrowerID uint64, expectedReturnTime time.Time, copyID uint64) error {
	err := b.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//先检查用户是否存在
		exist, err := checkUserIfExist(ctx, tx, borrowerID)
		if err != nil {
			return err
		}
		if !exist {
			return errors.New("user is not exist")
		}

		//添加积分
		err = addUserIntegral(ctx, tx, borrowerID, b.fixedPoints)
		if err != nil {
			return err
		}

		//检查对应书籍是否存在
		exist = checkBookCopyExist(ctx, tx, bookID, copyID)
		if !exist {
			return errors.New("data is not exist")
		}

		//更新书籍借阅状态
		err = updateCopyStatus(ctx, tx, bookID, copyID, false)
		if err != nil {
			return err
		}
		//添加借阅记录
		err = tx.Debug().Table(common.BookBorrowTableName).Create(&do.BookBorrow{
			BookID:             bookID,
			CopyID:             copyID,
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
		logger.LogPrinter.Errorf("DB:Add Book Borrow Record [bookID:%v borrowerID:%v expectedReturnTime:%v copyID:%v] failed, err: %v", bookID, borrowerID, expectedReturnTime, copyID, err)
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

func (b *BookDao) AddBookStock(ctx context.Context, bookID, userID uint64, num uint, donate bool) error {
	err := b.db.Transaction(func(tx *gorm.DB) error {
		if !checkBookStockIfExistByID(ctx, tx, bookID) {
			return errors.New("book stock is not exist")
		}
		if err := addStock(ctx, tx, bookID, num); err != nil {
			return err
		}
		//添加copy
		if err := createBookCopy(ctx, tx, bookID, num, donate); err != nil {
			return err
		}

		//创建捐赠记录
		if donate {
			if err := createDonateRecord(ctx, tx, bookID, userID, int(num)); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		logger.LogPrinter.Errorf("DB:Add Book Stock [id:%v num:%v] failed, err: %v", bookID, num, err)
	}
	return err
}

func (b *BookDao) RegisterAndAddBookStock(ctx context.Context, bookInfo do.BookInfo, userID uint64, addedNum uint, donate bool) error {
	createdTime := time.Now()
	err := b.db.Transaction(func(tx *gorm.DB) error {

		if err := createBookInfo(ctx, tx, bookInfo); err != nil {
			return err
		}

		if checkBookStockIfExistByID(ctx, tx, bookInfo.ID) {
			return errors.New("book stock unexpectedly exist")
		}

		if err := createBookStock(ctx, tx, do.BookStock{
			BookID:    bookInfo.ID,
			Stock:     addedNum,
			CreatedAt: createdTime,
			UpdatedAt: createdTime,
		}); err != nil {
			return err
		}
		//添加copy
		if err := createBookCopy(ctx, tx, bookInfo.ID, addedNum, donate); err != nil {
			return err
		}
		//创建捐赠记录
		if donate {
			if err := createDonateRecord(ctx, tx, bookInfo.ID, userID, int(addedNum)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		logger.LogPrinter.Errorf("DB:Register And Add Book Stock [bookInfo:%v addedNum:%v] failed, err: %v", bookInfo, addedNum, err)
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

	db := b.db.WithContext(ctx).Table(common.BookTableName).
		Select("DISTINCT id")

	for _, opt := range opts {
		opt(db)
	}

	var ids []uint64

	offset := (page - 1) * pageSize

	err := db.Debug().
		Order("id ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&ids).Error

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (b *BookDao) GetBookTotalNum(ctx context.Context, opts ...func(db *gorm.DB)) (int, error) {

	// 基础查询：联表查询
	db := b.db.WithContext(ctx).Table(common.BookTableName).
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

func checkBookCopyExist(ctx context.Context, db *gorm.DB, bookID, copyID uint64) bool {
	var count int64
	err := db.WithContext(ctx).Debug().Table(common.BookCopyTableName).
		Where("book_id = ? AND copy_id = ?", bookID, copyID).Count(&count).Error
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
	res := db.WithContext(ctx).Debug().Table(do.BookStock{}.TableName()).
		Where("book_id = ?", id).Update("where", where)
	if res.Error != nil {
		return res.Error
	}
	//if res.RowsAffected == 0 {
	//	return errors.New("no data updated")
	//}
	return nil
}

func createBookInfo(ctx context.Context, db *gorm.DB, bookInfo do.BookInfo) error {
	return db.WithContext(ctx).Debug().Create(&bookInfo).Error
}

func createBookStock(ctx context.Context, db *gorm.DB, bookStock do.BookStock) error {
	return db.WithContext(ctx).Debug().Create(&bookStock).Error
}

func createBookCopy(ctx context.Context, db *gorm.DB, bookID uint64, num uint, donate bool) error {
	for i := uint(0); i < num; i++ {
		err := db.WithContext(ctx).Debug().Table(common.BookCopyTableName).
			Create(&do.BookCopy{
				BookID: bookID,
				Status: true,
				Donate: donate,
			}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func getBookInfoByID(ctx context.Context, db *gorm.DB, ids ...uint64) ([]do.BookInfo, error) {
	var infos []do.BookInfo
	res := db.WithContext(ctx).Debug().Table(common.BookTableName).
		Where("id in (?)", ids).Find(&infos)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected != int64(len(ids)) {
		return nil, errors.New("some datas don't found")
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
	err := db.WithContext(ctx).Debug().Table(common.BookCopyTableName).
		Where("book_id =? AND copy_id =?", bookID, copyID).
		Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}
func updateBorrowRecordStatus(ctx context.Context, db *gorm.DB, bookID, copyID uint64, status string) error {
	err := db.WithContext(ctx).Debug().Table(common.BookBorrowTableName).Where("book_id = ? AND copy_id = ?", bookID, copyID).
		Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func updateBorrowRecordReturnTime(ctx context.Context, db *gorm.DB, bookID, copyID uint64, returnTime *time.Time) error {
	err := db.Debug().Table(common.BookBorrowTableName).Where("book_id = ? AND copy_id = ?", bookID, copyID).
		Update("return_time", returnTime).Error
	if err != nil {
		return err
	}
	return nil
}

func createDonateRecord(ctx context.Context, db *gorm.DB, bookID, userID uint64, num int) error {
	donateInfo := do.DonateInfo{
		BookID: bookID,
		UserID: userID,
		Num:    num,
	}
	return db.WithContext(ctx).Debug().Create(&donateInfo).Error
}
