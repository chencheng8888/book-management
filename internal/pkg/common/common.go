package common

// 借阅状态
const (
	BookStatusWaitingReturn = "waiting_return" //待归还
	BookStatusReturned      = "returned"       //已归还
	BookStatusOverdue       = "overdue"        //逾期
)

// 库存状态
const (
	Adequate     = "adequate"      //库存充足
	EarlyWarning = "early_warning" //库存预警
	Shortage     = "shortage"      //库存短缺
)

// 绘本类别
const (
	ChildrenStory    = "children_story"    //儿童故事
	ScienceKnowledge = "science_knowledge" //科普知识
	ArtEnlightenment = "art_enlightenment" //艺术启蒙
)

const (
	BookTableName       = "books"
	BookStockTableName  = "book_stocks"
	BookCopyTableName   = "book_copy"
	BookBorrowTableName = "book_borrow"
)
