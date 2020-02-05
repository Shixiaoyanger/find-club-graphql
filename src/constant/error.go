package constant

import "errors"

const (
	ErrorMsgParamWrong  = "param 为空或错误"
	ErrorMsgTypeWrong   = "type wrong"
	ErrorLogin          = "用户名或密码错误"
	ErrDatebaseInternal = "数据库内部错误"
)

var (
	ErrorEmpty       = errors.New("empty error")
	ErrorParamsWrong = errors.New("params wrong")
)
