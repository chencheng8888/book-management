package controller

type GetVolunteerInfosReq struct {
	Page     int     `json:"page" form:"page" binding:"required"`
	PageSize int     `json:"pageSize" form:"page_size" binding:"required"`
	ID       *uint64 `json:"id" form:"id"`
}

type VolunteerInfo struct {
	ID                    uint64 `json:"id"`                    //ID
	Name                  string `json:"name"`                  //姓名
	Phone                 string `json:"phone"`                 //手机号
	Age                   int    `json:"age"`                   //年龄
	ServiceTimePreference string `json:"serviceTimePreference"` //服务时间偏好
	ExpertiseArea         string `json:"expertiseArea"`         //专业领域
	CreatedAt             string `json:"createdAt"`             //创建时间
}

type GetVolunteerInfosResp struct {
	Code int `json:"code" binding:"required"`
	Data struct {
		Volunteers  []VolunteerInfo `json:"volunteers" binding:"required"`   // 志愿者列表
		TotalPage   int             `json:"total_page" binding:"required"`   // 总页数
		CurrentPage int             `json:"current_page" binding:"required"` // 当前页码
		Total       int             `json:"total" binding:"required"`        // 总记录数
	} `json:"data" binding:"required"`
	Msg string `json:"msg" binding:"required"` // 响应消息
}

type CreateVolunteerReq struct {
	Name                  string `json:"name" binding:"required"`                  // 姓名
	Phone                 string `json:"phone" binding:"required"`                 // 手机号
	Age                   int    `json:"age" binding:"required"`                   // 年龄
	ServiceTimePreference string `json:"serviceTimePreference" binding:"required"` // 服务时间偏好
	ExpertiseArea         string `json:"expertiseArea" binding:"required"`         // 专业领域
}

type CreateVolunteerResp struct {
	Code int `json:"code" binding:"required"`
	Data struct {
		VolunteerID uint64 `json:"volunteer_id" binding:"required"`
	} `json:"data" binding:"required"`
	Msg string `json:"msg" binding:"required"`
}

type GetVolunteerApplicationsReq struct {
	Page     int `json:"page" form:"page" binding:"required"`
	PageSize int `json:"pageSize" form:"page_size" binding:"required"`
}

type VolunteerApplicationInfo struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Age   int    `json:"age"`
}

type GetVolunteerApplicationsResp struct {
	Code int `json:"code" binding:"required"`
	Data struct {
		Applications []VolunteerApplicationInfo `json:"applications" binding:"required"` // 申请列表
		TotalPage    int                        `json:"total_page" binding:"required"`   // 总页数
		CurrentPage  int                        `json:"current_page" binding:"required"` // 当前页码
		Total        int                        `json:"total" binding:"required"`        // 总记录数
	} `json:"data" binding:"required"`
	Msg string `json:"msg" binding:"required"`
}
