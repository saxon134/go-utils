package saError

import (
	"fmt"
	"gorm.io/gorm"
	"runtime"
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
		s = fmt.Sprintf("%d %s \n%s", e.Code, e.Msg, e.Caller)
	}
	return s
}

// err只接收字符串和error类型
func NewError(err interface{}) error {
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

// err只接收字符串和error类型
func NewSensitiveError(err interface{}) error {
	if err == nil {
		return nil
	}

	var e = Error{Code: SensitiveErrorCode, Msg: "", Caller: ""}
	if s, ok := err.(string); ok {
		e.Msg = s
	} else if sae, ok := err.(Error); ok {
		e = sae
		e.Code = SensitiveErrorCode
	} else if sae, ok := err.(*Error); ok {
		e= *sae
		e.Code = SensitiveErrorCode
	} else if ev, ok := err.(error); ok {
		e.Msg = ev.Error()
		e.Code = SensitiveErrorCode
	}
	return e
}

/**
 * @params err 可接收string和error类型
 * @params params 可接收int,string,error类型
 * 注意：params参数会覆盖err中相同类型数据
 * 常见用法：err传字符串，params空 -> code则是NormalErrorCode，msg为传的字符串
 * 常见用法：err传字符串，params传code值 -> 则是Error{code, msg}类型
 * 常见用法：err传err，其他为空 -> code则是SensitiveErrorCode，msg为err.error()
 */
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
					resErr.Msg = s
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
						resErr.Msg = e.Msg
					}
					if e.Code > 0 {
						resErr.Code = e.Code
					}
					if e.Caller != "" {
						resErr.Caller = e.Caller + "\n" + resErr.Caller
					}
				} else if e, ok := err.(error); ok {
					resErr.Msg = e.Error()
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
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}
	return false
}
