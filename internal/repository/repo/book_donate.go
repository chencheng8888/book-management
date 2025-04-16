package repo

import (
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/tool"
	"book-management/internal/repository/do"
	"book-management/internal/service"
	"context"
	"time"
)

type BookDonateDao interface {
	GetBookDonateRecordsNum(ctx context.Context) (int, error)
	GetBookDonateInfos(ctx context.Context, pageSize int, page int) ([]do.DonateInfo, error)
	GetDonateRanking(ctx context.Context, top int) ([]DonateRanking, error)
}

type GetUserPhone interface {
	GetUserPhone(ctx context.Context, id ...uint64) (map[uint64]string, error)
}

type BookDonateRepo struct {
	donateDao    BookDonateDao
	infoDao      BookInfoDao
	getUserNamer GetUserNamer
	getUserPhone GetUserPhone
}

func NewBookDonateRepo(donateDao BookDonateDao, infoDao BookInfoDao, getUserNamer GetUserNamer, getUserPhone GetUserPhone) *BookDonateRepo {
	return &BookDonateRepo{
		donateDao:    donateDao,
		infoDao:      infoDao,
		getUserNamer: getUserNamer,
		getUserPhone: getUserPhone,
	}
}

func (b *BookDonateRepo) ListBookDonateRecordsReq(ctx context.Context, pageSize int, currentPage int, totalNum *int) ([]service.BookDonateRecord, error) {
	//获取全部页数
	total, err := b.donateDao.GetBookDonateRecordsNum(ctx)
	if err != nil {
		return nil, err
	}

	*totalNum = total
	if maxPage := tool.GetPage(total, pageSize); currentPage > maxPage {
		return nil, errcode.PageError
	}

	donateInfos, err := b.donateDao.GetBookDonateInfos(ctx, pageSize, currentPage)
	if err != nil {
		return nil, err
	}

	var userIDs = make([]uint64, 0, len(donateInfos))

	var bookIDs = make([]uint64, 0, len(donateInfos))

	for _, info := range donateInfos {
		userIDs = append(userIDs, info.UserID)
		bookIDs = append(bookIDs, info.BookID)
	}

	//去重
	userIDs = tool.Unique(userIDs)
	bookIDs = tool.Unique(bookIDs)

	userNameMap, err := b.getUserNamer.GetUserName(ctx, userIDs...)
	if err != nil {
		return nil, err
	}
	phoneMap, err := b.getUserPhone.GetUserPhone(ctx, userIDs...)
	if err != nil {
		return nil, err
	}
	infos, err := b.infoDao.GetBookInfoByID(ctx, bookIDs...)
	if err != nil {
		return nil, err
	}

	var bookNameMap = make(map[uint64]string)
	for _, info := range infos {
		bookNameMap[info.ID] = info.Name
	}
	return batchToServiceBookDonateRecord(donateInfos, userNameMap, phoneMap, bookNameMap), nil
}
func (b *BookDonateRepo) GetBookDonateRanking(ctx context.Context, top int) ([]service.Rank, error) {
	rankings, err := b.donateDao.GetDonateRanking(ctx, top)
	if err != nil {
		return nil, err
	}
	return convertToServiceRank(rankings), nil
}

type DonateRanking struct {
	UserID      uint64    `gorm:"user_id"`
	UserName    string    `gorm:"user_name"`
	DonateNum   int       `gorm:"donate_num"`
	DonateTimes int       `gorm:"donate_times"`
	UpdatedAt   time.Time `gorm:"updated_at"`
}

func convertToServiceRank(rankings []DonateRanking) []service.Rank {
	res := make([]service.Rank, 0, len(rankings))
	for _, rank := range rankings {
		var r service.Rank
		r.UserID = rank.UserID
		r.UserName = rank.UserName
		r.DonateNum = rank.DonateNum
		r.DonateTimes = rank.DonateTimes
		r.UpdatedAt = tool.ConvertTimeFormat(rank.UpdatedAt)
		res = append(res, r)
	}
	return res
}
