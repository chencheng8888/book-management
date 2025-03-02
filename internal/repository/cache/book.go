package cache

import (
	"book-management/internal/repository/do"
	"context"
	"github.com/redis/go-redis/v9"
)

type BookCache struct {
	rdb *redis.Client
}

func (b *BookCache) DeleteBookStock(ctx context.Context, id uint64) error {
	//TODO implement me
	panic("implement me")
}

func (b *BookCache) DeleteBookInfo(ctx context.Context, id uint64) error {
	//TODO implement me
	panic("implement me")
}

func (b *BookCache) DeleteBookTotalNum(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b *BookCache) GetBookInfoByID(ctx context.Context, ids ...uint64) ([]do.BookInfo, []uint64) {
	//TODO implement me
	panic("implement me")
}

func (b *BookCache) GetBookStockByID(ctx context.Context, ids ...uint64) ([]do.BookStock, []uint64) {
	//TODO implement me
	panic("implement me")
}

func (b *BookCache) GetBookTotalNum(ctx context.Context) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BookCache) SaveBookInfo(ctx context.Context, infos ...do.BookInfo) error {
	//TODO implement me
	panic("implement me")
}

func (b *BookCache) SaveBookStock(ctx context.Context, stocks ...do.BookStock) error {
	//TODO implement me
	panic("implement me")
}

func (b *BookCache) SaveBookTotalNum(ctx context.Context, num int) error {
	//TODO implement me
	panic("implement me")
}
