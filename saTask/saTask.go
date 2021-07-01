package saTask

import "fmt"

var _taskAry []Handle

type Handle struct {
	HandleFunc func()
	Name       string
	Spec       string // */2 * * * * * 每2秒执行一次
}

type handler struct {
	f func()
}

func (m handler) run() error {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("task panic: ", err)
			}
		}()

		m.f()
	}()
	return nil
}
