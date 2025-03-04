package dao

import (
	"book-management/internal/pkg/common"
	"book-management/internal/repository/do"
	"context"
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

	err = testDB.AutoMigrate(&do.BookInfo{}, &do.BookStock{})
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

func Test_getBookTotalNum(t *testing.T) {
	num, err := getBookTotalNum(context.Background(), testDB)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(num)
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
}
