package controller

type ListDonateRecordsReq struct {
	PageSize int `form:"page_size" binding:"required"` //每页大小
	Page     int `form:"page" binding:"required"`      //当前页
}

type DonateRecords struct {
	UserID     uint64 `json:"user_id" binding:"required"`     //用户ID
	UserName   string `json:"user_name" binding:"required"`   //用户名称
	Phone      string `json:"phone" binding:"required"`       //电话
	BookName   string `json:"book_name" binding:"required"`   //书籍名称
	DonateNum  int    `json:"donate_num" binding:"required"`  //捐赠数目
	DonateTime string `json:"donate_time" binding:"required"` //捐赠时间
}

type ListDonateRecordsResp struct {
	Code int `json:"code" binding:"required"`
	Data struct {
		DonateRecords []DonateRecords `json:"donate_records" binding:"required"` //捐赠记录
		CurrentPage   int             `json:"current_page" binding:"required"`   //当前页
		TotalPage     int             `json:"total_page" binding:"required"`     //总页数
		TotalNum      int             `json:"total_num" binding:"required"`      //总数量
	} `json:"data" binding:"required"`
	Msg string `json:"msg" binding:"required"`
}

type GetDonationRankingReq struct {
	Top int `form:"top" binding:"required"`
}

type Rank struct {
	UserID      uint64 `json:"user_id" binding:"required"`      //用户ID
	UserName    string `json:"user_name" binding:"required"`    //用户名称
	DonateNum   int    `json:"donate_num" binding:"required"`   //捐赠数目
	DonateTimes int    `json:"donate_times" binding:"required"` //捐赠次数
	UpdatedAt   string `json:"updated_at" binding:"required"`   //最近捐赠时间
}
type GetDonationRankingResp struct {
	Code int `json:"code" binding:"required"`
	Data struct {
		Rankings []Rank `json:"rankings" binding:"required"`
	} `json:"data" binding:"required"`
}
