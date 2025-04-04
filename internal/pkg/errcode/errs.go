package errcode

var (
	ParseRequestError       = NewErr(400, "Parse Request Error")
	ParamError              = NewErr(400, "param is illegal")
	GenerateVerifyCodeError = NewErr(500, "generate Verify Code Error")
	UserNotFoundError       = NewErr(400, "User Not Found")
	PasswordError           = NewErr(400, "Password Error")
	VerificationCodeError   = NewErr(400, "Verification Code Error")
	GenerateTokenError      = NewErr(500, "Generate Token Error")

	AddBookStockError     = NewErr(1001, "Add Book Stock Error")
	GenerateBookIDError   = NewErr(1002, "Generate Book ID Error")
	SearchDataError       = NewErr(1003, "Search data Error")
	PageError             = NewErr(1004, "page Error")
	InsufficientBookStock = NewErr(1005, "Insufficient book stock")
	AddBookBorrowError    = NewErr(1006, "Add Book Borrow Error")
)
