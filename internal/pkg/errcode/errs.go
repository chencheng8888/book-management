package errcode

var (
	ParseRequestError = NewErr(400, "Parse Request Error")
	ParamError        = NewErr(400, "param is illegal")

	AddBookStockError     = NewErr(1001, "Add Book Stock Error")
	GenerateBookIDError   = NewErr(1002, "Generate Book ID Error")
	SearchDataError       = NewErr(1003, "Search data Error")
	PageError             = NewErr(1004, "Page Error")
	InsufficientBookStock = NewErr(1005, "Insufficient book stock")
	AddBookBorrowError    = NewErr(1006, "Add Book Borrow Error")
)
