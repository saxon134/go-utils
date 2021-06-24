package saLog

var log _LOG
var logChan chan string
var logLevel LogLevel
var remoteUrl string

type _LOG interface {
	Log(a ...interface{})
}
