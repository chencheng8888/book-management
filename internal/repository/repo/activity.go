package repo

import (
	"book-management/internal/controller"
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/tool"
	"book-management/internal/repository/do"
	"context"
)

// ActivityDao 定义数据库操作接口
type ActivityDao interface {
	Create(ctx context.Context, activity do.Activity) error
	Update(ctx context.Context, activity do.Activity) error
	Query(ctx context.Context, pageSize, page int, status *string) ([]do.Activity, error)
	GetByID(ctx context.Context, id uint64) (*do.Activity, error)
	Count(ctx context.Context, status *string) (int, error)
}

type ActivityRepo struct {
	dao ActivityDao
}

func NewActivityRepo(dao ActivityDao) *ActivityRepo {
	return &ActivityRepo{dao: dao}
}

func (r *ActivityRepo) CreateActivity(ctx context.Context, activity controller.Activity) error {
	activityDo,err := convertToDO(activity)
	if err!=nil {
		return errcode.ParamError
	}
	return r.dao.Create(ctx, activityDo)
}

func (r *ActivityRepo) UpdateActivity(ctx context.Context, activity controller.Activity) error {
	activityDo,err := convertToDO(activity)
	if err!=nil {
		return errcode.ParamError
	}
	return r.dao.Update(ctx, activityDo)
}

func (r *ActivityRepo) QueryActivities(
	ctx context.Context,
	pageSize, page int,
	status *string,
	total *int,
) ([]controller.Activity, error) {
	// 参数校验
	if page < 1 || pageSize > 100 {
		return nil, errcode.ParamError
	}

	// 获取总数
	count, err := r.dao.Count(ctx, status)
	if err != nil {
		return nil, err
	}
	*total = count

	// 检查分页
	if maxPage := tool.GetPage(count, pageSize); page > maxPage {
		return nil, errcode.PageError
	}

	// 查询数据
	dos, err := r.dao.Query(ctx, pageSize, page, status)
	if err != nil {
		return nil, err
	}

	return batchToController(dos), nil
}

func (r *ActivityRepo) GetActivityByID(ctx context.Context, id uint64) (*controller.Activity, error) {
	do, err := r.dao.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return convertToController(*do), nil
}

// 数据转换方法
func convertToDO(c controller.Activity) (do.Activity,error) {
	startTime,err := tool.ParseToShanghaiTime(c.Info.StartTime, tool.Format1)
	if err!=nil {
		return do.Activity{},err
	}
	endTime,err := tool.ParseToShanghaiTime(c.Info.EndTime, tool.Format1)
	if err!=nil {
		return do.Activity{},err
	}

	return do.Activity{
		ID: c.ActivityID,
		Name: c.Info.ActivityName,
		Type: c.Info.ActivityType,
		StartTime: startTime,
		EndTime: endTime,
		Manager: c.Info.Manager,
		Phone: c.Info.Phone,
		Addr: c.Info.Addr,
	},nil
}

func convertToController(d do.Activity) *controller.Activity {
	return &controller.Activity{
		ActivityID: d.ID,
		Info: controller.ActivityInfo{
			ActivityName: d.Name,
			ActivityType: d.Type,
			Manager:      d.Manager,
			Phone:        d.Phone,
			Addr:         d.Addr,
			StartTime:    tool.ConvertTimeFormat(d.StartTime,tool.Format1),
			EndTime:      tool.ConvertTimeFormat(d.EndTime,tool.Format1),
		},
	}
}

func batchToController(dos []do.Activity) []controller.Activity {
	res := make([]controller.Activity, 0, len(dos))
	for _, d := range dos {
		res = append(res, *convertToController(d))
	}
	return res
}
