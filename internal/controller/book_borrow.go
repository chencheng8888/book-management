package controller

import (
	"book-management/internal/pkg/common"
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/req"
	"book-management/internal/pkg/resp"
	"book-management/internal/pkg/tool"
	"context"
	"github.com/gin-gonic/gin"
)

type BookBorrowSvc interface {
	AddBorrowBookRecord(ctx context.Context, req BorrowBookReq) (bookID uint64, copyID uint64, err error)
	QueryBookBorrowRecord(ctx context.Context, req QueryBookBorrowRecordReq, total *int) ([]BookBorrowRecord, error)
	UpdateBorrowStatus(ctx context.Context, req UpdateBorrowStatusReq) error
	GetStatisticBorrowRecords(ctx context.Context, req QueryStatisticsBorrowRecordsReq) (map[string]int, error)
	GetAvailableCopyBook(ctx context.Context, req GetAvailableCopyBookReq) ([]uint64, error)
}

type BookBorrowCtrl struct {
	borrowSvc BookBorrowSvc
}

func NewBookBorrowCtrl(borrowSvc BookBorrowSvc) *BookBorrowCtrl {
	return &BookBorrowCtrl{
		borrowSvc: borrowSvc,
	}
}

func (b *BookBorrowCtrl) RegisterRoute(r *gin.Engine) {
	g := r.Group("/api/v1/book/borrow")
	{
		g.POST("/add", b.BorrowBook)
		g.GET("/query", b.QueryBookBorrowRecord)
		g.PUT("/update_status", b.UpdateBorrowStatus)
		g.GET("/query_statistics", b.QueryStatisticsBorrowRecords)
		g.GET("/get_available", b.GetAvailableCopyBook)
	}
}

// GetAvailableCopyBook 获取可借用的书籍
// @Summary 获取可借用的书籍
// @Description 获取可借用的书籍,当返回的数量等于page_size+1时，则代表还有下一页，否则，没有
// @Tags 借书
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object query GetAvailableCopyBookReq true "请求"
// @Success 200 {object} GetAvailableCopyBookResp
// @Router /api/v1/book/borrow/get_available [get]
func (b *BookBorrowCtrl) GetAvailableCopyBook(c *gin.Context) {
	var getAvailableCopyBookReq GetAvailableCopyBookReq
	if err := req.ParseRequestBody(c, &getAvailableCopyBookReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	copyIds, err := b.borrowSvc.GetAvailableCopyBook(c, getAvailableCopyBookReq)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"copy_ids": copyIds,
	}))
}

// BorrowBook 借书
// @Summary 借书
// @Description 借书接口
// @Tags 借书
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object body BorrowBookReq true "借书请求"
// @Success 200 {object} BorrowBookResp
// @Router /api/v1/book/borrow/add [post]
func (b *BookBorrowCtrl) BorrowBook(c *gin.Context) {
	//解析参数
	var borrowBookReq BorrowBookReq
	if err := req.ParseRequestBody(c, &borrowBookReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	//执行
	bookID, copyID, err := b.borrowSvc.AddBorrowBookRecord(c, borrowBookReq)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	//发送响应
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"book_id": bookID,
		"copy_id": copyID,
	}))
}

// QueryBookBorrowRecord 查询借书记录
// @Summary 查询借书记录
// @Description 查询借书记录
// @Tags 借书
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object query QueryBookBorrowRecordReq true "查询请求"
// @Success 200 {object} QueryBookBorrowRecordResp
// @Router /api/v1/book/borrow/query [get]
func (b *BookBorrowCtrl) QueryBookBorrowRecord(c *gin.Context) {
	//解析参数
	var queryBookBorrowRecordReq QueryBookBorrowRecordReq

	if err := req.ParseRequestQuery(c, &queryBookBorrowRecordReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	if queryBookBorrowRecordReq.QueryStatus != nil {
		if *queryBookBorrowRecordReq.QueryStatus == "" {
			queryBookBorrowRecordReq.QueryStatus = nil
		} else {
			if !tool.CheckBorrowStatus(*queryBookBorrowRecordReq.QueryStatus) {
				resp.SendResp(c, resp.NewRespFromErr(errcode.ParamError))
				return
			}
		}
	}

	var totalNum int
	records, err := b.borrowSvc.QueryBookBorrowRecord(c, queryBookBorrowRecordReq, &totalNum)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	//发送响应
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"borrow_records": records,
		"current_page":   queryBookBorrowRecordReq.Page,
		"total_page":     tool.GetPage(totalNum, queryBookBorrowRecordReq.PageSize),
		"total_num":      totalNum,
	}))
}

// UpdateBorrowStatus 更新借阅状态
// @Summary 更新借阅状态
// @Description 更新借阅状态
// @Tags 借书
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object body UpdateBorrowStatusReq true "更新请求"
// @Success 200 {object} UpdateBorrowStatusResp
// @Router /api/v1/book/borrow/update_status [put]
func (b *BookBorrowCtrl) UpdateBorrowStatus(c *gin.Context) {
	var updatereq UpdateBorrowStatusReq
	if err := req.ParseRequestBody(c, &updatereq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	if !tool.CheckBorrowStatus(updatereq.Status) {
		resp.SendResp(c, resp.NewRespFromErr(errcode.ParamError))
		return
	}

	if err := b.borrowSvc.UpdateBorrowStatus(c, updatereq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.SuccessResp)
}

// QueryStatisticsBorrowRecords 获取统计借阅记录
// @Summary 获取统计借阅记录
// @Description 获取统计借阅记录
// @Tags 借书
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object query QueryStatisticsBorrowRecordsReq true "获取统计数据请求"
// @Success 200 {object} QueryStatisticsBorrowRecordsResp
// @Router /api/v1/book/borrow/query_statistics [get]
func (b *BookBorrowCtrl) QueryStatisticsBorrowRecords(c *gin.Context) {
	var statisticsBorrowRecordsReq QueryStatisticsBorrowRecordsReq
	if err := req.ParseRequestQuery(c, &statisticsBorrowRecordsReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	switch statisticsBorrowRecordsReq.Pattern {
	case "week", "month", "year":
	default:
		resp.SendResp(c, resp.NewRespFromErr(errcode.ParamError))
		return
	}
	result, err := b.borrowSvc.GetStatisticBorrowRecords(c, statisticsBorrowRecordsReq)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"children_story_num":    result[common.ChildrenStory],
		"science_knowledge_num": result[common.ScienceKnowledge],
		"art_enlightenment_num": result[common.ArtEnlightenment],
	}))
}
