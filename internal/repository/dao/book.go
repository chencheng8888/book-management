package dao

import (
	"book-management/internal/repository/do"
	"context"
	"gorm.io/gorm"
)

type BookDao struct {
	db *gorm.DB
}

func (b *BookDao) CheckBookIfExist(ctx context.Context, name, author, publisher, category string) (uint64, bool) {
	//TODO implement me
	panic("implement me")
}

func (b *BookDao) AddBookStock(ctx context.Context, id uint64, num uint, where *string) error {
	//TODO implement me
	panic("implement me")
}

func (b *BookDao) RegisterAndAddBookStock(ctx context.Context, bookInfo do.BookInfo, addedNum uint, where *string) error {
	//TODO implement me
	panic("implement me")
}

func (b *BookDao) GetBookInfoByID(ctx context.Context, id ...uint64) ([]do.BookInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookDao) GetBookStockByID(ctx context.Context, id ...uint64) ([]do.BookStock, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookDao) FuzzyQueryBookID(ctx context.Context, pageSize int, page int, opts ...func(db *gorm.DB)) ([]uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookDao) GetBookTotalNum(ctx context.Context) (int, error) {
	//TODO implement me
	panic("implement me")
}
