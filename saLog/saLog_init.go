package saLog

var log _LOG
var logChan chan string
var logLevel LogLevel       //当前日志等级
var settedLogLevel LogLevel //设置的日志等级，日志较少时需要恢复到设置的等级
var remoteUrl string

type _LOG interface {
	Log(a ...interface{})
}
