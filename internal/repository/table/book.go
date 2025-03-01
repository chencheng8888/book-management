package table

import (
	"book-management/internal/pkg/common"
	"time"
)

const BookTableName = "books"

// 绘本信息
type Book struct {
	ID        uint64    `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name"`
	Author    string    `gorm:"column:author"`
	Publisher string    `gorm:"column:publisher"`
	Category  string    `gorm:"column:category"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (b Book) TableName() string {
	return BookTableName
}

func (b Book) IsNormal() bool {
	return b.Status == common.BookStatusNormal
}
func (b Book) IsReturned() bool {
	return b.Status == common.BookStatusReturned
}
func (b Book) IsOverdue() bool {
	return b.Status == common.BookStatusOverdue
}

const BookStockTableName = "book_stocks"

// 绘本库存
type BookStock struct {
	BookID    uint64    `gorm:"primaryKey;column:book_id"`
	Stock     uint      `gorm:"column:stock"`
	Status    string    `gorm:"column:status"`
	Where     string    `gorm:"column:where"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (s BookStock) IsAdequate() bool {
	return s.Status == common.Adequate
}
func (s BookStock) IsEarlyWarning() bool {
	return s.Status == common.EarlyWarning
}
func (s BookStock) IsShortage() bool {
	return s.Status == common.Shortage
}
func (s BookStock) TableName() string {
	return BookStockTableName
}
