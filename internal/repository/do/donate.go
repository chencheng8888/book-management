package do

import (
	"time"

	"book-management/internal/pkg/common"
)

type DonateInfo struct {
	UserID    uint64    `gorm:"column:user_id"`
	BookID    uint64    `gorm:"column:book_id"`
	Num       int       `gorm:"column:num"` //捐赠数量
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (d DonateInfo) TableName() string {
	return common.DonateTableName
}
