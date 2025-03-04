package do

import (
	"book-management/internal/pkg/common"
	"time"
)

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
	return common.BookTableName
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
	return common.BookStockTableName
}

// 书籍副本信息
type BookCopy struct {
	BookID uint64 `gorm:"column:book_id;uniqueIndex:idx_book_copy"` // 书籍ID[相当于标记某一种书，比如《高等数学》]
	CopyID uint64 `gorm:"column:copy_id;uniqueIndex:idx_book_copy"` // 副本ID [这个相当于标识特定的某一本书，比如《高等数学》的第一本]
	Status string `gorm:"column:status"`                            // 借阅状态
}

func (bc BookCopy) TableName() string {
	return common.BookCopyTableName
}
func (bc BookCopy) IsNormal() bool {
	return bc.Status == common.BookStatusNormal
}
func (bc BookCopy) IsReturned() bool {
	return bc.Status == common.BookStatusReturned
}
func (bc BookCopy) IsOverdue() bool {
	return bc.Status == common.BookStatusOverdue
}

// 书籍借阅记录
type BookBorrow struct {
	BookID uint64 `gorm:"column:book_id;uniqueIndex:idx_book_copy"` // 书籍ID[相当于标记某一种书，比如《高等数学》]
	CopyID uint64 `gorm:"column:copy_id;uniqueIndex:idx_book_copy"` // 副本ID [这个相当于标识特定的某一本书，比如《高等数学》的第一本]

	BorrowerID         string     `gorm:"column:borrower_id"`          // 借阅者ID
	ExpectedReturnTime time.Time  `gorm:"column:expected_return_time"` // 预计归还时间
	CreatedTime        time.Time  `gorm:"column:created_time"`         // 借阅时间
	ReturnTime         *time.Time `gorm:"column:return_time"`          // 归还时间
}

func (bb BookBorrow) TableName() string {
	return common.BookBorrowTableName
}
