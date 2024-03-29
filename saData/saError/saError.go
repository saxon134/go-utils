/*
Package saError
一般建议：
各项目自行定义code值及对应的返回给前端的错误信息
业务代码内不固定错误文案
Msg是错误信息，只打印到日志
*/
package saError

import (
	"fmt"
	"gorm.io/gorm"
	"runtime"
	"strings"
)

type Error struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Caller string `json:"caller"`
}

// implements the error
func (e Error) Error() string {
	return e.String()
}

func (e Error) String() string {
	s := ""
	if e.Code > 0 {
		s = fmt.Sprintf("%d %s", e.Code, e.Msg)
		if e.Caller != "" {
			s += "\n " + e.Caller
		}
	}
	return s
}

func Msg(e error) string {
	if e == nil {
		return ""
	}

	if ee, ok := e.(Error); ok {
		return ee.Msg
	} else if ee, ok := e.(*Error); ok {
		return ee.Msg
	} else if ee, ok := e.(error); ok {
		return ee.Error()
	}
	return ""
}

// New 只接收字符串和error类型
func New(err interface{}) error {
	if err == nil {
		return nil
	}

	var e = Error{Code: NormalErrorCode, Msg: "", Caller: ""}
	if s, ok := err.(string); ok {
		e.Msg = s
	} else if sae, ok := err.(Error); ok {
		return sae
	} else if sae, ok := err.(*Error); ok {
		return *sae
	} else if ev, ok := err.(error); ok {
		e.Msg = ev.Error()
		e.Code = SensitiveErrorCode
	}
	return e
}

func NewBeDisplayedError(err string) error {
	return Error{
		Code: BeDisplayedErrorCode,
		Msg:  err,
	}
}

// NewSensitiveError err只接收字符串和error类型
func NewSensitiveError(err interface{}) error {
	if err == nil {
		return nil
	}

	var e = Error{Code: 0, Msg: "", Caller: ""}
	if s, ok := err.(string); ok {
		e.Msg = s
	} else if sae, ok := err.(Error); ok {
		e = sae
	} else if sae, ok := err.(*Error); ok {
		e = *sae
	} else if ev, ok := err.(error); ok {
		e.Msg = ev.Error()
	}

	e.Code = SensitiveErrorCode
	return e
}

// StackError
// @params err 可接收string和error类型
// @params params 可接收int,string,error类型
// 注意：params参数会覆盖err中相同类型数据
// 常见用法：err传字符串，params空 -> code则是NormalErrorCode，msg为传的字符串
// 常见用法：err传字符串，params传code值 -> 则是Error{code, msg}类型
// 常见用法：err传err，其他为空 -> code则是SensitiveErrorCode，msg为err.error()
func StackError(err interface{}, params ...interface{}) error {
	if err == nil {
		return nil
	}

	var resErr = Error{Code: NormalErrorCode, Msg: "", Caller: ""}
	if s, ok := err.(string); ok {
		resErr.Msg = s
		resErr.Code = NormalErrorCode
	} else {
		e, ok := err.(*Error)
		if ok == false {
			var e2 Error
			if e2, ok = err.(Error); ok {
				e = &e2
			}
		}

		if e != nil {
			if len(e.Msg) > 0 {
				resErr.Msg = e.Msg
			}
			if e.Code > 0 {
				resErr.Code = e.Code
			}
			if e.Caller != "" {
				if resErr.Caller == "" {
					resErr.Caller = e.Caller
				} else {
					resErr.Caller = e.Caller + "\n" + resErr.Caller
				}
			}
		} else if e, ok := err.(error); ok {
			resErr.Msg = e.Error()
			resErr.Code = SensitiveErrorCode
		} else {
			return nil
		}
	}

	if params != nil {
		for _, v := range params {
			if code, ok := v.(int); ok {
				if code > 0 {
					resErr.Code = code
				}
			} else if s, ok := v.(string); ok {
				if s != "" {
					resErr.Msg += s
				}
			} else {
				e, ok := err.(*Error)
				if ok == false {
					var e2 Error
					if e2, ok = err.(Error); ok {
						e = &e2
					}
				}

				if e != nil {
					if len(e.Msg) > 0 {
						resErr.Msg += s
					}
					if e.Code > 0 {
						resErr.Code = e.Code
					}
					if e.Caller != "" {
						resErr.Caller = e.Caller + "\n" + resErr.Caller
					}
				} else if e, ok := err.(error); ok {
					resErr.Msg += e.Error()
				}
			}
		}
	}

	//获取调用栈
	pc := make([]uintptr, 1)
	n := runtime.Callers(2, pc)
	if n >= 1 {
		f := runtime.FuncForPC(pc[0])
		file, line := f.FileLine(pc[0])
		resErr.Caller = fmt.Sprintf("%s:%d\n", file, line) + resErr.Caller
	}
	return &resErr
}

func DbErr(err error) bool {
	if err != nil && err != gorm.ErrRecordNotFound && strings.Contains(err.Error(), "no row") == false {
		return true
	}
	return false
}
