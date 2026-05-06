package saLog

import (
	"github.com/saxon134/go-utils/saOs"
	"io"
	golanglog "log"
	"os"
	"time"
)

type goLog struct {
	dir string
}

func (m *goLog) Log(a ...interface{}) {
	if m.dir != "" {
		var logpath = saOs.AbsPath(m.dir + "/" + time.Now().Format(time.DateOnly) + ".log")
		var info, _ = os.Stat(logpath)
		if info == nil || info.Name() == "" {
			file, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				golanglog.Fatal("日志文件打开失败:", err)
			}
			defer file.Close()

			multiWriter := io.MultiWriter(os.Stdout, file)
			golanglog.SetOutput(multiWriter)
		}
	}
	golanglog.Println(a...)
}

func initGoLog(dir string) *goLog {
	return &goLog{dir: dir}
}
