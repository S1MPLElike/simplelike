package utils

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	CodeSuccess         = 0
	CodeParamError      = 1001
	CodeUserNotFound    = 1002
	CodePwdError        = 1003
	CodeNotLogin        = 1004
	CodeTokenInvalid    = 1005
	CodeOperationFailed = 1006
	CodeDBError         = 2001
	CodeRedisError      = 2002
	CodeServerError     = 5000
)

func Success(data interface{}) Response {
	return Response{
		Code: CodeSuccess,
		Msg:  "操作成功",
		Data: data,
	}
}

func Fail(code int, msg string) Response {
	return Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}
