package saGo

import (
	"fmt"
	"github.com/saxon134/go-utils/saLog"
	"runtime/debug"
	"time"
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
		if e := recover(); e != nil {
			fmt.Println(e)
			debug.PrintStack()
			time.Sleep(time.Second)
			return
		}
		fn(args)
	}()
}

func IsStop() bool {
	return isStopped
}

func Stop() {
	isStopped = true
}
