package saError

const (
	SensitiveErrorCode = 1001 //错误信息包含敏感信息

	BeDisplayedErrorCode = 2000 //一般可以展示给C端用户
	NormalErrorCode      = 3000 //一般可以展示给B端用户

	LoggedFailErrorCode = 3100 //登录失效
	UnLoggedErrorCode   = 3101 //未登录
	UnAuthedErrorCode   = 3102 //未授权
	OutOfRange          = 3105 //超出范围
	ConflictErrorCode   = 3106 //有冲突
	ExistedErrorCode    = 3107 //已存在，重复
	NotExistedErrorCode = 3108 //不存在
)

var (
	ErrTooFrequent  = Error{Code: NormalErrorCode, Msg: "操作太频繁"}
	ErrParams       = Error{Code: NormalErrorCode, Msg: "缺少必要参数"}
	ErrExisted      = Error{Code: NormalErrorCode, Msg: "数据已存在"}
	ErrNotExisted   = Error{Code: NormalErrorCode, Msg: "数据不存在"}
	ErrData         = Error{Code: NormalErrorCode, Msg: "数据有误"}
	ErrNotSupport   = Error{Code: NormalErrorCode, Msg: "暂不支持"}
	ErrPassword     = Error{Code: NormalErrorCode, Msg: "账号、密码不匹配"}
	ErrLoggedFail   = Error{Code: NormalErrorCode, Msg: "登录失效"}
	ErrUnLogged     = Error{Code: NormalErrorCode, Msg: "未登录"}
	ErrUnauthorized = Error{Code: NormalErrorCode, Msg: "未授权"}
)
