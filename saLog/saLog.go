package saLog

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"net/http"
	"strings"
	"time"
)

type LogType int8

const (
	NullType LogType = iota
	LocalType
	ZapType
)

type LogLevel int8

const (
	NullLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	endLevel
)

func Init(l LogLevel, t LogType) {
	if t == LocalType {
		log = initLocalLog()
		log.Log("local log初始化成功~")
	} else if t == ZapType {
		log = initZapLog()
		log.Log("zap log初始化成功~")
	}

	if log == nil {
		panic("log初始化失败~")
	}

	logLevel = l
	settedLogLevel = l
	logChan = make(chan string, 12)
	go func() {
		for {
			if s, ok := <-logChan; ok {
				//向远端发送日志
				if strings.HasPrefix(remoteUrl, "http") == true {
					_, _ = http.Post(remoteUrl, "text/plain", strings.NewReader(s))
				}
				log.Log(s)
			}
		}
	}()
}

func SetLogLevel(l LogLevel) {
	if l != NullLevel {
		logLevel = l
		settedLogLevel = l
	}
}

func SetRemoteUrl(url string) {
	remoteUrl = url
}

func Err(a ...interface{}) {
	if log == nil {
		return
	}

	//日志太多时，升等级
	if len(logChan) >= 5 {
		if logLevel == InfoLevel {
			logLevel = WarnLevel
		} else if logLevel == WarnLevel {
			logLevel = ErrorLevel
		}
		return
	}

	//日志少了之后，恢复等级
	if len(logChan) == 0 && logLevel > settedLogLevel {
		logLevel = settedLogLevel
	}

	logChan <- saData.TimeStr(time.Now(), saData.TimeFormat_Default) + " E\n" + fmt.Sprint(a...)
}

func Warn(a ...interface{}) {
	if log == nil {
		return
	}

	if logLevel > WarnLevel {
		return
	}

	//日志太多时，升等级
	if len(logChan) >= 5 {
		if logLevel == InfoLevel {
			logLevel = WarnLevel
		} else if logLevel == WarnLevel {
			logLevel = ErrorLevel
		}
		return
	}

	//日志少了之后，恢复等级
	if len(logChan) == 0 && logLevel > settedLogLevel {
		logLevel = settedLogLevel
	}

	//输出日志
	logChan <- saData.TimeStr(time.Now(), saData.TimeFormat_Default) + " W\n" + fmt.Sprint(a...)
	if len(logChan) >= 5 {
		if logLevel == InfoLevel {
			logLevel = WarnLevel
		} else if logLevel == WarnLevel {
			logLevel = ErrorLevel
		}
	}
}

func Info(a ...interface{}) {
	if logLevel > InfoLevel {
		return
	}

	//日志太多时，升等级
	if len(logChan) >= 5 {
		if logLevel == InfoLevel {
			logLevel = WarnLevel
		} else if logLevel == WarnLevel {
			logLevel = ErrorLevel
		}
		return
	}

	//日志少了之后，恢复等级
	if len(logChan) == 0 && logLevel > settedLogLevel {
		logLevel = settedLogLevel
	}

	//输出日志
	logChan <- saData.TimeStr(time.Now(), saData.TimeFormat_Default) + " W\n" + fmt.Sprint(a...)
	if len(logChan) >= 5 {
		if logLevel == InfoLevel {
			logLevel = WarnLevel
		} else if logLevel == WarnLevel {
			logLevel = ErrorLevel
		}
	}
}
