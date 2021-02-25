package saLog

var log _LOG
var logChan chan string
var logLevel LogLevel

type _LOG interface {
	Log(a ...interface{})
}
