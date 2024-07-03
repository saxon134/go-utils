package saLog

var log _LOG
var logLevel LogLevel //当前日志等级
var pkg string        //设置后只会打印指定项目日志（file路径包含pkg）
var ignores []string  //忽略带有该字符的路径

// 打印消息太多，则忽略等级低的
var lastLogTimestamp int64 //最后打印日志的时间戳
var loggedCnt int          //当前日志打印次数，1秒清零

// 远端方法，走channel，会保证不并发调用
var remoteFun func([][4]string)
var remoteChan = make(chan [4]string, 20)

type _LOG interface {
	Log(a ...interface{})
}
