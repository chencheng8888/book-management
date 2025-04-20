package controller

// AddActivityReq 定义了新增活动请求的参数结构
// 包含活动名称、活动类型、开始时间、结束时间、负责人、联系电话和地址等字段
type AddActivityReq struct {
	Info ActivityInfo `json:"info" binding:"required"`
}

// AddActivityResp 定义了新增活动的响应结构
// 包含状态码、返回数据（活动ID）和消息
type AddActivityResp struct {
	// Code 响应状态码，必填项
	Code int `json:"code" binding:"required"`

	// Data 返回的数据结构
	Data struct {
		// ActivityID 新增活动的唯一标识符，必填项
		ActivityID uint64 `json:"activity_id" binding:"required"`
	} `json:"data" binding:"required"`

	// Msg 响应消息，必填项
	Msg string `json:"msg" binding:"required"`
}

// UpdateActivityReq 定义了更新活动请求的参数结构
// 包含活动ID、活动名称、活动类型、开始时间、结束时间和负责人等字段
type UpdateActivityReq struct {
	// ActivityID 活动唯一标识符，必填项
	ActivityID uint64 `json:"activity_id" binding:"required"`

	Info ActivityInfo `json:"info" binding:"required"`
}

// UpdateActivityResp 定义了更新活动的响应结构
// 包含状态码、返回数据（活动ID）和消息
type UpdateActivityResp struct {
	// Code 响应状态码，必填项
	Code int `json:"code" binding:"required"`

	// Msg 响应消息，必填项
	Msg string `json:"msg" binding:"required"`

	// Data 返回的数据结构
	Data struct {
	} `json:"data" binding:"required"`
}

// QueryActivityReq 定义了查询活动请求的参数结构
// 包含分页信息（页码、每页大小）和可选的状态过滤条件
type QueryActivityReq struct {
	// Page 当前页码，必填项
	Page int `json:"page" form:"page" binding:"required"`

	// PageSize 每页显示的数量，必填项
	PageSize int `json:"page_size" form:"page_size" binding:"required"`

	// Status 活动状态，可选项,活动状态有pending，ongoing，ended
	Status *string `json:"status" form:"status"`
}

type ActivityInfo struct {
	// ActivityName 活动名称，必填项
	ActivityName string `json:"activity_name" binding:"required"`

	// ActivityType 活动类型，必填项，暂时有parent_child_interactions
	ActivityType string `json:"activity_type" binding:"required"`

	// StartTime 活动开始时间，必填项，时间格式为"2025-01-01 12:00:00"
	StartTime string `json:"start_time" binding:"required"`

	// EndTime 活动结束时间，必填项，时间格式为"2025-01-01 12:00:00"
	EndTime string `json:"end_time" binding:"required"`

	// Manager 负责人姓名，必填项
	Manager string `json:"manager" binding:"required"`

	// Phone 联系电话，必填项
	Phone string `json:"phone" binding:"required"`

	// Addr 活动地址，必填项
	Addr string `json:"addr" binding:"required"`
}

type Activity struct {
	ActivityID uint64       `json:"activity_id" binding:"required"`
	Info       ActivityInfo `json:"info" binding:"required"`
}

type QueryActivityResp struct {
	Code int `json:"code" binding:"required"`
	Data struct {
		Activitys   []Activity `json:"activitys" binding:"required"`
		CurrentPage int        `json:"current_page" binding:"required"`
		TotalPage   int        `json:"total_page" binding:"required"`
		Total       int        `json:"total" binding:"required"`
	} `json:"data" binding:"required"`
}
