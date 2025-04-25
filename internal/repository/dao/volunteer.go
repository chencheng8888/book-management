package dao

import (
	"book-management/internal/repository/do"
	"context"

	"gorm.io/gorm"
)

type VolunteerDao struct {
	db *gorm.DB
}

func NewVolunteerDao(db *gorm.DB) *VolunteerDao {
	return &VolunteerDao{db: db}
}

// Create 创建志愿者
func (d *VolunteerDao) Create(ctx context.Context, volunteer do.Volunteer) error {
	return d.db.WithContext(ctx).Create(&volunteer).Error
}

// Query 查询志愿者信息（支持分页和按ID查询）
func (d *VolunteerDao) Query(ctx context.Context, pageSize, page int, id *uint64) ([]do.Volunteer, error) {
	var volunteers []do.Volunteer
	query := d.db.WithContext(ctx).Model(&do.Volunteer{})

	if id != nil {
		query = query.Where("id = ?", *id)
	}

	err := query.Scopes(paginate(page, pageSize)).Find(&volunteers).Error
	return volunteers, err
}

// Count 统计志愿者数量
func (d *VolunteerDao) Count(ctx context.Context, id *uint64) (int, error) {
	var count int64
	query := d.db.WithContext(ctx).Model(&do.Volunteer{})

	if id != nil {
		query = query.Where("id = ?", *id)
	}

	err := query.Count(&count).Error
	return int(count), err
}

func (d *VolunteerDao) QueryApplications(ctx context.Context, pageSize, page int) ([]do.VolunteerApplication, error) {
	var applications []do.VolunteerApplication
	err := d.db.WithContext(ctx).
		Model(&do.VolunteerApplication{}).
		Scopes(paginate(page, pageSize)).
		Find(&applications).Error
	return applications, err
}

func (d *VolunteerDao) CountApplications(ctx context.Context) (int, error) {
	var count int64
	err := d.db.WithContext(ctx).
		Model(&do.VolunteerApplication{}).
		Count(&count).Error
	return int(count), err
}
