package saError

const (
	BeDisplayedErrorCode = 2000 //一般可以展示给C端用户

	NormalErrorCode     = 3000 //一般可以展示给B端用户
	SensitiveErrorCode  = 3001 //错误信息包含敏感信息
	LoggedFailErrorCode = 3100 //登录失效
	UnLoggedErrorCode   = 3101 //未登录
	UnAuthedErrorCode   = 3102 //未授权
	OutOfRange          = 3105 //超出范围
	ConflictErrorCode   = 3106 //有冲突
)
