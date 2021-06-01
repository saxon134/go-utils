package saTask

import (
	"github.com/astaxie/beego/toolbox"
	"github.com/pkg/errors"
)

func Init(handlers ...Handle) {
	if len(handlers) > 0 {
		for _, h := range handlers {
			toolbox.AddTask(h.Name, toolbox.NewTask(h.Name, h.Spec, handler{f: h.HandleFunc}.run))
		}
	}
	_taskAry = handlers

	toolbox.StartTask()
}

func Fire(name string) error {
	if name == "" {
		return errors.New("任务名称不能空")
	}

	if len(_taskAry) == 0 {
		return errors.New("任务空")
	}

	existed := false
	for _, t := range _taskAry {
		if t.Name == name {
			existed = true
			go t.HandleFunc()
			break
		}
	}
	if existed == false {
		return errors.New("未找到当前任务")
	}
	return nil
}
