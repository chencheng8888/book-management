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
}

type UserRepo struct {
	userDao UserDao
}

func NewUserRepo(userDao UserDao) *UserRepo {
	return &UserRepo{userDao: userDao}
}

func (u *UserRepo) SearchUser(ctx context.Context, totalPage *int, SearchInfo service.SearchUserOpts) ([]service.User, error) {
	if totalPage == nil || SearchInfo.Page <= 0 || SearchInfo.PageSize <= 0 {
		return nil, errcode.PageError
	}
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
	}

	num, err := u.userDao.GetUserNum(ctx, opts...)
	if err != nil {
		return nil, err
	}

	*totalPage = tool.GetPage(num, SearchInfo.PageSize)

	if SearchInfo.Page > *totalPage {
		return nil, errcode.PageError
	}

	users, err := u.userDao.SearchUserID(ctx, SearchInfo.Page, SearchInfo.PageSize, opts...)
	if err != nil {
		return nil, err
	}
	return batchToServiceUser(users), nil
}
