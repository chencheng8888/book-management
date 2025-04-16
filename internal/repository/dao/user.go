package dao

import (
	"book-management/internal/pkg/common"
	"book-management/internal/repository/do"
	"context"
	"errors"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (u *UserDao) GetUserPhone(ctx context.Context, id ...uint64) (map[uint64]string, error) {

	db := u.db.WithContext(ctx).Table(common.UserTableName)

	var res []struct {
		ID    uint64 `gorm:"column:id"`
		Phone string `gorm:"column:phone"`
	}
	err := db.Debug().Select("id, phone").Where("id in (?)", id).Find(&res).Error
	if err != nil {
		return nil, err
	}
	var mp = make(map[uint64]string, len(res))
	for _, v := range res {
		mp[v.ID] = v.Phone
	}
	return mp, nil
}

func (u *UserDao) SearchUserID(ctx context.Context, currentPage, pageSize int, opts ...func(db *gorm.DB)) ([]do.User, error) {
	if currentPage <= 0 {
		currentPage = 1
	}
	if pageSize <= 0 {
		pageSize = 1
	}

	var users []do.User
	db := u.db.WithContext(ctx).Table(common.UserTableName)
	for _, opt := range opts {
		opt(db)
	}
	err := db.Debug().Where("id >= (?)", db.Select("id").Order("id").Offset((currentPage-1)*pageSize).Limit(1)).
		Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserDao) GetUserNum(ctx context.Context, opts ...func(db *gorm.DB)) (int, error) {
	db := u.db.WithContext(ctx)

	for _, opt := range opts {
		opt(db)
	}
	var num int64
	err := db.Debug().Count(&num).Error
	if err != nil {
		return 0, err
	}
	return int(num), nil
}

func (u *UserDao) GetUserName(ctx context.Context, id ...uint64) (map[uint64]string, error) {
	if len(id) == 0 {
		return nil, nil
	}
	db := u.db.WithContext(ctx).Table(common.UserTableName)

	type tmp struct {
		ID   uint64 `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}

	var res []tmp
	err := db.Debug().Select("id, name").Where("id IN (?)", id).Find(&res).Error
	if err != nil {
		return nil, err
	}

	var mp = make(map[uint64]string, len(res))

	for _, v := range res {
		mp[v.ID] = v.Name
	}
	return mp, nil
}

func checkUserIfExist(ctx context.Context, db *gorm.DB, id uint64) (bool, error) {
	var cnt int64
	err := db.WithContext(ctx).Table(common.UserTableName).Where("id = ?", id).Count(&cnt).Error
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func addUserIntegral(ctx context.Context, db *gorm.DB, id uint64, integral int) error {
	res := db.WithContext(ctx).Table(common.UserTableName).Where("id = ?", id).Update("integral", gorm.Expr("integral + ?", integral))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("no data updated")
	}
	return nil
}
