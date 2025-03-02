package controller

import (
	"book-management/internal/pkg/common"
	"book-management/internal/pkg/errcode"
	"book-management/internal/pkg/req"
	"book-management/internal/pkg/resp"
	"context"
	"github.com/gin-gonic/gin"
)

type BookStockSvc interface {
	AddStock(ctx context.Context, req AddStockReq) error
	SearchBookStockByID(ctx context.Context, req SearchStockByBookIDReq) (Book, error)
	FuzzyQueryBookStock(ctx context.Context, req FuzzyQueryBookStockReq, total *int) ([]Book, error)
	ListBookStock(ctx context.Context, req ListBookStockReq, total *int) ([]Book, error)
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
	if err := req.ParseRequest(c, &addStockReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	if ok := checkCategory(addStockReq.Category); !ok {
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
// @Param object body SearchStockByBookIDReq true "查询请求"
// @Success 200 {object} SearchStockByBookIDResp
// @Router /api/v1/book/stock/searchByID [get]
func (b *BookStockCtrl) SearchStockByBookID(c *gin.Context) {
	var searchByBookIDReq SearchStockByBookIDReq
	if err := req.ParseRequest(c, &searchByBookIDReq); err != nil {
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
// @Description 模糊查询库存信息
// @Tags 库存
// @Accept application/json
// @Produce application/json
// @Param object body FuzzyQueryBookStockReq true "查询请求"
// @Success 200 {object} FuzzyQueryBookStockResp
// @Router /api/v1/book/stock/fuzzy_query [get]
func (b *BookStockCtrl) FuzzyQueryBookStock(c *gin.Context) {
	var fuzzyQueryBookStockReq FuzzyQueryBookStockReq
	if err := req.ParseRequest(c, &fuzzyQueryBookStockReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	if ok := checkCategory(fuzzyQueryBookStockReq.Category); !ok {
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

func (b *BookStockCtrl) ListBookStock(c *gin.Context) {
	var listBookStockReq ListBookStockReq
	if err := req.ParseRequest(c, &listBookStockReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	var totalPage int
	books, err := b.stockSvc.ListBookStock(c, listBookStockReq, &totalPage)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"books":        books,
		"current_page": listBookStockReq.Page,
		"total_page":   totalPage,
	}))
}

func checkCategory(category string) bool {
	switch category {
	case common.ChildrenStory, common.ScienceKnowledge, common.ArtEnlightenment:
		return true
	default:
		return false
	}
}
