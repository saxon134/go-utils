package saError

const (
	NormalErrorCode       = 3000
	UnAuthorizedErrorCode = 3100
)

var (
	ErrorDate         = Error{Code: NormalErrorCode, Msg: "数据有误"}
	ErrorParams       = Error{Code: NormalErrorCode, Msg: "参数有误"}
	ErrorUnauthorized = Error{Code: UnAuthorizedErrorCode, Msg: "未授权"}
)
