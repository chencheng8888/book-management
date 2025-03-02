package controller

//Resp主要是为了帮助生成接口文档

type AddStockReq struct {
	Name          string  `json:"name" binding:"required"`           //书本名称
	Author        string  `json:"author" binding:"required"`         // 作者
	Publisher     string  `json:"publisher" binding:"required"`      //出版社
	Category      string  `json:"category" binding:"required"`       //类别
	QuantityAdded uint    `json:"quantity_added" binding:"required"` // 添加的库存数目
	Where         *string `json:"where"`                             // 库存位置
}

type AddStockResp struct {
	Code int    `json:"code" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
	Data struct {
		BookID uint64 `json:"book_id" binding:"required"` //书本ID
	} `json:"data"`
}

type Book struct {
	BookID      uint64 `json:"book_id" binding:"required"`       //书本ID
	Name        string `json:"name" binding:"required"`          //书本名称
	Author      string `json:"author" binding:"required"`        // 作者
	Publisher   string `json:"publisher" binding:"required"`     //出版社
	Category    string `json:"category" binding:"required"`      //类别
	Stock       uint   `json:"stock" binding:"required"`         //库存数量
	StockStatus string `json:"stock_status"  binding:"required"` //库存状态
	StockWhere  string `json:"stock_where" binding:"required"`   //库存位置
	CreatedAt   string `json:"created_at" binding:"required"`    //入库时间
}

type SearchStockByBookIDReq struct {
	BookID uint64 `json:"book_id" binding:"required"` // 书本ID
}

type SearchStockByBookIDResp struct {
	Code int    `json:"code" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
	Data Book   `json:"data" binding:"required"` //数据
}

type FuzzyQueryBookStockReq struct {
	Name         *string `json:"name"`                         //书本名称
	Author       *string `json:"author"`                       //作者
	AddStockTime *string `json:"add_stock_time"`               //入库时间
	Category     string  `json:"category" binding:"required"`  //类别
	Page         int     `json:"page" binding:"required"`      //第几页
	PageSize     int     `json:"page_size" binding:"required"` //每页大小
}
type FuzzyQueryBookStockResp struct {
	Code int    `json:"code" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
	Data struct {
		Books       []Book `json:"books" binding:"required"`        //数据
		CurrentPage int    `json:"current_page" binding:"required"` //当前页
		TotalPage   int    `json:"total_page" binding:"required"`   //总数
	} `json:"data" binding:"required"` //数据
}
type ListBookStockReq struct {
	Page     int `json:"page" binding:"required"`      //第几页
	PageSize int `json:"page_size" binding:"required"` //每页大小
}
type ListBookStockResp struct {
	Code int    `json:"code" binding:"required"`
	Msg  string `json:"msg" binding:"required"`
	Data struct {
		Books       []Book `json:"books" binding:"required"`        //数据
		CurrentPage int    `json:"current_page" binding:"required"` //当前页
		TotalPage   int    `json:"total_page" binding:"required"`   //总数
	}
}
