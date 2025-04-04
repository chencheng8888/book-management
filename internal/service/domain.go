package service

import (
	"book-management/internal/controller"
	"time"
)

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
	BorrowerID   uint64    `json:"borrower_id"`   // 借阅者ID
	Borrower     string    `json:"borrower"`      // 借阅者姓名
	CopyID       uint64    `json:"copy_id"`       // 副本ID
	BorrowTime   time.Time `json:"borrow_time"`   // 借阅时间
	ExpectedTime time.Time `json:"expected_time"` // 预计归还时间
	ReturnStatus string    `json:"return_status"` // 归还状态
}

type User struct {
	controller.User
}

func toControllerBook(book Book) controller.Book {
	return controller.Book{
		BookID:      book.BookID,
		Name:        book.Info.Name,
		Author:      book.Info.Author,
		Publisher:   book.Info.Publisher,
		Category:    book.Info.Category,
		Stock:       book.Stock.Stock,
		StockStatus: book.Stock.Status,
		StockWhere:  book.Stock.Where,
		CreatedAt:   convertTimeFormat(book.Stock.AddedTime),
	}
}
func toControllerBookBorrowRecord(record BookBorrowRecord) controller.BookBorrowRecord {
	return controller.BookBorrowRecord{
		BookID:           record.BookID,
		CopyID:           record.CopyID,
		UserID:           record.BorrowerID,
		UserName:         record.Borrower,
		ShouldReturnTime: convertTimeFormat(record.ExpectedTime),
		ReturnStatus:     record.ReturnStatus,
	}
}

func batchToControllerBook(books []Book) []controller.Book {
	var result []controller.Book
	for _, book := range books {
		result = append(result, toControllerBook(book))
	}
	return result
}

func batchToControllerBookBorrowRecord(records []BookBorrowRecord) []controller.BookBorrowRecord {
	var result []controller.BookBorrowRecord
	for _, record := range records {
		result = append(result, toControllerBookBorrowRecord(record))
	}
	return result
}

func batchToControllerUser(users []User) []controller.User {
	var result = make([]controller.User, 0, len(users))
	for _, user := range users {
		result = append(result, user.User)
	}
	return result
}
