package saError

import (
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
	_, file, line, ok := runtime.Caller(1)
	if ok {
		if ary := strings.Split(file, "/"); len(ary) > 0 {
			if len(ary) >= 2 {
				e.Caller = ary[len(ary)-2] + "/" + ary[len(ary)-1]
			} else if len(ary) >= 1 {
				e.Caller = ary[len(ary)-1]
			}
		}
		e.Caller += ":" + strconv.Itoa(line)
	}

	if s, ok := err.(string); ok {
		e.Err = s
		e.Msg = s
	} else if ev, ok := err.(error); ok {
		e.Msg = ""
		e.Err = ev.Error()
	}
	return e
}

// 会跟踪error的调用位置；
// err只接收字符串和error类型；字符串会覆盖msg
// params可传Code以及Msg，注意：会覆盖前面的
func StackError(err interface{}, params ...interface{}) error {
	if err == nil {
		return nil
	}

	var e = Error{Code: NormalErrorCode, Msg: "", Err: "", Caller: ""}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		if ary := strings.Split(file, "/"); len(ary) > 0 {
			if len(ary) >= 3 {
				e.Caller = ary[len(ary)-3] + "/" + ary[len(ary)-2] + "/" + ary[len(ary)-1]
			} else if len(ary) >= 2 {
				e.Caller = ary[len(ary)-2] + "/" + ary[len(ary)-1]
			} else if len(ary) >= 1 {
				e.Caller = ary[len(ary)-1]
			}
		}

		e.Caller += ":" + strconv.Itoa(line)
	}

	if s, ok := err.(string); ok {
		e.Err = s
		e.Msg = s
	} else {
		ev, ok := err.(*Error)
		if ok == false {
			var ev2 Error
			if ev2, ok = err.(Error); ok {
				ev = &ev2
			}
		}

		if ev != nil {
			e.Err = ev.Err
			if len(ev.Msg) > 0 {
				e.Msg = ev.Msg
			}
			if ev.Code > 0 {
				e.Code = ev.Code
			}
			if ev.Caller != "" {
				e.Caller = e.Caller + " " + ev.Caller
			}
		} else if e_v, ok := err.(error); ok {
			e.Err = e_v.Error()
		} else {
			return nil
		}
	}

	if params != nil {
		for _, v := range params {
			if code, ok := v.(int); ok {
				if code > 0 {
					e.Code = code
				}
			} else if s, ok := v.(string); ok {
				if s != "" {
					e.Msg = s
				}
			}
		}
	}

	return &e
}

func DbErr(err error) bool {
	if err != nil && err != gorm.ErrRecordNotFound {
		return true
	}
	return false
}
