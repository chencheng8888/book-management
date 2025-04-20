package controller

import (
	"book-management/internal/pkg/errcode"
	treq "book-management/internal/pkg/req"
	"book-management/internal/pkg/resp"
	"book-management/internal/pkg/tool"
	"context"

	"github.com/gin-gonic/gin"
)

type ActivitySvc interface {
	AddActivity(ctx context.Context, req AddActivityReq) (uint64, error)
	UpdateActivity(ctx context.Context, req UpdateActivityReq) error
	QueryActivities(ctx context.Context, req QueryActivityReq, total *int) ([]Activity, error)
}

type ActivityCtrl struct {
	svc ActivitySvc
}

func NewActivityCtrl(svc ActivitySvc) *ActivityCtrl {
	return &ActivityCtrl{
		svc: svc,
	}
}

func (a *ActivityCtrl) RegisterRoute(r *gin.Engine) {
	g := r.Group("/api/v1/activity")
	{
		g.POST("/add", a.AddActivity)
		g.PUT("/update", a.UpdateActivity)
		g.GET("/query", a.QueryActivities)
	}
}

// AddActivity 新增活动
// @Summary 新增活动
// @Description 创建新的图书漂流活动
// @Tags 活动管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object body AddActivityReq true "新增活动请求"
// @Success 200 {object} AddActivityResp
// @Router /api/v1/activity/add [post]
func (a *ActivityCtrl) AddActivity(c *gin.Context) {
	var req AddActivityReq
	if err := treq.ParseRequestBody(c, &req); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	if !a.checkoutActivityType(req.Info.ActivityType) {
		resp.SendResp(c, resp.NewRespFromErr(errcode.ParamError))
		return
	}

	activityID, err := a.svc.AddActivity(c, req)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"activity_id": activityID,
	}))
}

// UpdateActivity 更新活动
// @Summary 更新活动信息
// @Description 更新已存在的图书漂流活动
// @Tags 活动管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object body UpdateActivityReq true "更新活动请求"
// @Success 200 {object} UpdateActivityResp
// @Router /api/v1/activity/update [put]
func (a *ActivityCtrl) UpdateActivity(c *gin.Context) {
	var req UpdateActivityReq

	if err := treq.ParseRequestBody(c, &req); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	if !a.checkoutActivityType(req.Info.ActivityType) {
		resp.SendResp(c, resp.NewRespFromErr(errcode.ParamError))
		return
	}

	if err := a.svc.UpdateActivity(c, req); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	resp.SendResp(c, resp.SuccessResp)
}

// QueryActivities 查询活动
// @Summary 查询活动列表
// @Description 分页查询图书漂流活动
// @Tags 活动管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object query QueryActivityReq true "查询参数"
// @Success 200 {object} QueryActivityResp
// @Router /api/v1/activity/query [get]
func (a *ActivityCtrl) QueryActivities(c *gin.Context) {
	var req QueryActivityReq
	if err := treq.ParseRequestQuery(c, &req); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	var total int
	activities, err := a.svc.QueryActivities(c, req, &total)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"activitys":    activities,
		"total_page":   tool.GetPage(total, req.PageSize),
		"current_page": req.Page,
		"total":        total,
	}))
}

func (a *ActivityCtrl) checkoutActivityType(Type string) bool {
	switch Type {
	case "parent_child_interactions":
		return true
	default:
		return false
	}
}
