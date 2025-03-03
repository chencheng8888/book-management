package cache

import (
	"book-management/internal/repository/do"
	"book-management/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type BookCache struct {
	rdb                *redis.Client
	bookInfoExpire     time.Duration
	bookStockExpire    time.Duration
	bookTotalNumExpire time.Duration
}

func (b *BookCache) DeleteBookStock(ctx context.Context, id uint64) error {
	key := b.generateBookStockKey(id)
	err := b.rdb.Del(ctx, key).Err()
	if err != nil {
		logger.LogPrinter.Warnf("delete book stock in cache failed, id: %d, err: %v", id, err)
	}
	return err
}

func (b *BookCache) DeleteBookInfo(ctx context.Context, id uint64) error {
	key := b.generateBookInfoKey(id)
	err := b.rdb.Del(ctx, key).Err()
	if err != nil {
		logger.LogPrinter.Warnf("delete book info in cache failed, id: %d, err: %v", id, err)
	}
	return err
}

func (b *BookCache) DeleteBookTotalNum(ctx context.Context) error {
	key := b.generateBookTotalNumKey()
	err := b.rdb.Del(ctx, key).Err()
	if err != nil {
		logger.LogPrinter.Warnf("delete book total num in cache failed, err: %v", err)
	}
	return err
}

func (b *BookCache) GetBookInfoByID(ctx context.Context, ids ...uint64) ([]do.BookInfo, []uint64) {
	if len(ids) == 0 {
		return nil, nil
	}

	var keys = make([]string, 0, len(ids))

	for _, id := range ids {
		keys = append(keys, b.generateBookInfoKey(id))
	}
	vals, err := b.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, ids
	}

	var infos []do.BookInfo
	var leftIDs []uint64

	for i, val := range vals {
		if val == nil {
			leftIDs = append(leftIDs, ids[i])
			continue
		}

		var info do.BookInfo
		err = json.Unmarshal([]byte(val.(string)), &info)
		if err != nil {
			leftIDs = append(leftIDs, ids[i])
			continue
		}
		infos = append(infos, info)
	}
	return infos, leftIDs
}

func (b *BookCache) GetBookStockByID(ctx context.Context, ids ...uint64) ([]do.BookStock, []uint64) {
	if len(ids) == 0 {
		return nil, nil
	}

	var keys = make([]string, 0, len(ids))

	for _, id := range ids {
		keys = append(keys, b.generateBookStockKey(id))
	}
	vals, err := b.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, ids
	}

	var stocks []do.BookStock
	var leftIDs []uint64

	for i, val := range vals {
		if val == nil {
			leftIDs = append(leftIDs, ids[i])
			continue
		}

		var stock do.BookStock
		err = json.Unmarshal([]byte(val.(string)), &stock)
		if err != nil {
			leftIDs = append(leftIDs, ids[i])
			continue
		}
		stocks = append(stocks, stock)
	}
	return stocks, leftIDs
}

func (b *BookCache) GetBookTotalNum(ctx context.Context) (int, error) {
	key := b.generateBookTotalNumKey()
	val, err := b.rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (b *BookCache) SaveBookInfo(ctx context.Context, infos ...do.BookInfo) error {
	if len(infos) == 0 {
		return nil
	}
	var marshalData = make(map[string]string, len(infos))

	for _, info := range infos {
		key := b.generateBookInfoKey(info.ID)
		val, err := json.Marshal(&info)
		if err != nil {
			logger.LogPrinter.Warnf("marshal book info failed, id: %d, err: %v", info.ID, err)
			continue
		}
		marshalData[key] = string(val)
	}
	var pipe = b.rdb.Pipeline()

	for k, v := range marshalData {
		pipe.Set(ctx, k, v, b.bookInfoExpire)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.LogPrinter.Warnf("save book info in cache failed, err: %v", err)
	}
	return err
}

func (b *BookCache) SaveBookStock(ctx context.Context, stocks ...do.BookStock) error {
	if len(stocks) == 0 {
		return nil
	}
	var marshalData = make(map[string]string, len(stocks))

	for _, stock := range stocks {
		key := b.generateBookStockKey(stock.BookID)
		val, err := json.Marshal(&stock)
		if err != nil {
			logger.LogPrinter.Warnf("marshal book stock failed, id: %d, err: %v", stock.BookID, err)
			continue
		}
		marshalData[key] = string(val)
	}
	var pipe = b.rdb.Pipeline()

	for k, v := range marshalData {
		pipe.Set(ctx, k, v, b.bookStockExpire)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.LogPrinter.Warnf("save book stock in cache failed, err: %v", err)
	}
	return err
}

func (b *BookCache) SaveBookTotalNum(ctx context.Context, num int) error {
	key := b.generateBookTotalNumKey()
	err := b.rdb.Set(ctx, key, num, b.bookTotalNumExpire).Err()
	if err != nil {
		logger.LogPrinter.Warnf("save book total num in cache failed, err: %v", err)
	}
	return err
}

func (b *BookCache) generateBookInfoKey(id uint64) string {
	return fmt.Sprintf("book-info:%d", id)
}
func (b *BookCache) generateBookStockKey(id uint64) string {
	return fmt.Sprintf("book-stock:%d", id)
}
func (b *BookCache) generateBookTotalNumKey() string {
	return "book-total-num"
}
