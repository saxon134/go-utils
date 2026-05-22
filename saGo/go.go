package saGo

import (
	"github.com/saxon134/go-utils/saLog"
	"runtime/debug"
)

var isStopped = false

func Go(fn func()) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				saLog.Err(e)
				saLog.Err(string(debug.Stack()))
			}
		}()
		fn()
	}()
}

func GoWithParams(args interface{}, fn func(params interface{})) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				saLog.Err(e)
				saLog.Err(string(debug.Stack()))
			}
		}()
		if fn != nil {
			fn(args)
		}
	}()
}

func IsStop() bool {
	return isStopped
}

func Stop() {
	isStopped = true
}
