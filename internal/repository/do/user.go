package do

import (
	"book-management/internal/pkg/common"
	"time"
)

type User struct {
	ID        string    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"column:name"`
	Phone     string    `gorm:"column:phone;uniqueIndex"`
	Integral  uint      `gorm:"column:integral"`
	Gender    string    `gorm:"column:gender"`
	IsVip     bool      `gorm:"column:is_vip"`
	VipLevels *string   `gorm:"column:vip_levels"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (u *User) TableName() string {
	return common.UserTableName
}
