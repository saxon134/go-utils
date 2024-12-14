package saGo

import (
	"fmt"
	"runtime/debug"
)

var isStopped = false

func Go(fn func()) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				fmt.Println(e)
				debug.PrintStack()
				return
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
