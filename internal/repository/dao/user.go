package dao

import (
	"book-management/internal/pkg/common"
	"context"
	"errors"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func (u *UserDao) GetUserName(ctx context.Context, id ...string) (map[string]string, error) {
	if len(id) == 0 {
		return nil, nil
	}
	db := u.db.WithContext(ctx).Table(common.UserTableName)

	type tmp struct {
		ID   string `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}

	var res []tmp
	err := db.Debug().Select("id, name").Where("id IN (?)", id).Find(&res).Error
	if err != nil {
		return nil, err
	}

	var mp = make(map[string]string, len(res))

	for _, v := range res {
		mp[v.ID] = v.Name
	}
	return mp, nil
}

func checkUserIfExist(ctx context.Context, db *gorm.DB, id string) (bool, error) {
	var cnt int64
	err := db.WithContext(ctx).Table(common.UserTableName).Where("id = ?", id).Count(&cnt).Error
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func addUserIntegral(ctx context.Context, db *gorm.DB, id string, integral int) error {
	res := db.WithContext(ctx).Table(common.UserTableName).Where("id = ?", id).Update("integral", gorm.Expr("integral + ?", integral))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("no data updated")
	}
	return nil
}
