package do

import (
	"book-management/internal/pkg/common"
	"time"
)

type User struct {
	ID        uint64    `gorm:"column:id;primaryKey"`
	Name      string    `gorm:"type:varchar(255);column:name"`
	Phone     string    `gorm:"type:varchar(30);column:phone;uniqueIndex"`
	Integral  uint      `gorm:"column:integral"`
	Gender    string    `gorm:"type:varchar(10);column:gender"`
	IsVip     bool      `gorm:"column:is_vip"`
	VipLevels *string   `gorm:"type:varchar(30);column:vip_levels"`
	Status    string    `gorm:"type:varchar(30);column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (u User) TableName() string {
	return common.UserTableName
}
