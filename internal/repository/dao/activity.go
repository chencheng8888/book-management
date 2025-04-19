package dao

import (
	"book-management/internal/pkg/common"
	"book-management/internal/repository/do"
	"context"
	"time"

	"gorm.io/gorm"
)

type ActivityDao struct {
	db *gorm.DB
}

func NewActivityDao(db *gorm.DB) *ActivityDao {
	return &ActivityDao{db: db}
}

func (d *ActivityDao) Create(ctx context.Context, activity do.Activity) error {
	return d.db.WithContext(ctx).Create(&activity).Error
}

func (d *ActivityDao) Update(ctx context.Context, activity do.Activity) error {
	return d.db.WithContext(ctx).
		Model(&do.Activity{}).
		Where("id = ?", activity.ID).
		Updates(activity).Error
}

func (d *ActivityDao) Query(ctx context.Context, pageSize, page int, status *string) ([]do.Activity, error) {
	var activities []do.Activity

	query := d.db.WithContext(ctx).Model(&do.Activity{})

	if status != nil {
		query = d.applyStatusFilter(query, *status)
	}

	err := query.Scopes(paginate(page, pageSize)).
		Order("start_time DESC").
		Find(&activities).Error

	return activities, err
}

func (d *ActivityDao) GetByID(ctx context.Context, id uint64) (*do.Activity, error) {
	var activity do.Activity
	err := d.db.WithContext(ctx).
		Where("id = ?", id).
		First(&activity).Error
	return &activity, err
}

func (d *ActivityDao) Count(ctx context.Context, status *string) (int, error) {
	var count int64

	query := d.db.WithContext(ctx).Model(&do.Activity{})

	if status != nil {
		query = d.applyStatusFilter(query, *status)
	}

	err := query.Count(&count).Error
	return int(count), err
}

// 私有方法用于构建状态过滤条件
func (d *ActivityDao) applyStatusFilter(query *gorm.DB, status string) *gorm.DB {
	now := time.Now()

	switch status {
	case common.ActivityStatusPending:
		return query.Where("start_time > ?",now)
	case common.ActivityStatusOngoing:
		return query.Where("start_time <= ? AND end_time >= ?",now,now)
	case common.ActivityStatusEnded:
		return query.Where("end_time < ?",now)
	default:
		return query
	}
}

// 分页scope
func paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}