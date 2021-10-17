package saLog

import (
	"github.com/beego/beego/v2/core/logs"
)

type beegoLog struct {
}

func initBeegoLog() *beegoLog {
	//_ = logs.SetLogger("file")
	//_ = logs.SetLogger(logs.AdapterFile, `{"filename":"app.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	logs.SetLogFuncCall(true)
	logs.SetLevel(logs.LevelDebug)
	logs.SetPrefix("")
	logs.SetLogFuncCallDepth(3)
	return &beegoLog{}
}

func (m *beegoLog) Log(a ...interface{}) {
	logs.Debug(a)
}
