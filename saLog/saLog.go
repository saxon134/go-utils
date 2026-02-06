package saLog

import (
	"fmt"
	"github.com/saxon134/go-utils/saData/saError"
	"github.com/saxon134/go-utils/saData/saTime"
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

func SetRemoteFun(f func([][4]string)) {
	if remoteFun != nil || f == nil {
		return
	}

	remoteFun = f
	go func() {
		var msgAry = make([][4]string, 0, 10)
		var lastCallAt int64
		for {
			var ary [4]string
			var ok bool
			ary, ok = <-remoteChan
			if ok == false {
				continue
			}

			msgAry = append(msgAry, ary)
			var timestamps = time.Now().Unix()
			if len(msgAry) >= 10 || timestamps-lastCallAt > 5 {
				remoteFun(msgAry)
				msgAry = make([][4]string, 0, 10)
			}
		}
	}()
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

	_log("W", a...)
}

func Info(a ...interface{}) {
	if logLevel > InfoLevel {
		return
	}

	_log("I", a...)
}

func _log(level string, a ...interface{}) {
	var now = time.Now()
	var timestamp = now.Unix()
	if timestamp <= lastLogTimestamp+1 {
		if loggedCnt >= 10 {
			if level != "E" {
				return
			}
		}
		loggedCnt++
	} else {
		lastLogTimestamp = timestamp
		loggedCnt = 0
	}

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
			file = formatCaller(file)
			if file != "" {
				caller += fmt.Sprintf("%s:%d => ", file, line)
			}
		}
		caller = strings.TrimSuffix(caller, "=> ")
	}

	for i, v := range a {
		if i == 0 && s != "" {
			continue
		}
		s += fmt.Sprint(v) + " "
	}

	var t = saTime.TimeToStr(now, saTime.FormatDefault)
	log.Log(t + " " + level + " " + caller + "\n" + s)

	//调用远端方法
	if remoteFun != nil {
		remoteChan <- [4]string{t, level, caller, s}
	}
}

func formatCaller(file string) string {
	if strings.Contains(file, "go/pkg/mod") || strings.Contains(file, "/go/src/") {
		return ""
	}

	var ignore = false
	for _, s := range ignores {
		if strings.Contains(file, s) {
			ignore = true
			break
		}
	}
	if ignore {
		return ""
	}

	if pkg != "" {
		var ary = strings.Split(file, pkg)
		if len(ary) == 2 {
			return ary[1]
		} else {
			return ""
		}
	}

	return file
}
