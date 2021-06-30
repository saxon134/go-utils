package saLog

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"net/http"
	"runtime"
	"strconv"
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
	}
}

func SetRemoteUrl(url string) {
	remoteUrl = url
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

	s := "E " + saData.TimeStr(time.Now(), saData.TimeFormat_Default) + " "
	for _, v := range a {
		s += fmt.Sprint(v) + " "
	}

	caller := ""
	_, file, line, ok := runtime.Caller(1)
	if ok {
		if ary := strings.Split(file, "/"); len(ary) > 0 {
			if len(ary) >= 3 {
				caller = ary[len(ary)-3] + "/" + ary[len(ary)-2] + "/" + ary[len(ary)-1]
			} else if len(ary) >= 2 {
				caller = ary[len(ary)-2] + "/" + ary[len(ary)-1]
			} else if len(ary) >= 1 {
				caller = ary[len(ary)-1]
			}
		}
		caller += ":" + strconv.Itoa(line)
	}

	if len(caller) > 0 {
		logChan <- caller + " " + s
	} else {
		logChan <- s
	}
}

func Warn(a ...interface{}) {
	if log == nil {
		return
	}

	if logLevel <= WarnLevel {
		s := "W " + saData.TimeStr(time.Now(), saData.TimeFormat_Default) + " "
		for _, v := range a {
			s += fmt.Sprint(v) + " "
		}

		caller := ""
		_, file, line, ok := runtime.Caller(1)
		if ok {
			if ary := strings.Split(file, "/"); len(ary) > 0 {
				if len(ary) >= 3 {
					caller = ary[len(ary)-3] + "/" + ary[len(ary)-2] + "/" + ary[len(ary)-1]
				} else if len(ary) >= 2 {
					caller = ary[len(ary)-2] + "/" + ary[len(ary)-1]
				} else if len(ary) >= 1 {
					caller = ary[len(ary)-1]
				}
			}
			caller += ":" + strconv.Itoa(line)
		}

		if len(caller) > 0 {
			logChan <- caller + " " + s
		} else {
			logChan <- s
		}

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
		s := "I " + saData.TimeStr(time.Now(), saData.TimeFormat_Default) + " "
		for _, v := range a {
			s += fmt.Sprint(v) + " "
		}

		caller := ""
		_, file, line, ok := runtime.Caller(1)
		if ok {
			if ary := strings.Split(file, "/"); len(ary) > 0 {
				if len(ary) >= 3 {
					caller = ary[len(ary)-3] + "/" + ary[len(ary)-2] + "/" + ary[len(ary)-1]
				} else if len(ary) >= 2 {
					caller = ary[len(ary)-2] + "/" + ary[len(ary)-1]
				} else if len(ary) >= 1 {
					caller = ary[len(ary)-1]
				}
			}
			caller += ":" + strconv.Itoa(line)
		}

		if len(caller) > 0 {
			logChan <- caller + " " + s
		} else {
			logChan <- s
		}

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

		caller := ""
		_, file, line, ok := runtime.Caller(1)
		if ok {
			if ary := strings.Split(file, "/"); len(ary) > 0 {
				if len(ary) >= 3 {
					caller = ary[len(ary)-3] + "/" + ary[len(ary)-2] + "/" + ary[len(ary)-1]
				} else if len(ary) >= 2 {
					caller = ary[len(ary)-2] + "/" + ary[len(ary)-1]
				} else if len(ary) >= 1 {
					caller = ary[len(ary)-1]
				}
			}
			caller += ":" + strconv.Itoa(line)
		}

		if len(caller) > 0 {
			log.Log(caller + " " + s)
		} else {
			log.Log(s)
		}
	}
}
