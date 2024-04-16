package saLog

var log _LOG
var logChan chan string
var logLevel LogLevel //当前日志等级
var remoteUrl string
var pkg string       //设置后只会打印指定项目日志（file路径包含pkg）
var ignores []string //忽略带有该字符的路径

var lastLogTimestamp int //最后打印日志的时间戳
var loggedCnt int        //当前日志打印次数，1秒清零

type _LOG interface {
	Log(a ...interface{})
}
