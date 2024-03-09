package saLog

import (
	"fmt"
	"github.com/saxon134/go-utils/saData/saTime"
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

	logChan = make(chan string, 100)
	go func() {
		for {
			now := time.Now().Second()
			if now == lastLogTimestamp {
				if loggedCnt >= 10 {
					time.Sleep(time.Microsecond * 300)
				}
				loggedCnt++
			} else {
				lastLogTimestamp = now
				loggedCnt = 0
			}

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
	}
}

func SetRemoteUrl(url string) {
	remoteUrl = url
}

func Log(a ...interface{}) {
	if log == nil {
		return
	}

	//输出日志
	var s = ""
	for _, v := range a {
		s += fmt.Sprint(v) + " "
	}
	//logChan <- saTime.TimeToStr(time.Now(), saTime.FormatDefault) + " L " + s
	log.Log(saTime.TimeToStr(time.Now(), saTime.FormatDefault) + " L " + s)
}

func Err(a ...interface{}) {
	if log == nil {
		return
	}

	//输出日志
	var s = ""
	for _, v := range a {
		s += fmt.Sprint(v) + " "
	}

	//logChan <- saTime.TimeToStr(time.Now(), saTime.FormatDefault) + " E " + s
	log.Log(saTime.TimeToStr(time.Now(), saTime.FormatDefault) + " E " + s)
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

	//输出日志
	var s = ""
	for _, v := range a {
		s += fmt.Sprint(v) + " "
	}
	//logChan <- saTime.TimeToStr(time.Now(), saTime.FormatDefault) + " W " + s
	log.Log(saTime.TimeToStr(time.Now(), saTime.FormatDefault) + " W " + s)
}

func Info(a ...interface{}) {
	if logLevel > InfoLevel {
		return
	}

	//输出日志
	var s = ""
	for _, v := range a {
		s += fmt.Sprint(v) + " "
	}
	//logChan <- saTime.TimeToStr(time.Now(), saTime.FormatDefault) + " I " + s
	log.Log(saTime.TimeToStr(time.Now(), saTime.FormatDefault) + " I " + s)
}
