package service

import "time"

type BookInfo struct {
	//BookID    uint64 `json:"book_id"`   //书本ID
	Name      string `json:"name" `     //书本名称
	Author    string `json:"author"`    // 作者
	Publisher string `json:"publisher"` //出版社
	Category  string `json:"category"`  //类别
}

type BookStock struct {
	//BookID    uint64    `json:"book_id"`    // 书本ID
	Stock     uint      `json:"stock"`      //库存
	Status    string    `json:"status"`     //库存状态
	Where     string    `json:"where"`      //库存位置
	AddedTime time.Time `json:"added_time"` //入库时间
}
type Book struct {
	BookID uint64 `json:"book_id"` // 书本ID
	Info   BookInfo
	Stock  BookStock
}

type BookBorrowRecord struct {
	BookID       uint64    `json:"book_id"`       // 书本ID
	BorrowerID   string    `json:"borrower_id"`   // 借阅者ID
	Borrower     string    `json:"borrower"`      // 借阅者姓名
	CopyID       uint64    `json:"copy_id"`       // 副本ID
	BorrowTime   time.Time `json:"borrow_time"`   // 借阅时间
	ExpectedTime time.Time `json:"expected_time"` // 预计归还时间
	ReturnStatus string    `json:"return_status"` // 归还状态
}
