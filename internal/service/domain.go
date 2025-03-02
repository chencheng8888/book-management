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
