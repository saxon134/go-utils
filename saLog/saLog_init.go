package saLog

var log _LOG
var logChan chan string
var logLevel LogLevel       //当前日志等级
var settedLogLevel LogLevel //设置的日志等级，日志较少时需要恢复到设置的等级
var remoteUrl string

var logsPerSecond int    //1秒最多打印日志次数
var lastLogTimestamp int //最后打印日志的时间戳
var loggedCnt int        //当前日志打印次数，1秒清零

type _LOG interface {
	Log(a ...interface{})
}
