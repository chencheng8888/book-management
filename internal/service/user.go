package service

import (
	"book-management/internal/controller"
	"book-management/internal/pkg/errcode"
	"context"
)

type UserRepo interface {
	SearchUser(ctx context.Context, totalPage *int, SearchInfo SearchUserOpts) ([]User, error)
}

type UserSvc struct {
	userRepo UserRepo
}

func NewUserSvc(userRepo UserRepo) *UserSvc {
	return &UserSvc{
		userRepo: userRepo,
	}
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
		searchUserOptSlice = append(searchUserOptSlice, WithUserFieldOpts(NewUserFieldOpts(userFieldOptSlice...)))
	}

	var totalPage int
	users, err := u.userRepo.SearchUser(ctx, &totalPage, NewSearchUserOpts(searchUserOptSlice...))
	if err != nil {
		return nil, 0, errcode.SearchDataError
	}
	return batchToControllerUser(users), totalPage, nil
}

type UserFieldOpt func(opts *UserFieldOpts)

type UserFieldOpts struct {
	UserName string
	Phone    string
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
		opts.ByID = true
		opts.Opts = userFieldOpts
	}
}
