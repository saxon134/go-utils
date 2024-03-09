package saLog

var log _LOG
var logChan chan string
var logLevel LogLevel       //当前日志等级
var remoteUrl string

var lastLogTimestamp int //最后打印日志的时间戳
var loggedCnt int        //当前日志打印次数，1秒清零

type _LOG interface {
	Log(a ...interface{})
}
