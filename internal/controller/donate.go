package controller

import (
	"book-management/internal/pkg/req"
	"book-management/internal/pkg/resp"
	"book-management/internal/pkg/tool"
	"context"

	"github.com/gin-gonic/gin"
)

type BookDonateSvc interface {
	ListDonateRecordsReq(ctx context.Context, req ListDonateRecordsReq, total *int) ([]DonateRecords, error)
	GetDonationRanking(ctx context.Context, req GetDonationRankingReq) ([]Rank, error)
}

type BookDonateCtrl struct {
	donateSvc BookDonateSvc
}

func NewBookDonateCtrl(donateSvc BookDonateSvc) *BookDonateCtrl {
	return &BookDonateCtrl{
		donateSvc: donateSvc,
	}
}

func (b *BookDonateCtrl) RegisterRoute(r *gin.Engine) {
	g := r.Group("/api/v1/book/donate")
	{
		g.GET("/list", b.ListDonateRecords)
		g.GET("/get_ranking", b.GetDonationRanking)
	}
}

// ListDonateRecords 列出捐赠记录
// @Summary 列出捐赠记录
// @Description 列出捐赠记录
// @Tags 捐赠
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object query ListDonateRecordsReq true "增加库存请求"
// @Success 200 {object} ListDonateRecordsResp
// @Router /api/v1/book/donate/list [get]
func (b *BookDonateCtrl) ListDonateRecords(c *gin.Context) {
	var listDonateRecordsReq ListDonateRecordsReq
	if err := req.ParseRequestQuery(c, &listDonateRecordsReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	var totalNum int
	records, err := b.donateSvc.ListDonateRecordsReq(c, listDonateRecordsReq, &totalNum)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"donate_records": records,
		"total_page":     tool.GetPage(totalNum, listDonateRecordsReq.PageSize),
		"current_page":   listDonateRecordsReq.Page,
		"total_num":      totalNum,
	}))
}

// GetDonationRanking 获取捐赠排名
// @Summary 获取捐赠排名
// @Description 获取捐赠排名
// @Tags 捐赠
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object query GetDonationRankingReq true "增加库存请求"
// @Success 200 {object} GetDonationRankingResp
// @Router /api/v1/book/donate/get_ranking [get]
func (b *BookDonateCtrl) GetDonationRanking(c *gin.Context) {
	var getDonationRankingReq GetDonationRankingReq
	if err := req.ParseRequestQuery(c, &getDonationRankingReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	rankings, err := b.donateSvc.GetDonationRanking(c, getDonationRankingReq)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"rankings": rankings,
	}))
}
