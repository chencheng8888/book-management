package errcode

import "fmt"

type Err struct {
	code int
	msg  string
}

func (e *Err) Error() string {
	return fmt.Sprintf("errorCode: %d, message: %s", e.code, e.msg)
}

func (e *Err) Code() int {
	return e.code
}
func (e *Err) Msg() string {
	return e.msg
}

func (e *Err) WrapErr(err error) *Err {
	return &Err{code: e.code, msg: fmt.Sprintf("%s: %s", e.msg, err.Error())}
}

func NewErr(code int, msg string) *Err {
	return &Err{code: code, msg: msg}
}
