package errcode

var (
	ParseRequestError = NewErr(400, "Parse Request Error")
	ParamError        = NewErr(400, "param is illegal")

	AddBookStockError   = NewErr(1001, "Add Book Stock Error")
	GenerateBookIDError = NewErr(1002, "Generate Book ID Error")
	SearchBookError     = NewErr(1003, "Search Book Error")
	PageError           = NewErr(1004, "Page Error")
)
