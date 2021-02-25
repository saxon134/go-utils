package saTask

import (
	"github.com/astaxie/beego/toolbox"
)

type Handle struct {
	HandleFunc func() error
	Name       string
	Spec       string // */2 * * * * * 每2秒执行一次
}

type handler struct {
	f func() error
}

func (m *handler) run() error {
	go m.f()
	return nil
}

func Init(handlers ...Handle) {
	if len(handlers) > 0 {
		for _, h := range handlers {
			toolbox.AddTask(h.Name, toolbox.NewTask(h.Name, h.Spec, handler{f: h.HandleFunc}.run))
		}
	}

	toolbox.StartTask()
}
