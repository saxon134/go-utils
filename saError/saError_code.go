package saError

const (
	NormalErrorCode       = 1000 //通常错误
	SensitiveErrorCode    = 1001 //此类错误信息一般不对外暴露，可能会有一些敏感数据
	TooFrequentErrCode    = 1010 //操作太频繁
	UnLoggedErrorCode     = 1103 //未登录
	UnauthorizedErrorCode = 1104 //未授权
	OutOfRange            = 1105 //超出范围
	ConflictErrorCode     = 1106 //有冲突
	BeDisplayedErrorCode  = 2000 //错误信息可以显示给用户
)

const (
	MissingParams     = "缺少参数"
	ErrorParams       = "参数有误"
	ErrorIo           = "接口错误"
	ErrorID           = "ID有误"
	ErrorDate         = "数据有误"
	ErrorPhone        = "手机格式有误"
	ErrorConfig       = "配置有误"
	ErrorUnauthorized = "未授权"
)
