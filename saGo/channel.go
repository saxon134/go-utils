package saGo

import "fmt"

type Chan struct {
	Chan   chan struct{}
	name   string
	handle func(params ...interface{})
}

var _channel map[string]*Chan

/**
name: channel唯一标识
routineMaxCnt: 最大并发数量，默认1，小于100
handle: 注册处理方法，如果已经注册过，则不会覆盖之前的方法
*/
func Channel(name string, routineMaxCnt int, handle func(params ...interface{})) (c *Chan) {
	if name == "" {
		return
	}

	if routineMaxCnt < 1 {
		routineMaxCnt = 1
	}
	if routineMaxCnt > 100 {
		routineMaxCnt = 100
	}

	if _channel == nil || len(_channel) == 0 {
		_channel = map[string]*Chan{}
	}

	channel := _channel[name]
	if channel == nil {
		channel = &Chan{
			Chan:   make(chan struct{}, routineMaxCnt),
			name:   name,
			handle: handle,
		}

		_channel[name] = channel
	}
	return channel
}

func (m *Chan) Do(params ...interface{}) {
	if m == nil || m.Chan == nil {
		fmt.Println("未注册channel")
		return
	}

	m.Chan <- struct{}{}
	m.handle(params)
	<-m.Chan
}
