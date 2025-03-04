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
	AddStock(ctx context.Context, req AddStockReq) error
	SearchBookStockByID(ctx context.Context, req SearchStockByBookIDReq) (Book, error)
	FuzzyQueryBookStock(ctx context.Context, req FuzzyQueryBookStockReq, total *int) ([]Book, error)
	//ListBookStock(ctx context.Context, req ListBookStockReq, total *int) ([]Book, error)
}

type BookStockCtrl struct {
	stockSvc BookStockSvc
}

func (b *BookStockCtrl) RegisterRoute(r *gin.Engine) {
	g := r.Group("/api/v1/book/stock")
	{
		g.POST("/add", b.AddStock)
		g.GET("/searchByID", b.SearchStockByBookID)
		g.GET("/fuzzy_query", b.FuzzyQueryBookStock)
		//g.GET("/list", b.ListBookStock)
	}
}

// AddStock 添加库存
// @Summary 添加库存
// @Description 添加库存接口，参数的where是可选参数
// @Tags 库存
// @Accept application/json
// @Produce application/json
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

	//执行
	if err := b.stockSvc.AddStock(c, addStockReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	//发送响应
	resp.SendResp(c, resp.SuccessResp)
}

// SearchStockByBookID 精确查询库存信息
// @Summary 根据ID查询库存信息
// @Description 根据ID查询库存信息
// @Tags 库存
// @Accept application/json
// @Produce application/json
// @Param object query SearchStockByBookIDReq true "查询请求"
// @Success 200 {object} SearchStockByBookIDResp
// @Router /api/v1/book/stock/searchByID [get]
func (b *BookStockCtrl) SearchStockByBookID(c *gin.Context) {
	var searchByBookIDReq SearchStockByBookIDReq
	if err := req.ParseRequestQuery(c, &searchByBookIDReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	book, err := b.stockSvc.SearchBookStockByID(c, searchByBookIDReq)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	resp.SendResp(c, resp.WithData(resp.SuccessResp, book))
}

// FuzzyQueryBookStock 模糊查询库存信息
// @Summary 模糊查询库存信息
// @Description 模糊查询库存信息,没有任何查询条件就是直接列出数据
// @Tags 库存
// @Accept application/json
// @Produce application/json
// @Param object query FuzzyQueryBookStockReq true "查询请求"
// @Success 200 {object} FuzzyQueryBookStockResp
// @Router /api/v1/book/stock/fuzzy_query [get]
func (b *BookStockCtrl) FuzzyQueryBookStock(c *gin.Context) {
	var fuzzyQueryBookStockReq FuzzyQueryBookStockReq
	if err := req.ParseRequestQuery(c, &fuzzyQueryBookStockReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	if fuzzyQueryBookStockReq.Category != nil && tool.CheckCategory(*fuzzyQueryBookStockReq.Category) {
		resp.SendResp(c, resp.NewRespFromErr(errcode.ParamError))
		return
	}

	if fuzzyQueryBookStockReq.AddStockTime != nil && tool.IsTimeFormatValid(*fuzzyQueryBookStockReq.AddStockTime, tool.Format2) {
		resp.SendResp(c, resp.NewRespFromErr(errcode.ParamError))
		return
	}

	var totalPage int

	books, err := b.stockSvc.FuzzyQueryBookStock(c, fuzzyQueryBookStockReq, &totalPage)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"books":        books,
		"current_page": fuzzyQueryBookStockReq.Page,
		"total_page":   totalPage,
	}))
}

//// ListBookStock 列出所有库存信息
//// @Summary 列出所有库存信息
//// @Description 列出所有库存信息
//// @Tags 库存
//// @Accept application/json
//// @Produce application/json
//// @Param object query ListBookStockReq true "查询请求"
//// @Success 200 {object} ListBookStockResp
//// @Router /api/v1/book/stock/list [get]
//func (b *BookStockCtrl) ListBookStock(c *gin.Context) {
//	var listBookStockReq ListBookStockReq
//	if err := req.ParseRequestBody(c, &listBookStockReq); err != nil {
//		resp.SendResp(c, resp.NewRespFromErr(err))
//		return
//	}
//	var totalPage int
//	books, err := b.stockSvc.ListBookStock(c, listBookStockReq, &totalPage)
//	if err != nil {
//		resp.SendResp(c, resp.NewRespFromErr(err))
//		return
//	}
//	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
//		"books":        books,
//		"current_page": listBookStockReq.Page,
//		"total_page":   totalPage,
//	}))
//}
