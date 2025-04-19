package service

import (
	"book-management/internal/controller"
	"book-management/internal/pkg/errcode"
	"book-management/pkg/logger"
	"context"
)

type ActivityRepo interface {
	CreateActivity(ctx context.Context, activity controller.Activity) error
	UpdateActivity(ctx context.Context, activity controller.Activity) error
	QueryActivities(ctx context.Context, pageSize, page int, status *string, total *int) ([]controller.Activity, error)
	GetActivityByID(ctx context.Context, id uint64) (*controller.Activity, error)
}

type ActivitySvc struct {
	repo ActivityRepo
	ider *MyIDer
}

func NewActivitySvc(repo ActivityRepo) *ActivitySvc {
	return &ActivitySvc{
		repo: repo,
		ider: NewMyIDer(),
	}
}

// AddActivity 创建新活动
func (s *ActivitySvc) AddActivity(ctx context.Context, req controller.AddActivityReq) (uint64, error) {
	activityID, err := s.ider.GenerateID(ctx)
	if err != nil {
		logger.LogPrinter.Errorf("生成活动ID失败: %v", err)
		return 0, errcode.GenerateIDError
	}

	activity := controller.Activity{
		ActivityID: activityID,
		Info:       req.Info,
	}

	if err := s.repo.CreateActivity(ctx, activity); err != nil {
		logger.LogPrinter.Errorf("创建活动失败[ID:%d]: %v", activityID, err)
		return 0, errcode.AddActivityError
	}

	return activityID, nil
}

// UpdateActivity 更新活动信息
func (s *ActivitySvc) UpdateActivity(ctx context.Context, req controller.UpdateActivityReq) error {
	_, err := s.repo.GetActivityByID(ctx, req.ActivityID)
	if err != nil {
		logger.LogPrinter.Errorf("查询活动失败[ID:%d]: %v", req.ActivityID, err)
		return errcode.SearchDataError
	}

	updated := controller.Activity{
		ActivityID: req.ActivityID,
		Info:       req.Info,
	}

	if err := s.repo.UpdateActivity(ctx, updated); err != nil {
		logger.LogPrinter.Errorf("更新活动失败[ID:%d]: %v", req.ActivityID, err)
		return errcode.UpdateDataError
	}

	return nil
}

// QueryActivities 查询活动列表
func (s *ActivitySvc) QueryActivities(ctx context.Context, req controller.QueryActivityReq, total *int) ([]controller.Activity, error) {
	var statusFilter *string
	if req.Status != nil {
		if !isValidStatus(*req.Status) {
			return nil, errcode.ParamError
		}
		statusFilter = req.Status
	}

	activities, err := s.repo.QueryActivities(
		ctx,
		req.PageSize,
		req.Page,
		statusFilter,
		total,
	)
	if err != nil {
		logger.LogPrinter.Errorf("查询活动列表失败: %v", err)
		return nil, errcode.SearchDataError
	}

	return activities, nil
}

// 辅助函数校验状态合法性
func isValidStatus(status string) bool {
	switch status {
	case "pending", "ongoing", "ended":
		return true
	default:
		return false
	}
}
