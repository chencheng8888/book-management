package do

import (
	"time"
)

const BookTableName = "books"

// 绘本信息
type BookInfo struct {
	ID        uint64 `gorm:"primaryKey;column:id"`
	Name      string `gorm:"column:name"`
	Author    string `gorm:"column:author"`
	Publisher string `gorm:"column:publisher"`
	Category  string `gorm:"column:category"`
	//Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (b BookInfo) TableName() string {
	return BookTableName
}

//func (b BookInfo) IsNormal() bool {
//	return b.Status == common.BookStatusNormal
//}
//func (b BookInfo) IsReturned() bool {
//	return b.Status == common.BookStatusReturned
//}
//func (b BookInfo) IsOverdue() bool {
//	return b.Status == common.BookStatusOverdue
//}

const BookStockTableName = "book_stocks"

// 绘本库存
type BookStock struct {
	BookID uint64 `gorm:"primaryKey;column:book_id"`
	Stock  uint   `gorm:"column:stock"`
	//Status    string    `gorm:"column:status"`
	Where     string    `gorm:"column:where"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (s BookStock) IsAdequate() bool {
	return s.Stock > 20
}
func (s BookStock) IsEarlyWarning() bool {
	return s.Stock <= 20 && s.Stock > 0
}
func (s BookStock) IsShortage() bool {
	return s.Stock == 0
}
func (s BookStock) TableName() string {
	return BookStockTableName
}
