package service

type Book struct {
	BookID      uint64 `json:"book_id"`       //书本ID
	Name        string `json:"name" `         //书本名称
	Author      string `json:"author"`        // 作者
	Publisher   string `json:"publisher"`     //出版社
	Category    string `json:"category"`      //类别
	Stock       uint   `json:"stock"`         //库存数量
	StockStatus string `json:"stock_status" ` //库存状态
	StockWhere  string `json:"stock_where"`   //库存位置
	CreatedAt   string `json:"created_at"`    //入库时间
}
