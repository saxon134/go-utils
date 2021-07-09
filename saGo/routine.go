package saGo

import (
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saLog"
	"time"
)

type Routine struct {
	routineMaxCnt    int              //goroutine最大并发数量
	routineMaxSecond int              //单个goroutine执行最大时间，如果非0，则会导致goroutine数量翻倍，暂不支持超时
	routineChan      chan struct{}    //控制goroutine并发数量
	paramsChan       chan interface{} //传输参数
	handle           func(params interface{})
}

/**
通过channel，分发事务，控制事务并发数量
*/
func NewRoutine(routineMaxCnt int, routineMaxSecond int, handle func(params interface{})) *Routine {
	if routineMaxCnt <= 0 {
		routineMaxCnt = 20
	}

	m := &Routine{
		routineMaxCnt:    routineMaxCnt,
		routineMaxSecond: routineMaxSecond,
		routineChan:      make(chan struct{}, routineMaxCnt),
		paramsChan:       make(chan interface{}, routineMaxCnt+1),
		handle:           handle,
	}
	return m
}

/**
任务数量大于最大协程数时会阻塞
*/
func (r *Routine) Do(params interface{}) {
	if r.routineMaxCnt == 0 {
		panic("请通过NewRoutine初始化")
		return
	}

	r.paramsChan <- params
	if len(r.routineChan) < r.routineMaxCnt {
		r.routineChan <- struct{}{}
		go func() {
			defer func() {
				_ = <-r.routineChan
			}()

			for {
				if v, ok := <-r.paramsChan; ok {
					//控制每个goroutine最大执行时间
					if r.routineMaxSecond > 0 {
						var quitTick = time.Tick(time.Second * time.Duration(r.routineMaxSecond))
						var handleDoneChan chan bool

						go func() {
							r.handle(v)
							handleDoneChan <- true
						}()

						select {
						case <-handleDoneChan:
							break
						case <-quitTick:
							str, _ := saData.DataToJson(v)
							saLog.Err("Goroutine执行超时：", str)
							break
						}
					} else {
						r.handle(v)
					}
				} else {
					return
				}
			}
		}()
	}
}
