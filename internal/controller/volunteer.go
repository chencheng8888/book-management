package controller

import (
	"book-management/internal/pkg/req"
	"book-management/internal/pkg/resp"
	"book-management/internal/pkg/tool"
	"context"
	"github.com/gin-gonic/gin"
)

type VolunteerSvc interface {
	GetVolunteerInfos(ctx context.Context, req GetVolunteerInfosReq, totalNum *int) ([]VolunteerInfo, error)
	CreateVolunteer(ctx context.Context, req CreateVolunteerReq) (uint64, error)
	GetApplications(ctx context.Context, req GetVolunteerApplicationsReq, total *int) ([]VolunteerApplicationInfo, error)
}

type VolunteerController struct {
	svc VolunteerSvc
}

func NewVolunteerController(svc VolunteerSvc) *VolunteerController {
	return &VolunteerController{svc: svc}
}

func (v *VolunteerController) RegisterRoute(r *gin.Engine) {
	g := r.Group("/api/v1/volunteer")
	{
		g.GET("/list_application", v.GetApplications)
		g.GET("/query", v.GetVolunteerInfos)
		g.POST("/create", v.CreateVolunteer)
	}
}

// GetVolunteerInfos 查询志愿者信息
// @Summary 查询志愿者信息
// @Description 分页查询志愿者信息或根据ID查询单个志愿者
// @Tags 志愿者管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object query GetVolunteerInfosReq true "查询参数"
// @Success 200 {object} GetVolunteerInfosResp
// @Router /api/v1/volunteer/query [get]
func (v *VolunteerController) GetVolunteerInfos(c *gin.Context) {
	var rreq GetVolunteerInfosReq
	if err := req.ParseRequestQuery(c, &rreq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	if rreq.ID != nil && *rreq.ID == 0 {
		rreq.ID = nil
	}

	var totalNum int
	infos, err := v.svc.GetVolunteerInfos(c, rreq, &totalNum)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"volunteers":   infos,
		"total_page":   tool.GetPage(totalNum, rreq.PageSize),
		"current_page": rreq.Page,
		"total":        totalNum,
	}))
}

// CreateVolunteer 创建志愿者
// @Summary 创建志愿者
// @Description 新增志愿者信息
// @Tags 志愿者管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object body CreateVolunteerReq true "创建志愿者请求"
// @Success 200 {object} CreateVolunteerResp
// @Router /api/v1/volunteer/create [post]
func (v *VolunteerController) CreateVolunteer(c *gin.Context) {
	var rreq CreateVolunteerReq
	if err := req.ParseRequestBody(c, &rreq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	volunteerID, err := v.svc.CreateVolunteer(c, rreq)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"volunteer_id": volunteerID,
	}))
}

// GetApplications 获取申请志愿者列表
// @Summary 获取申请志愿者列表
// @Description 分页获取申请志愿者的信息，包括姓名、电话号码和年龄
// @Tags 志愿者管理
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param request query GetVolunteerApplicationsReq true "请求参数"
// @Success 200 {object} GetVolunteerApplicationsResp
// @Router /api/v1/volunteer/list_application [get]
func (v *VolunteerController) GetApplications(c *gin.Context) {
	var rreq GetVolunteerApplicationsReq
	if err := req.ParseRequestQuery(c, &rreq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	var total int
	applications, err := v.svc.GetApplications(c, rreq, &total)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"applications": applications,
		"total_page":   tool.GetPage(total, rreq.PageSize),
		"current_page": rreq.Page,
		"total":        total,
	}))
}
