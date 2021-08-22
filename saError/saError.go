package saError

import (
	"fmt"
	"gorm.io/gorm"
	"runtime"
	"strconv"
	"strings"
)

type Error struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Err    string `json:"err"`
	Caller string `json:"caller"`
}

// implements the error
func (e Error) Error() string {
	return e.String()
}

func (e Error) String() string {
	s := ""
	if e.Code > 0 {
		s = strings.Join([]string{strconv.Itoa(e.Code), e.Msg, e.Err, e.Caller}, " ")
	}
	return s
}

// err只接收字符串和error类型
func NewError(err interface{}) error {
	if err == nil {
		return nil
	}

	var e = Error{Code: NormalErrorCode, Msg: "", Err: "", Caller: ""}
	if s, ok := err.(string); ok {
		e.Msg = s
	} else if sae, ok := err.(Error); ok {
		return sae
	} else if sae, ok := err.(*Error); ok {
		return *sae
	} else if ev, ok := err.(error); ok {
		e.Err = ev.Error()
	}
	return e
}

/**
 * 会打印调用栈信息，建议只在类似controller最上层位置调用
 * @params err 可接收提示信息(err.Msg)和error类型
 * @params params 可接收错误信息(err.Err)和error类型
 * 注意：params参数会覆盖err参数
 */
func StackError(err interface{}, params ...interface{}) error {
	if err == nil {
		return nil
	}

	var resErr = Error{Code: NormalErrorCode, Msg: "", Err: "", Caller: "\n"}
	if s, ok := err.(string); ok {
		resErr.Msg = s
	} else {
		e, ok := err.(*Error)
		if ok == false {
			var e2 Error
			if e2, ok = err.(Error); ok {
				e = &e2
			}
		}

		if e != nil {
			resErr.Err = e.Err
			if len(e.Msg) > 0 {
				resErr.Msg = e.Msg
			}
			if e.Code > 0 {
				resErr.Code = e.Code
			}
			if e.Caller != "" {
				resErr.Caller += e.Caller + "\n"
			}
		} else if e, ok := err.(error); ok {
			resErr.Err = e.Error()
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
					resErr.Err = s
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
					resErr.Err = e.Err
					if len(e.Msg) > 0 {
						resErr.Msg = e.Msg
					}
					if e.Code > 0 {
						resErr.Code = e.Code
					}
					if e.Caller != "" {
						resErr.Caller += e.Caller + "\n"
					}
				} else if e, ok := err.(error); ok {
					resErr.Err = e.Error()
				}
			}
		}
	}

	//获取调用栈
	pc := make([]uintptr, 1)
	n := runtime.Callers(1, pc)
	if n >= 1 {
		f := runtime.FuncForPC(pc[0])
		file, line := f.FileLine(pc[0])
		resErr.Caller += fmt.Sprintf("%s:%d\n", file, line)
	}
	return &resErr
}

func DbErr(err error) bool {
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}
	return false
}
