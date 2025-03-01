package errcode

var (
	ParseRequestBodyError = NewErr(400, "Parse Request Body Error")
	ParamError            = NewErr(400, "param is illegal")
)
