package controller

type SearchUserReq struct {
	UserID   *uint64 `form:"user_id" binding:"required"`
	UserName *string `form:"user_name" binding:"required"`
	Phone    *string `form:"phone" binding:"required"`
	Page     int     `form:"page" binding:"required"`
	PageSize int     `form:"page_size" binding:"required"`
}

type User struct {
	ID        uint64  `json:"id" binding:"required"`       //用户ID
	Name      string  `json:"name" binding:"required"`     //姓名
	Phone     string  `json:"phone" binding:"required"`    //手机号
	Integral  uint    `json:"integral" binding:"required"` //积分
	Gender    string  `json:"gender" binding:"required"`   //性别
	IsVip     bool    `json:"is_vip" binding:"required"`   //是否是会员
	VipLevels *string `json:"vip_levels"`                  //会员等级
	Status    string  `json:"status" binding:"required"`   //用户状态
}

type SearchUserResp struct {
	Code int `json:"code" binding:"required"`
	Data struct {
		Users       []User `json:"users" binding:"required"`
		TotalPage   int    `json:"total_page" binding:"required"`
		CurrentPage int    `json:"current_page" binding:"required"`
	} `json:"data" binding:"required"`
	Msg string `json:"msg" binding:"required"`
}
