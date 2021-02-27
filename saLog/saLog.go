package saLog

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
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

	logChan = make(chan string, 12)
	go func() {
		for {
			if s, ok := <-logChan; ok {
				log.Log(s)
			}
		}
	}()
}

func Err(a ...interface{}) {
	if log == nil {
		return
	}

	if len(logChan) >= 10 {
		if logLevel == InfoLevel {
			logLevel = WarnLevel
		} else if logLevel == WarnLevel {
			logLevel = ErrorLevel
		}
		return
	}
	s := fmt.Sprint(a...)
	s = "E " + saData.TimeStr(time.Now(), saData.TimeFormat_Default) + " " + s
	logChan <- s
}

func Warn(a ...interface{}) {
	if log == nil {
		return
	}

	if logLevel <= WarnLevel {
		s := fmt.Sprint(a...)
		s = "W " + saData.TimeStr(time.Now(), saData.TimeFormat_Default) + " " + s
		logChan <- s

		if len(logChan) >= 10 {
			if logLevel == InfoLevel {
				logLevel = WarnLevel
			} else if logLevel == WarnLevel {
				logLevel = ErrorLevel
			}
			return
		}
	}
}

func Info(a ...interface{}) {
	if logLevel <= InfoLevel {
		s := fmt.Sprint(a...)
		s = "I " + saData.TimeStr(time.Now(), saData.TimeFormat_Default) + " " + s
		logChan <- s

		if len(logChan) >= 10 {
			if logLevel == InfoLevel {
				logLevel = WarnLevel
			}
			return
		}
	}
}

func Log(a ...interface{}) {
	if log != nil {
		s := fmt.Sprint(a...)
		s = "L " + saData.TimeStr(time.Now(), saData.TimeFormat_Default) + " " + s
		log.Log(s)
	}
}
