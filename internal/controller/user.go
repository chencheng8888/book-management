package controller

import (
	"book-management/internal/pkg/req"
	"book-management/internal/pkg/resp"
	"context"
	"github.com/gin-gonic/gin"
)

type UserSvc interface {
	SearchUser(ctx context.Context, req SearchUserReq) ([]User, int, error)
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

	users, totalPage, err := u.userSvc.SearchUser(c, searchUserReq)
	if err != nil {
		resp.SendResp(c, resp.NewRespFromErr(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"users":        users,
		"total_page":   totalPage,
		"current_page": searchUserReq.Page,
	}))
}
