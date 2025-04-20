package do

import (
	"book-management/internal/pkg/common"
	"time"
)

type Activity struct {
	ID        uint64    `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"type:varchar(255);column:name"`
	Type      string    `gorm:"type:varchar(255);column:type"`
	StartTime time.Time `gorm:"column:start_time"`
	EndTime   time.Time `gorm:"column:end_time"`
	Manager   string    `gorm:"type:varchar(255);column:manager"`
	Phone     string    `gorm:"type:varchar(30);column:phone"`
	Addr      string    `gorm:"type:varchar(255);column:addr"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (a Activity) TableName() string {
	return common.ActivityTableName
}
