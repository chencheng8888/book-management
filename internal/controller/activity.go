package controller

import (
	"book-management/internal/pkg/common"
	"book-management/internal/pkg/errcode"
	treq "book-management/internal/pkg/req"
	"book-management/internal/pkg/resp"
	"book-management/internal/pkg/tool"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ActivitySvc interface {
	AddActivity(ctx context.Context, req AddActivityReq) (uint64, error)
	UpdateActivity(ctx context.Context, req UpdateActivityReq) error
	QueryActivities(ctx context.Context, req QueryActivityReq, total *int) ([]Activity, error)
}

type ActivityCtrl struct {
	svc ActivitySvc
	db  *gorm.DB
}

func NewActivityCtrl(svc ActivitySvc, db *gorm.DB) *ActivityCtrl {
	return &ActivityCtrl{
		svc: svc,
		db:  db,
	}
}

func (a *ActivityCtrl) RegisterRoute(r *gin.Engine) {
	g := r.Group("/api/v1/activity")
	{
		g.POST("/add", a.AddActivity)
		g.PUT("/update", a.UpdateActivity)
		g.GET("/query", a.QueryActivities)
		g.GET("/get_statics", a.GetActivityStatics)
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

// GetActivityStatics 获取活动统计信息
// @Summary 获取活动统计信息
// @Description 获取活动的总数、报名人数、参与率、已结束、进行中和即将开始的活动数量
// @Tags 活动管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Success 200 {object} GetActivityStaticsResp
// @Router /api/v1/activity/get_statics [get]
func (a *ActivityCtrl) GetActivityStatics(c *gin.Context) {
	var totalNum int64
	a.db.WithContext(c).Table(common.ActivityTableName).Count(&totalNum)

	var totalApplicants int64
	a.db.WithContext(c).Table(common.ActivityTableName).Select("sum(sign_up_num)").Scan(&totalApplicants)

	var totalParticipate int64
	a.db.WithContext(c).Table(common.ActivityTableName).Select("sum(people_num)").Scan(&totalParticipate)

	var endedNum int64
	a.db.WithContext(c).Table(common.ActivityTableName).Where("end_time < ?", timestamppb.Now()).Count(&endedNum)

	var ongoingNum int64
	a.db.WithContext(c).Table(common.ActivityTableName).Where("start_time < ? and end_time > ?", timestamppb.Now(), timestamppb.Now()).Count(&ongoingNum)

	var upcomingNum int64
	a.db.WithContext(c).Table(common.ActivityTableName).Where("start_time > ?", timestamppb.Now()).Count(&upcomingNum)

	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"total_num":                   totalNum,
		"total_applicants":            totalApplicants,
		"activity_participation_rate": totalParticipate * 100 / totalApplicants,
		"ended_num":                   endedNum,
		"ongoing_num":                 ongoingNum,
		"upcoming_num":                upcomingNum,
	}))
}

func (a *ActivityCtrl) checkoutActivityType(Type string) bool {
	switch Type {
	case "parent_child_interactions", "handmade_diy", "theme_experience", "role_play":
		return true
	default:
		return false
	}
}
