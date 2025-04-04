package controller

//Resp主要是为了帮助生成接口文档

type BorrowBookReq struct {
	BookID             uint64 `json:"book_id" binding:"required"`              // 书本ID【这个你可以理解为一类书，比如《高等数学》】
	BorrowerID         uint64 `json:"borrower_id" binding:"required"`          // 借阅者ID
	ExpectedReturnTime string `json:"expected_return_time" binding:"required"` // 预计归还时间,格式为"2006-01-02"
}

type BorrowBookResp struct {
	Code int    `json:"code" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
	Data struct {
		BookID uint64 `json:"book_id" binding:"required"` // 实际借阅的书本ID
		CopyID uint64 `json:"copy_id" binding:"required"` // 副本ID	【这个你可以理解为具体一本书，比如《高等数学》的第一本】
	} `json:"data" binding:"required"`
}

type BookBorrowRecord struct {
	BookID           uint64 `json:"book_id" binding:"required"`            // 书本ID【这个你可以理解为一类书，比如《高等数学》】
	CopyID           uint64 `json:"copy_id" binding:"required"`            // 副本ID	【这个你可以理解为具体一本书，比如《高等数学》的第一本
	UserID           uint64 `json:"user_id" binding:"required"`            //用户ID
	UserName         string `json:"user_name" binding:"required"`          //用户名
	ShouldReturnTime string `json:"should_return_time" binding:"required"` //应该归还的时间
	ReturnStatus     string `json:"return_status" binding:"required"`      //归还状态
}

type QueryBookBorrowRecordReq struct {
	QueryStatus *string `json:"query_status"`                                  //借阅状态的查询条件
	Page        int     `json:"page" form:"page" binding:"required"`           //第几页
	PageSize    int     `json:"page_size" form:"page_size" binding:"required"` //每页大小
}

type QueryBookBorrowRecordResp struct {
	Code int    `json:"code" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
	Data struct {
		BorrowRecords []BookBorrowRecord `json:"borrow_records" binding:"required"`
		CurrentPage   int                `json:"current_page" binding:"required"` //当前页
		TotalPage     int                `json:"total_page" binding:"required"`   //总数
	} `json:"data" binding:"required"`
}

type UpdateBorrowStatusReq struct {
	BookID uint64 `json:"book_id" binding:"required"` // 书本ID【这个你可以理解为一类书，比如《高等数学》】
	CopyID uint64 `json:"copy_id" binding:"required"` // 副本ID	【这个你可以理解为具体一本书，比如《高等数学》的第一本】
	Status string `json:"status" binding:"required"`  // 要更新的状态,取值有[waiting_return,returned,overdue]
}

type UpdateBorrowStatusResp struct {
	Code int         `json:"code" binding:"required"`
	Msg  string      `json:"msg" binding:"required"`
	Data interface{} `json:"data" binding:"required"`
}

type QueryStatisticsBorrowRecordsReq struct {
	Pattern string `json:"pattern" binding:"required"` //示例值 "week" "month" "year"
}

type QueryStatisticsBorrowRecordsResp struct {
	Code int    `json:"code" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
	Data struct {
		ChildrenStoryNum    int `json:"children_story_num" binding:"required"`
		ScienceKnowledgeNum int `json:"science_knowledge_num" binding:"required"`
		ArtEnlightenmentNum int `json:"art_enlightenment_num" binding:"required"`
	}
}
