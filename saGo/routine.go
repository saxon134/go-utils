package saGo

type Routine struct {
	maxGoroutine       int              //goroutine最大并发数量
	maxGoroutineSecond int              //单个goroutine执行最大时间，如果非0，则会导致goroutine数量翻倍，暂不支持超时
	routineChan        chan struct{}    //控制goroutine并发数量
	paramsChan         chan interface{} //传输参数
	handle             func(params interface{})
}

/**
通过channel，分发事务，控制事务并发数量
*/
func NewRoutine(maxGoroutine int, handle func(params interface{})) *Routine {
	if maxGoroutine <= 0 {
		maxGoroutine = 10
	}

	m := &Routine{
		maxGoroutine:       maxGoroutine,
		maxGoroutineSecond: 0,
		routineChan:        make(chan struct{}, maxGoroutine),
		paramsChan:         make(chan interface{}, maxGoroutine+1),
		handle:             handle,
	}
	return m
}

/**
任务数量大于最大协程数时会阻塞
*/
func (r *Routine) Do(params interface{}) {
	if r.maxGoroutine == 0 {
		panic("请通过NewRoutine初始化")
		return
	}

	r.paramsChan <- params
	if len(r.routineChan) < r.maxGoroutine {
		r.routineChan <- struct{}{}
		go func() {
			defer func() {
				_ = <-r.routineChan
			}()

			for {
				if v, ok := <-r.paramsChan; ok {
					r.handle(v)
				} else {
					return
				}
			}
		}()
	}
}
