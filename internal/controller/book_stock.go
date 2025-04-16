package controller

import (
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/req"
	"book-management/internal/pkg/resp"
	"book-management/internal/pkg/tool"
	"context"

	"github.com/gin-gonic/gin"
)

type BookStockSvc interface {
	AddStock(ctx context.Context, req AddStockReq, bookID *uint64) error
	FuzzyQueryBookStock(ctx context.Context, req FuzzyQueryBookStockReq, total *int) ([]Book, error)
	//ListBookStock(ctx context.Context, req ListBookStockReq, total *int) ([]Book, error)
}

type BookStockCtrl struct {
	stockSvc BookStockSvc
}

func NewBookStockCtrl(stockSvc BookStockSvc) *BookStockCtrl {
	return &BookStockCtrl{
		stockSvc: stockSvc,
	}
}

func (b *BookStockCtrl) RegisterRoute(r *gin.Engine) {
	g := r.Group("/api/v1/book/stock")
	{
		g.POST("/add", b.AddStock)
		g.GET("/fuzzy_query", b.FuzzyQueryBookStock)
		//g.GET("/list", b.ListBookStock)
	}
}

// AddStock 添加库存
// @Summary 添加库存
// @Description 添加库存接口
// @Tags 库存
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object body AddStockReq true "增加库存请求"
// @Success 200 {object} AddStockResp
// @Router /api/v1/book/stock/add [post]
func (b *BookStockCtrl) AddStock(c *gin.Context) {
	//解析参数
	var addStockReq AddStockReq
	if err := req.ParseRequestBody(c, &addStockReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	if ok := tool.CheckCategory(addStockReq.Category); !ok {
		resp.SendResp(c, resp.NewRespFromErr(errcode.ParamError))
		return
	}

	if addStockReq.UserID != nil && *addStockReq.UserID == 0 {
		addStockReq.UserID = nil
	}

	//执行
	var bookID uint64
	if err := b.stockSvc.AddStock(c, addStockReq, &bookID); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	//发送响应
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"book_id": bookID,
	}))
}

// FuzzyQueryBookStock 模糊查询库存信息
// @Summary 模糊查询库存信息
// @Description 模糊查询库存信息,没有任何查询条件就是直接列出数据
// @Tags 库存
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "鉴权"
// @Param object query FuzzyQueryBookStockReq true "查询请求"
// @Success 200 {object} FuzzyQueryBookStockResp
// @Router /api/v1/book/stock/fuzzy_query [get]
func (b *BookStockCtrl) FuzzyQueryBookStock(c *gin.Context) {
	var fuzzyQueryBookStockReq FuzzyQueryBookStockReq
	if err := req.ParseRequestQuery(c, &fuzzyQueryBookStockReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	if fuzzyQueryBookStockReq.BookID != nil && *fuzzyQueryBookStockReq.BookID == 0 {
		fuzzyQueryBookStockReq.BookID = nil
	}
	if fuzzyQueryBookStockReq.Category != nil && *fuzzyQueryBookStockReq.Category == "" {
		fuzzyQueryBookStockReq.Category = nil
	}
	if fuzzyQueryBookStockReq.Name != nil && *fuzzyQueryBookStockReq.Name == "" {
		fuzzyQueryBookStockReq.Name = nil
	}
	if fuzzyQueryBookStockReq.Author != nil && *fuzzyQueryBookStockReq.Author == "" {
		fuzzyQueryBookStockReq.Author = nil
	}

	if fuzzyQueryBookStockReq.Category != nil {
		if !tool.CheckCategory(*fuzzyQueryBookStockReq.Category) {
			resp.SendResp(c, resp.NewRespFromErr(errcode.ParamError))
			return
		}
	}

	var totalNum int

	books, err := b.stockSvc.FuzzyQueryBookStock(c, fuzzyQueryBookStockReq, &totalNum)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"books":        books,
		"current_page": fuzzyQueryBookStockReq.Page,
		"total_page":   tool.GetPage(totalNum, fuzzyQueryBookStockReq.PageSize),
		"total_num":    totalNum,
	}))
}
