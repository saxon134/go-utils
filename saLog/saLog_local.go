package saLog

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"os"
	"time"
)

type localLog struct {
	dir string
}

func (m *localLog) Log(a ...interface{}) {
	if m.dir != "" {
		var fullpath = saData.AbsPath(m.dir + "/" + time.Now().Format(time.DateOnly) + ".txt")
		var fileinfo, err = os.Stat(fullpath)
		if err != nil || fileinfo.Name() == "" {
			var file, _ = os.OpenFile(fullpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if file != nil {
				defer file.Close()
				os.Stdout = file
			}
		}
	}
	fmt.Println(a...)
}

func initLocalLog(path string) *localLog {
	var l = &localLog{}
	l.dir = path
	return l
}
