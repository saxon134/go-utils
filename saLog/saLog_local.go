package saLog

import "fmt"

type localLog struct {
}

func (m *localLog) Log(a ...interface{}) {
	fmt.Println(a)
}

func initLocalLog() *localLog {
	return &localLog{}
}
