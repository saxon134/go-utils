package saTask

var _taskAry []Handle

type Handle struct {
	HandleFunc func()
	Name       string
	Spec       string // */2 * * * * * 每2秒执行一次
}

type handler struct {
	f func()
}

func (m *handler) run() error {
	go m.f()
	return nil
}
