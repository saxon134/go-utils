package saLog

import (
	"fmt"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saTime"
	"net/http"
	"runtime"
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

func SetPkg(pkgPath string, ignore ...string) {
	pkg = pkgPath
	ignores = ignore
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

	_log("L", a...)
}

func Err(a ...interface{}) {
	if log == nil {
		return
	}

	_log("E", a...)
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

	_log("W", a...)
}

func Info(a ...interface{}) {
	if logLevel > InfoLevel {
		return
	}

	_log("I", a...)
}

func _log(level string, a ...interface{}) {
	//输出日志
	var s = ""

	//获取调用栈
	var caller = ""
	if len(a) > 0 {
		if e, ok := a[0].(saError.Error); ok == true {
			if e.Caller != "" {
				caller = e.Caller
			}

			s = e.Msg
		}
	}
	if caller == "" {
		var pc = make([]uintptr, 10)
		var n = runtime.Callers(3, pc)
		for i := n - 1; i >= 0; i-- {
			var f = runtime.FuncForPC(pc[i])
			var file, line = f.FileLine(pc[i])
			if strings.Contains(file, "go/pkg/mod") || strings.Contains(file, "/go/src/") {
				continue
			}

			var ignore = false
			for _, s := range ignores {
				if strings.Contains(file, s) {
					ignore = true
					break
				}
			}
			if ignore {
				continue
			}

			if pkg != "" {
				var ary = strings.Split(file, pkg)
				if len(ary) == 2 {
					file = ary[1]
				} else {
					continue
				}
			}
			caller += fmt.Sprintf("%s:%d => ", file, line)
		}
		caller = strings.TrimSuffix(caller, "=> ")
	}

	for _, v := range a {
		s += fmt.Sprint(v) + " "
	}

	//logChan <- saTime.TimeToStr(time.Now(), saTime.FormatDefault) + " L " + s
	log.Log(saTime.TimeToStr(time.Now(), saTime.FormatDefault) + " " + level + " " + caller + "\n" + s)
}
