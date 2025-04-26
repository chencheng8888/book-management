package controller

import (
	"book-management/internal/pkg/req"
	"book-management/internal/pkg/resp"
	"book-management/internal/pkg/tool"
	"context"
	"github.com/gin-gonic/gin"
	"reflect"
)

type UserSvc interface {
	SearchUser(ctx context.Context, req SearchUserReq) ([]User, int, error)
	GetVIPStatics(ctx context.Context) (map[string]int, error)
}

type UserCtrl struct {
	userSvc UserSvc
}

func NewUserCtrl(userSvc UserSvc) *UserCtrl {
	return &UserCtrl{
		userSvc: userSvc,
	}
}

func (u *UserCtrl) RegisterRoute(r *gin.Engine) {
	g := r.Group("/api/v1/user")
	{
		g.GET("/search", u.SearchUser)
		g.GET("/vip_statics", u.GetVIPStatics)
	}
}

// SearchUser 查询用户
// @Summary 查询用户
// @Description 查询用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param Authorization header string true "鉴权"
// @Param object query SearchUserReq true "请求参数"
// @Success 200 {object} SearchUserResp "查询成功"
// @Router /api/v1/user/search [get]
func (u *UserCtrl) SearchUser(c *gin.Context) {
	var searchUserReq SearchUserReq
	if err := req.ParseRequestQuery(c, &searchUserReq); err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}

	NormalizePointerFields(&searchUserReq)

	users, totalNum, err := u.userSvc.SearchUser(c, searchUserReq)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"users":        users,
		"total_page":   tool.GetPage(totalNum, searchUserReq.PageSize),
		"current_page": searchUserReq.Page,
		"total_num":    totalNum,
	}))
}

// GetVIPStatics 获取会员的统计数据
// @Summary  获取会员的统计数据
// @Description  获取会员的统计数据
// @Tags 用户
// @Accept json
// @Produce json
// @Param Authorization header string true "鉴权"
// @Param object query GetVIPStaticsReq true "请求参数"
// @Success 200 {object} GetVIPStaticsResp "查询成功"
// @Router /api/v1/user/vip_statics [get]
func (u *UserCtrl) GetVIPStatics(c *gin.Context) {
	mp, err := u.userSvc.GetVIPStatics(c)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"normal_num": mp["normal"],
		"gold_num":   mp["gold"],
		"silver_num": mp["silver"],
	}))
}

// NormalizePointerFields 检查结构体中的指针字段，如果值为零值则设置为 nil
func NormalizePointerFields(v interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return // 确保 v 是非 nil 的指针
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return // 确保 v 指向的是结构体
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// 只处理指针字段
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			// 检查指针指向的值是否为零值
			zeroValue := reflect.Zero(field.Elem().Type()).Interface()
			if reflect.DeepEqual(field.Elem().Interface(), zeroValue) {
				field.Set(reflect.Zero(field.Type())) // 设置为 nil
			}
		}
	}
}
