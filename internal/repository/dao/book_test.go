package dao

import (
	"book-management/internal/pkg/common"
	"book-management/internal/repository/do"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = gorm.Open(mysql.Open("root:12345678@tcp(127.0.0.1:13306)/bookTest?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = testDB.AutoMigrate(&do.BookInfo{}, &do.BookStock{}, &do.BookBorrow{}, &do.BookCopy{})
	if err != nil {
		panic(err)
	}
	m.Run()
}

func Test_createBookInfo(t *testing.T) {
	testTime := time.Now()
	err := createBookInfo(context.Background(), testDB, do.BookInfo{
		ID:        111,
		Name:      "test",
		Author:    "test",
		Publisher: "test",
		Category:  common.ChildrenStory,
		CreatedAt: testTime,
		UpdatedAt: testTime,
	})
	if err != nil {
		t.Error(err)
	}
}

func Test_createBookStock(t *testing.T) {
	testTime := time.Now()
	err := createBookStock(context.Background(), testDB, do.BookStock{
		BookID:    111,
		Stock:     10,
		Where:     "test",
		CreatedAt: testTime,
		UpdatedAt: testTime,
	})
	if err != nil {
		t.Error(err)
	}
}

func Test_getBookInfoByID(t *testing.T) {
	info, err := getBookInfoByID(context.Background(), testDB, 111)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}

func Test_getBookStockByID(t *testing.T) {
	stock, err := getBookStockByID(context.Background(), testDB, 111)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(stock)
}

func Test_updateStockWhere(t *testing.T) {
	err := updateStockWhere(context.Background(), testDB, 111, "test_where1")
	if err != nil {
		t.Error(err)
	}
}

func Test_addStock(t *testing.T) {
	err := addStock(context.Background(), testDB, 111, 10)
	if err != nil {
		t.Error(err)
	}
}

func Test_checkBookIfExist(t *testing.T) {
	id, exist := checkBookIfExist(context.Background(), testDB, "test", "test", "test", common.ChildrenStory)
	if exist {
		t.Log(id)
	} else {
		t.Log("not exist")
	}
}

func Test_checkBookStockIfExistByID(t *testing.T) {
	exist := checkBookStockIfExistByID(context.Background(), testDB, 111)
	if exist {
		t.Log("exist")
	} else {
		t.Log("not exist")
	}
}

func TestBookDao_FuzzyQueryBookID(t *testing.T) {
	bookDao := &BookDao{db: testDB}
	t.Run("no query", func(t *testing.T) {
		ids, err := bookDao.FuzzyQueryBookID(context.Background(), 10, 1)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ids)
	})
	t.Run("have query", func(t *testing.T) {
		ids, err := bookDao.FuzzyQueryBookID(context.Background(), 10, 1, func(db *gorm.DB) {
			db = db.Where(fmt.Sprintf("%s.name = ?", common.BookTableName), "test")
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ids)
	})
	t.Run("have many query", func(t *testing.T) {
		ids, err := bookDao.FuzzyQueryBookID(context.Background(), 10, 1, func(db *gorm.DB) {
			db = db.Where(fmt.Sprintf("%s.author = ?", common.BookTableName), "test")
		}, func(db *gorm.DB) {
			db = db.Where(fmt.Sprintf("%s.where = ?", common.BookStockTableName), "test_where1")
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(ids)
	})
}

func TestBookDao_GetBookTotalNum(t *testing.T) {
	bookDao := &BookDao{db: testDB}
	t.Run("no query", func(t *testing.T) {
		nums, err := bookDao.GetBookTotalNum(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		t.Log(nums)
	})
	t.Run("have query", func(t *testing.T) {
		nums, err := bookDao.GetBookTotalNum(context.Background(), func(db *gorm.DB) {
			db = db.Where(fmt.Sprintf("%s.name = ?", common.BookTableName), "test")
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(nums)
	})
	t.Run("have many  query", func(t *testing.T) {
		nums, err := bookDao.GetBookTotalNum(context.Background(), func(db *gorm.DB) {
			db = db.Where(fmt.Sprintf("%s.author = ?", common.BookTableName), "test")
		}, func(db *gorm.DB) {
			db = db.Where(fmt.Sprintf("%s.where = ?", common.BookStockTableName), "test_where1")
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(nums)
	})
}

func TestBookDao_AddBookBorrowRecord(t *testing.T) {
	bookDao := &BookDao{db: testDB}

	var copyID uint64

	err := bookDao.AddBookBorrowRecord(context.Background(), 111, "test", time.Now(), &copyID)
	if err != nil {
		t.Error(err)
	}
}

func TestBookDao_GetBookRecordTotalNum(t *testing.T) {
	bookDao := &BookDao{db: testDB}

	num, err := bookDao.GetBookRecordTotalNum(context.Background(), func(db *gorm.DB) {
		db.Where("status = ?", common.BookStatusWaitingReturn)
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(num)
}

func TestBookDao_UpdateBorrowStatus(t *testing.T) {
	bookDao := &BookDao{db: testDB}

	err := bookDao.UpdateBorrowStatus(context.Background(), 111, 2, common.BookStatusOverdue)
	if err != nil {
		t.Error(err)
	}
}
