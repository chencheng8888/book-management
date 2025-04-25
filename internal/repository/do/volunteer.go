package do

import (
	"book-management/internal/pkg/common"
	"time"
)

type Volunteer struct {
	ID                    uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Name                  string    `gorm:"type:varchar(255);column:name"`
	Phone                 string    `gorm:"type:varchar(30);column:phone"`
	Age                   int       `gorm:"column:age"`
	ServiceTimePreference string    `gorm:"type:varchar(255);column:service_time_preference"`
	ExpertiseArea         string    `gorm:"type:varchar(255);column:expertise_area"`
	CreatedAt             time.Time `gorm:"column:created_at"`
	UpdatedAt             time.Time `gorm:"column:updated_at"`
}

func (v Volunteer) TableName() string {
	return common.VolunteerTableName
}

type VolunteerApplication struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	Name      string    `gorm:"type:varchar(255);column:name"`
	Phone     string    `gorm:"type:varchar(30);column:phone"`
	Age       int       `gorm:"column:age"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (v VolunteerApplication) TableName() string {
	return common.VolunteerApplicationTableName
}
