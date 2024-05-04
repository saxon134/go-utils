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
	"github.com/saxon134/go-utils/saData/saHit"
	"gorm.io/gorm"
	"net/url"
	"runtime"
	"strings"
)

var pkg = ""
var ignores []string

type Error struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Caller string `json:"caller"`
}

func SetPkg(pkgPath string, ignore ...string) {
	pkg = pkgPath
	ignores = ignore
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

// Stack
// @params err 可接收string和error类型
// @params params 可接收int,string,error类型
// 注意：params参数会覆盖err中相同类型数据
// 常见用法：err传字符串，params空 -> code则是NormalErrorCode，msg为传的字符串
// 常见用法：err传字符串，params传code值 -> 则是Error{code, msg}类型
// 常见用法：err传err，其他为空 -> code则是SensitiveErrorCode，msg为err.error()
func Stack(errs ...interface{}) error {
	if len(errs) == 0 {
		return nil
	}

	var resErr = Error{Code: NormalErrorCode, Msg: "", Caller: ""}
	if len(errs) == 1 && errs[0] == gorm.ErrRecordNotFound {
		return nil
	}

	for _, v := range errs {
		//字符串
		if s, ok := v.(string); ok {
			if s != "" {
				resErr.Msg += saHit.Str(resErr.Msg != "", " ", "") + s
			}
			continue
		}

		//saError
		{
			e, ok := v.(*Error)
			if ok == false {
				var e2 Error
				if e2, ok = v.(Error); ok {
					e = &e2
				}
			}
			if e != nil {
				if len(e.Msg) > 0 {
					resErr.Msg += saHit.Str(resErr.Msg != "", " ", "") + e.Msg
				}
				if e.Code > 0 && resErr.Code == 0 {
					resErr.Code = e.Code
				}
				if e.Caller != "" {
					resErr.Caller = e.Caller + "\n" + resErr.Caller
				}
				continue
			}
		}

		//url.Error
		if e, ok := v.(*url.Error); ok {
			resErr.Msg += saHit.Str(resErr.Msg != "", " ", "") + e.Err.Error()
			resErr.Caller = e.URL + "\n" + resErr.Caller
			continue
		}

		//error
		if e, ok := v.(error); ok {
			resErr.Msg += saHit.Str(resErr.Msg != "", " ", "") + e.Error()
			continue
		}

		//其他
		var str = fmt.Sprint(v)
		if str != "" {
			resErr.Msg += saHit.Str(resErr.Msg != "", "\n", "") + str
		}
	}

	//获取调用栈，已经存在caller，只获取一层；否则获取全部调用路径
	var pc = make([]uintptr, 10)
	var n = runtime.Callers(2, pc)
	if resErr.Caller != "" {
		if len(pc) > 0 {
			var f = runtime.FuncForPC(pc[0])
			var file, line = f.FileLine(pc[0])
			resErr.Caller += fmt.Sprintf(" => %s:%d", file, line)
		}
	} else {
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
			resErr.Caller += fmt.Sprintf("%s:%d => ", file, line)
		}
		resErr.Caller = strings.TrimSuffix(resErr.Caller, " => ")
	}

	return &resErr
}

func DbErr(err error) bool {
	if err != nil && err != gorm.ErrRecordNotFound && strings.Contains(err.Error(), "no row") == false {
		return true
	}
	return false
}
