package service

import (
	"book-management/internal/controller"
	"book-management/internal/pkg/errcode"
	"context"
)

type UserRepo interface {
	SearchUser(ctx context.Context, totalNum *int, SearchInfo SearchUserOpts) ([]User, error)
	GetVIPStatics(ctx context.Context) (map[string]int, error)
}

type UserSvc struct {
	userRepo UserRepo
}

func NewUserSvc(userRepo UserRepo) *UserSvc {
	return &UserSvc{
		userRepo: userRepo,
	}
}

func (u *UserSvc) GetVIPStatics(ctx context.Context) (map[string]int, error) {
	return u.userRepo.GetVIPStatics(ctx)
}

func (u *UserSvc) SearchUser(ctx context.Context, req controller.SearchUserReq) ([]controller.User, int, error) {
	var searchUserOptSlice []SearchUserOpt
	searchUserOptSlice = append(searchUserOptSlice, WithPage(req.Page), WithPageSize(req.PageSize))

	if req.UserID != nil {
		searchUserOptSlice = append(searchUserOptSlice, WithUserID(*req.UserID))
	} else {
		var userFieldOptSlice []UserFieldOpt
		if req.UserName != nil {
			userFieldOptSlice = append(userFieldOptSlice, WithUserName(*req.UserName))
		}
		if req.Phone != nil {
			userFieldOptSlice = append(userFieldOptSlice, WithPhone(*req.Phone))
		}
		if req.IsVIP != nil {
			userFieldOptSlice = append(userFieldOptSlice, WithIsVIP(*req.IsVIP))
		}
		if req.Level != nil {
			userFieldOptSlice = append(userFieldOptSlice, WithLevel(*req.Level))
		}
		searchUserOptSlice = append(searchUserOptSlice, WithUserFieldOpts(NewUserFieldOpts(userFieldOptSlice...)))
	}

	var totalNum int
	users, err := u.userRepo.SearchUser(ctx, &totalNum, NewSearchUserOpts(searchUserOptSlice...))
	if err != nil {
		return nil, 0, errcode.SearchDataError
	}
	return batchToControllerUser(users), totalNum, nil
}

type UserFieldOpt func(opts *UserFieldOpts)

type UserFieldOpts struct {
	UserName string
	Phone    string
	IsVIP    bool
	Level    string
}

func NewUserFieldOpts(opts ...UserFieldOpt) UserFieldOpts {
	var userFieldOpts UserFieldOpts
	for _, opt := range opts {
		opt(&userFieldOpts)
	}
	return userFieldOpts
}

func WithUserName(userName string) UserFieldOpt {
	return func(opts *UserFieldOpts) {
		opts.UserName = userName
	}
}
func WithPhone(phone string) UserFieldOpt {
	return func(opts *UserFieldOpts) {
		opts.Phone = phone
	}
}

func WithIsVIP(isvip bool) UserFieldOpt {
	return func(opts *UserFieldOpts) {
		opts.IsVIP = isvip
	}
}

func WithLevel(level string) UserFieldOpt {
	return func(opts *UserFieldOpts) {
		opts.Level = level
	}
}

type SearchUserOpt func(opts *SearchUserOpts)

type SearchUserOpts struct {
	Page     int
	PageSize int
	ByID     bool
	UserID   uint64
	Opts     UserFieldOpts
}

func NewSearchUserOpts(opts ...SearchUserOpt) SearchUserOpts {
	var searchUserOpts SearchUserOpts
	for _, opt := range opts {
		opt(&searchUserOpts)
	}
	if searchUserOpts.Page == 0 {
		searchUserOpts.Page = 1
	}
	if searchUserOpts.PageSize == 0 {
		searchUserOpts.PageSize = 10
	}
	return searchUserOpts
}
func WithPage(page int) SearchUserOpt {
	return func(opts *SearchUserOpts) {
		opts.Page = page
	}
}
func WithPageSize(pageSize int) SearchUserOpt {
	return func(opts *SearchUserOpts) {
		opts.PageSize = pageSize
	}
}
func WithUserID(userID uint64) SearchUserOpt {
	return func(opts *SearchUserOpts) {
		opts.UserID = userID
	}
}
func WithUserFieldOpts(userFieldOpts UserFieldOpts) SearchUserOpt {
	return func(opts *SearchUserOpts) {
		opts.Opts = userFieldOpts
	}
}
