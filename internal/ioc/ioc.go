package ioc

import (
	"book-management/configs"
	"book-management/internal/repository/do"
	"context"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewGormDB, NewCache)

func NewGormDB(conf configs.AppConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(conf.DB.Addr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&do.User{}, &do.BookInfo{}, &do.BookStock{},
		&do.BookBorrow{}, &do.BookCopy{}, &do.DonateInfo{}, &do.Activity{}, &do.Volunteer{}, &do.VolunteerApplication{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewCache(conf configs.AppConfig) (*redis.Client, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: conf.Cache.Addr,
	})
	err := cli.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return cli, nil
}
