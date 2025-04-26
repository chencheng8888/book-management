package repo

import (
	"book-management/internal/pkg/common"
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/tool"
	"book-management/internal/repository/do"
	"book-management/internal/service"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type UserDao interface {
	SearchUserID(ctx context.Context, currentPage, pageSize int, opts ...func(db *gorm.DB)) ([]do.User, error)
	GetUserNum(ctx context.Context, opts ...func(db *gorm.DB)) (int, error)
	GetVIPStatics(ctx context.Context) (map[string]int, error)
}

type UserRepo struct {
	userDao UserDao
}

func NewUserRepo(userDao UserDao) *UserRepo {
	return &UserRepo{userDao: userDao}
}

func (u *UserRepo) SearchUser(ctx context.Context, totalNum *int, SearchInfo service.SearchUserOpts) ([]service.User, error) {
	var opts []func(db *gorm.DB)
	if SearchInfo.ByID {
		opts = append(opts, func(db *gorm.DB) {
			db = db.Where(fmt.Sprintf("%s.id = ?", common.UserTableName), SearchInfo.UserID)
		})
	} else {
		if SearchInfo.Opts.UserName != "" {
			opts = append(opts, func(db *gorm.DB) {
				db = db.Where(fmt.Sprintf("%s.name = ?", common.UserTableName), SearchInfo.Opts.UserName)
			})
		}
		if SearchInfo.Opts.Phone != "" {
			opts = append(opts, func(db *gorm.DB) {
				db = db.Where(fmt.Sprintf("%s.phone = ?", common.UserTableName), SearchInfo.Opts.Phone)
			})
		}

		if SearchInfo.Opts.IsVIP {
			if SearchInfo.Opts.Level != "" {
				opts = append(opts, func(db *gorm.DB) {
					db = db.Where(fmt.Sprintf("%s.is_vip = ? AND %s.vip_levels = ?", common.UserTableName, common.UserTableName), SearchInfo.Opts.IsVIP, SearchInfo.Opts.Level)
				})
			} else {
				opts = append(opts, func(db *gorm.DB) {
					db = db.Where(fmt.Sprintf("%s.is_vip = ?", common.UserTableName), SearchInfo.Opts.IsVIP)
				})
			}
		}
	}

	num, err := u.userDao.GetUserNum(ctx, opts...)
	if err != nil {
		return nil, err
	}

	*totalNum = num

	if SearchInfo.Page > tool.GetPage(num, SearchInfo.PageSize) {
		return nil, errcode.PageError
	}

	users, err := u.userDao.SearchUserID(ctx, SearchInfo.Page, SearchInfo.PageSize, opts...)
	if err != nil {
		return nil, err
	}
	return batchToServiceUser(users), nil
}

func (u *UserRepo) GetVIPStatics(ctx context.Context) (map[string]int, error) {
	return u.userDao.GetVIPStatics(ctx)
}
