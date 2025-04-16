package do

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"book-management/internal/pkg/common"
)

// 绘本信息
type BookInfo struct {
	ID        uint64 `gorm:"primaryKey;column:id"`
	Name      string `gorm:"type:varchar(255);column:name;index"`
	Author    string `gorm:"type:varchar(100);column:author;index"`
	Publisher string `gorm:"type:varchar(255);column:publisher"`
	Category  string `gorm:"type:varchar(50);column:category;index"`
	//Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (b BookInfo) TableName() string {
	return common.BookTableName
}

// 绘本库存
type BookStock struct {
	BookID uint64 `gorm:"primaryKey;column:book_id"`
	Stock  uint   `gorm:"column:stock"`
	//Status    string    `gorm:"column:status"`
	// Where     string    `gorm:"type:varchar(255);column:where"`
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
	Status bool   `gorm:"column:status"`                            // 是否在库中
	Donate bool   `gorm:"column:donate"`                            //是否捐赠
}

func (b *BookCopy) BeforeCreate(tx *gorm.DB) (err error) {
	if b.CopyID == 0 { // 只有当CopyID未设置时才自动生成
		var maxCopyID uint64
		err = tx.Model(&BookCopy{}).
			Where("book_id = ?", b.BookID).
			Select("IFNULL(MAX(copy_id), 0)").
			Scan(&maxCopyID).Error
		if err != nil {
			return fmt.Errorf("failed to query max copy_id: %w", err)
		}

		b.CopyID = maxCopyID + 1 // 从1开始计数
	}
	return nil
}
func (b *BookCopy) TableName() string {
	return common.BookCopyTableName
}

// 书籍借阅记录
type BookBorrow struct {
	BookID uint64 `gorm:"column:book_id"` // 书籍ID[相当于标记某一种书，比如《高等数学》]
	CopyID uint64 `gorm:"column:copy_id"` // 副本ID [这个相当于标识特定的某一本书，比如《高等数学》的第一本]

	BorrowerID         uint64     `gorm:"column:borrower_id"`             // 借阅者ID
	ExpectedReturnTime time.Time  `gorm:"column:expected_return_time"`    // 预计归还时间
	CreatedTime        time.Time  `gorm:"column:created_time"`            // 借阅时间
	ReturnTime         *time.Time `gorm:"column:return_time"`             // 归还时间
	Status             string     `gorm:"type:varchar(50);column:status"` // 借阅状态
}

func (bb BookBorrow) TableName() string {
	return common.BookBorrowTableName
}
