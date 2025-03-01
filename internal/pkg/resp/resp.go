package resp

import (
	"book-management/internal/pkg/errcode"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"net/http"
)

var (
	SuccessResp = NewResponse(200, "success", nil)
	BadResp     = NewResponse(400, "Request failed", nil)
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func WithData(resp *Response, data interface{}) *Response {
	resp.Data = data
	return resp
}

func NewResponse(code int, msg string, data interface{}) *Response {
	return &Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func NewRespFromErr(err error) *Response {
	if err == nil {
		return SuccessResp
	}

	var MyErr *errcode.Err
	if !errors.As(err, &MyErr) {
		return BadResp
	}

	code := http.StatusBadRequest
	if MyErr.Code() <= 600 {
		code = MyErr.Code()
	}
	return NewResponse(code, MyErr.Msg(), nil)
}

func SendResp(c *gin.Context, resp *Response) {
	c.JSON(resp.Code, resp)
}
