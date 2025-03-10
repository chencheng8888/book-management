package repo

import (
	"book-management/internal/pkg/mocks"
	"book-management/internal/repository/do"
	"book-management/internal/service"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestUpdateBorrowStatus 测试 UpdateBorrowStatus 函数
func TestUpdateBorrowStatus(t *testing.T) {
	// 创建一个mock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建mock对象
	mockBookDao := mocks.NewMockBookDao(ctrl)

	// 初始化BookRepo
	bookRepo := &BookRepo{
		bookDao: mockBookDao,
	}

	// 定义测试用例
	tests := []struct {
		name      string
		bookID    uint64
		copyID    uint64
		status    string
		mockSetup func()
		wantErr   bool
	}{
		{
			name:   "成功更新借阅状态",
			bookID: 1,
			copyID: 1,
			status: "borrowed",
			mockSetup: func() {
				mockBookDao.EXPECT().UpdateBorrowStatus(gomock.Any(), uint64(1), uint64(1), "borrowed").Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "数据库更新失败",
			bookID: 2,
			copyID: 2,
			status: "returned",
			mockSetup: func() {
				mockBookDao.EXPECT().UpdateBorrowStatus(gomock.Any(), uint64(2), uint64(2), "returned").Return(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name:   "无效的状态值",
			bookID: 3,
			copyID: 3,
			status: "invalid_status",
			mockSetup: func() {
				mockBookDao.EXPECT().UpdateBorrowStatus(gomock.Any(), uint64(3), uint64(3), "invalid_status").Return(errors.New("invalid status"))
			},
			wantErr: true,
		},
	}

	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置mock行为
			tt.mockSetup()

			// 调用待测函数
			err := bookRepo.UpdateBorrowStatus(context.Background(), tt.bookID, tt.copyID, tt.status)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBookRepo_QueryBookRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBookDao := mocks.NewMockBookDao(ctrl)
	mockUserDao := mocks.NewMockUserDao(ctrl)

	bookRepo := &BookRepo{
		bookDao: mockBookDao,
		userDao: mockUserDao,
	}

	// 使用固定时间戳
	testTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)

	type args struct {
		ctx         context.Context
		pageSize    int
		currentPage int
		totalPage   *int
		opts        []func(db *gorm.DB)
	}
	tests := []struct {
		name      string
		b         *BookRepo
		args      args
		mockSetup func()
		want      []service.BookBorrowRecord
		wantErr   bool
	}{
		{
			name: "成功查询借阅记录",
			b:    bookRepo,
			args: args{
				ctx:         context.Background(),
				pageSize:    10,
				currentPage: 1,
				totalPage:   new(int),
				opts:        nil,
			},
			mockSetup: func() {
				mockBookDao.EXPECT().GetBookRecordTotalNum(gomock.Any(), gomock.Any()).Return(10, nil)
				mockBookDao.EXPECT().FuzzyQueryBookBorrowRecord(gomock.Any(), 10, 1, gomock.Any()).Return([]do.BookBorrow{
					{BookID: 1, BorrowerID: "user1", CopyID: 1, CreatedTime: testTime, ExpectedReturnTime: testTime.Add(7 * 24 * time.Hour), Status: "borrowed"},
				}, nil)
				mockUserDao.EXPECT().GetUserName(gomock.Any(), "user1").Return(map[string]string{"user1": "User One"}, nil)
			},
			want: []service.BookBorrowRecord{
				{BookID: 1, BorrowerID: "user1", Borrower: "User One", CopyID: 1, BorrowTime: testTime, ExpectedTime: testTime.Add(7 * 24 * time.Hour), ReturnStatus: "borrowed"},
			},
			wantErr: false,
		},
		{
			name: "查询借阅记录失败",
			b:    bookRepo,
			args: args{
				ctx:         context.Background(),
				pageSize:    10,
				currentPage: 1,
				totalPage:   new(int),
				opts:        nil,
			},
			mockSetup: func() {
				mockBookDao.EXPECT().GetBookRecordTotalNum(gomock.Any(), gomock.Any()).Return(0, errors.New("database error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "获取用户名失败",
			b:    bookRepo,
			args: args{
				ctx:         context.Background(),
				pageSize:    10,
				currentPage: 1,
				totalPage:   new(int),
				opts:        nil,
			},
			mockSetup: func() {
				mockBookDao.EXPECT().GetBookRecordTotalNum(gomock.Any(), gomock.Any()).Return(10, nil)
				mockBookDao.EXPECT().FuzzyQueryBookBorrowRecord(gomock.Any(), 10, 1, gomock.Any()).Return([]do.BookBorrow{
					{BookID: 1, BorrowerID: "user1", CopyID: 1, CreatedTime: testTime, ExpectedReturnTime: testTime.Add(7 * 24 * time.Hour), Status: "borrowed"},
				}, nil)
				mockUserDao.EXPECT().GetUserName(gomock.Any(), "user1").Return(nil, errors.New("user not found"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "pageSize 为负数",
			b:    bookRepo,
			args: args{
				ctx:         context.Background(),
				pageSize:    -1,
				currentPage: 1,
				totalPage:   new(int),
				opts:        nil,
			},
			mockSetup: func() {},
			want:      nil,
			wantErr:   true,
		},
		{
			name: "currentPage 为负数",
			b:    bookRepo,
			args: args{
				ctx:         context.Background(),
				pageSize:    10,
				currentPage: -1,
				totalPage:   new(int),
				opts:        nil,
			},
			mockSetup: func() {},
			want:      nil,
			wantErr:   true,
		},
		{
			name: "totalPage 为 nil",
			b:    bookRepo,
			args: args{
				ctx:         context.Background(),
				pageSize:    10,
				currentPage: 1,
				totalPage:   nil,
				opts:        nil,
			},
			mockSetup: func() {
				mockBookDao.EXPECT().GetBookRecordTotalNum(gomock.Any(), gomock.Any()).Return(10, nil)
				mockBookDao.EXPECT().FuzzyQueryBookBorrowRecord(gomock.Any(), 10, 1, gomock.Any()).Return([]do.BookBorrow{
					{BookID: 1, BorrowerID: "user1", CopyID: 1, CreatedTime: testTime, ExpectedReturnTime: testTime.Add(7 * 24 * time.Hour), Status: "borrowed"},
				}, nil)
				mockUserDao.EXPECT().GetUserName(gomock.Any(), "user1").Return(map[string]string{"user1": "User One"}, nil)
			},
			want: []service.BookBorrowRecord{
				{BookID: 1, BorrowerID: "user1", Borrower: "User One", CopyID: 1, BorrowTime: testTime, ExpectedTime: testTime.Add(7 * 24 * time.Hour), ReturnStatus: "borrowed"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			got, err := tt.b.QueryBookRecord(tt.args.ctx, tt.args.pageSize, tt.args.currentPage, tt.args.totalPage, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("BookRepo.QueryBookRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, "got and want should be equal")
		})
	}
}

func TestAddBookBorrowRecord(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建mock对象
	mockBookDao := mocks.NewMockBookDao(ctrl)

	// 初始化BookRepo
	bookRepo := &BookRepo{
		bookDao: mockBookDao,
	}

	// 定义测试用例
	tests := []struct {
		name               string
		bookID             uint64
		borrowerID         string
		expectedReturnTime time.Time
		copyID             *uint64
		mockSetup          func()
		wantErr            bool
	}{
		{
			name:               "成功添加借阅记录",
			bookID:             1,
			borrowerID:         "user123",
			expectedReturnTime: time.Now().AddDate(0, 0, 14),
			copyID:             func() *uint64 { id := uint64(1); return &id }(),
			mockSetup: func() {
				mockBookDao.EXPECT().AddBookBorrowRecord(gomock.Any(), uint64(1), "user123", gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name:               "数据库添加失败",
			bookID:             2,
			borrowerID:         "user456",
			expectedReturnTime: time.Now().AddDate(0, 0, 14),
			copyID:             func() *uint64 { id := uint64(2); return &id }(),
			mockSetup: func() {
				expectedError := errors.New("database error")
				mockBookDao.EXPECT().AddBookBorrowRecord(gomock.Any(), uint64(2), "user456", gomock.Any(), gomock.Any()).Return(expectedError)
			},
			wantErr: true,
		},
	}

	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置mock行为
			tt.mockSetup()

			// 调用待测函数
			err := bookRepo.AddBookBorrowRecord(context.Background(), tt.bookID, tt.borrowerID, tt.expectedReturnTime, tt.copyID)

			// 验证结果
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckBookInfoIfExist(t *testing.T) {
	// 创建一个gomock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建MockBookDao实例
	mockBookDao := mocks.NewMockBookDao(ctrl)

	// 创建BookRepo实例，注入MockBookDao
	bookRepo := &BookRepo{
		bookDao: mockBookDao,
	}

	// 测试用例1: 书籍存在
	ctx := context.Background()
	mockBookDao.EXPECT().CheckBookIfExist(ctx, "Book1", "Author1", "Publisher1", "Category1").Return(uint64(1), true)
	id, exists := bookRepo.CheckBookInfoIfExist(ctx, "Book1", "Author1", "Publisher1", "Category1")
	assert.Equal(t, uint64(1), id)
	assert.True(t, exists)

	// 测试用例2: 书籍不存在
	mockBookDao.EXPECT().CheckBookIfExist(ctx, "Book2", "Author2", "Publisher2", "Category2").Return(uint64(0), false)
	id, exists = bookRepo.CheckBookInfoIfExist(ctx, "Book2", "Author2", "Publisher2", "Category2")
	assert.Equal(t, uint64(0), id)
	assert.False(t, exists)
}

func TestAddBookStock_TableDriven(t *testing.T) {
	// 创建测试用例表
	testCases := []struct {
		name        string
		setupMocks  func(*mocks.MockBookCache, *mocks.MockBookDao)
		wantErr     bool
		errContains string
	}{
		{
			name: "Cache delete failure",
			setupMocks: func(cache *mocks.MockBookCache, dao *mocks.MockBookDao) {
				cache.EXPECT().DeleteBookStock(gomock.Any(), uint64(1)).
					Return(errors.New("cache error")).Times(1)
			},
			wantErr:     true,
			errContains: "cache error",
		},
		{
			name: "DAO operation failure",
			setupMocks: func(cache *mocks.MockBookCache, dao *mocks.MockBookDao) {
				cache.EXPECT().DeleteBookStock(gomock.Any(), uint64(1)).Return(nil)
				dao.EXPECT().AddBookStock(gomock.Any(), uint64(1), uint(10), (*string)(nil)).
					Return(errors.New("dao error"))
			},
			wantErr:     true,
			errContains: "dao error",
		},
		{
			name: "Full success flow",
			setupMocks: func(cache *mocks.MockBookCache, dao *mocks.MockBookDao) {
				// 同步删除调用1次
				cache.EXPECT().DeleteBookStock(gomock.Any(), uint64(1)).Return(nil)
				// DAO操作成功
				dao.EXPECT().AddBookStock(gomock.Any(), uint64(1), uint(10), (*string)(nil)).Return(nil)
				// 异步删除会再次调用
				cache.EXPECT().DeleteBookStock(gomock.Any(), uint64(1)).Return(nil)
			},
			wantErr: false,
		},
	}

	//// 使用gomonkey加速定时器（全局设置）
	//patch := gomonkey.ApplyFunc(time.AfterFunc, func(d time.Duration, f func()) *time.Timer {
	//	f() // 立即执行回调
	//	return nil
	//})
	//defer patch.Reset()

	// 遍历执行所有测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCache := mocks.NewMockBookCache(ctrl)
			mockDao := mocks.NewMockBookDao(ctrl)

			// 配置mock行为
			tc.setupMocks(mockCache, mockDao)

			repo := &BookRepo{bookDao: mockDao, bookCache: mockCache}
			err := repo.AddBookStock(context.Background(), 1, 10, nil)

			// 错误断言
			if tc.wantErr {
				if err == nil || err.Error() != tc.errContains {
					t.Errorf("Expected error contains %q, got %v", tc.errContains, err)
				}
				return
			}

			// 成功断言
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			time.Sleep(time.Millisecond * 2000)
		})
	}
}

func TestBookRepo_getBookInID(t *testing.T) {

	type args struct {
		ids []uint64
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(*mocks.MockBookDao, *mocks.MockBookCache)
		want     []do.BookInfo
		want1    []do.BookStock
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "缓存全部命中",
			args: args{
				ids: []uint64{1, 2, 3},
			},
			mockFunc: func(dao *mocks.MockBookDao, cache *mocks.MockBookCache) {
				cache.EXPECT().GetBookInfoByID(gomock.Any(), gomock.Any()).Return([]do.BookInfo{
					{
						ID: 1,
					},
					{
						ID: 2,
					},
					{
						ID: 3,
					},
				}, nil)
				cache.EXPECT().GetBookStockByID(gomock.Any(), gomock.Any()).Return([]do.BookStock{
					{
						BookID: 1,
					},
					{
						BookID: 2,
					},
					{
						BookID: 3,
					},
				}, nil)
			},
			want: []do.BookInfo{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
				{
					ID: 3,
				},
			},
			want1: []do.BookStock{
				{
					BookID: 1,
				},
				{
					BookID: 2,
				},
				{
					BookID: 3,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "缓存没有完全命中",
			args: args{
				ids: []uint64{1, 2, 3},
			},
			mockFunc: func(dao *mocks.MockBookDao, cache *mocks.MockBookCache) {
				cache.EXPECT().GetBookInfoByID(gomock.Any(), gomock.Any()).Return([]do.BookInfo{
					{
						ID: 1,
					},
					{
						ID: 2,
					},
				}, []uint64{3})
				dao.EXPECT().GetBookInfoByID(gomock.Any(), gomock.Any()).Return([]do.BookInfo{
					{
						ID: 3,
					},
				}, nil)
				cache.EXPECT().GetBookStockByID(gomock.Any(), gomock.Any()).Return([]do.BookStock{
					{
						BookID: 1,
					},
					{
						BookID: 2,
					},
					{
						BookID: 3,
					},
				}, nil)
			},
			want: []do.BookInfo{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
				{
					ID: 3,
				},
			},
			want1: []do.BookStock{
				{
					BookID: 1,
				},
				{
					BookID: 2,
				},
				{
					BookID: 3,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "缓存miss，同时数据库查询失败",
			args: args{
				ids: []uint64{1, 2, 3},
			},
			mockFunc: func(dao *mocks.MockBookDao, cache *mocks.MockBookCache) {
				cache.EXPECT().GetBookInfoByID(gomock.Any(), gomock.Any()).Return([]do.BookInfo{
					{
						ID: 1,
					},
					{
						ID: 2,
					},
				}, []uint64{3})
				dao.EXPECT().GetBookInfoByID(gomock.Any(), gomock.Any()).Return([]do.BookInfo{
					{
						ID: 3,
					},
				}, errors.New("some err"))
			},
			want:    nil,
			want1:   nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			bookDao := mocks.NewMockBookDao(ctrl)
			cacheDao := mocks.NewMockBookCache(ctrl)

			tt.mockFunc(bookDao, cacheDao)
			bookRepo := &BookRepo{
				bookDao:   bookDao,
				bookCache: cacheDao,
			}
			got, got1, err := bookRepo.getBookInID(context.Background(), tt.args.ids...)
			tt.wantErr(t, err)
			assert.ElementsMatch(t, got, tt.want)
			assert.ElementsMatch(t, got1, tt.want1)
		})
	}
}

func Test_batchToServiceBookRecord(t *testing.T) {
	fixedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)

	type args struct {
		borrows []do.BookBorrow
		user    map[string]string
	}
	tests := []struct {
		name string
		args args
		want []service.BookBorrowRecord
	}{
		{
			name: "正常转换",
			args: args{
				borrows: []do.BookBorrow{
					{
						BookID:             1,
						BorrowerID:         "user1",
						CopyID:             1001,
						CreatedTime:        fixedTime,
						ExpectedReturnTime: fixedTime.Add(7 * 24 * time.Hour),
						Status:             "borrowed",
					},
					{
						BookID:             2,
						BorrowerID:         "user2",
						CopyID:             2001,
						CreatedTime:        fixedTime,
						ExpectedReturnTime: fixedTime.Add(14 * 24 * time.Hour),
						Status:             "returned",
					},
				},
				user: map[string]string{
					"user1": "张三",
					"user2": "李四",
				},
			},
			want: []service.BookBorrowRecord{
				{
					BookID:       1,
					BorrowerID:   "user1",
					Borrower:     "张三",
					CopyID:       1001,
					BorrowTime:   fixedTime,
					ExpectedTime: fixedTime.Add(7 * 24 * time.Hour),
					ReturnStatus: "borrowed",
				},
				{
					BookID:       2,
					BorrowerID:   "user2",
					Borrower:     "李四",
					CopyID:       2001,
					BorrowTime:   fixedTime,
					ExpectedTime: fixedTime.Add(14 * 24 * time.Hour),
					ReturnStatus: "returned",
				},
			},
		},
		{
			name: "部分用户不存在",
			args: args{
				borrows: []do.BookBorrow{
					{
						BookID:             3,
						BorrowerID:         "user3",
						CopyID:             3001,
						CreatedTime:        fixedTime,
						ExpectedReturnTime: fixedTime.Add(7 * 24 * time.Hour),
						Status:             "overdue",
					},
				},
				user: map[string]string{
					"user1": "张三", // 故意不包含user3
				},
			},
			want: []service.BookBorrowRecord{
				{
					BookID:       3,
					BorrowerID:   "user3",
					Borrower:     "", // 预期用户名为空
					CopyID:       3001,
					BorrowTime:   fixedTime,
					ExpectedTime: fixedTime.Add(7 * 24 * time.Hour),
					ReturnStatus: "overdue",
				},
			},
		},
		{
			name: "空记录",
			args: args{
				borrows: []do.BookBorrow{},
				user:    map[string]string{},
			},
			want: []service.BookBorrowRecord{},
		},
		{
			name: "用户名为空",
			args: args{
				borrows: []do.BookBorrow{
					{
						BookID:             4,
						BorrowerID:         "user4",
						CopyID:             4001,
						CreatedTime:        fixedTime,
						ExpectedReturnTime: fixedTime.Add(7 * 24 * time.Hour),
						Status:             "processing",
					},
				},
				user: nil, // 测试空map情况
			},
			want: []service.BookBorrowRecord{
				{
					BookID:       4,
					BorrowerID:   "user4",
					Borrower:     "",
					CopyID:       4001,
					BorrowTime:   fixedTime,
					ExpectedTime: fixedTime.Add(7 * 24 * time.Hour),
					ReturnStatus: "processing",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, batchToServiceBookRecord(tt.args.borrows, tt.args.user),
				"batchToServiceBookRecord(%v, %v)", tt.args.borrows, tt.args.user)
		})
	}
}
