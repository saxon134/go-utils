package saGo

import (
	"context"
	"github.com/saxon134/go-utils/saLog"
	"time"
)

type Routine struct {
	routineMaxCnt  int              //goroutine最大并发数量
	routineMaxTime time.Duration    //所有任务执行最大时间
	routineChan    chan struct{}    //控制goroutine并发数量
	paramsChan     chan interface{} //传输参数
	handle         func(params interface{})
}

/**
通过channel，分发事务，控制事务并发数量
*/
func NewRoutine(routineMaxCnt int, routineMaxTime time.Duration, handle func(params interface{})) *Routine {
	if routineMaxCnt <= 0 {
		routineMaxCnt = 20
	}

	m := &Routine{
		routineMaxCnt:  routineMaxCnt,
		routineMaxTime: routineMaxTime,
		routineChan:    make(chan struct{}, routineMaxCnt),
		paramsChan:     make(chan interface{}, routineMaxCnt+1),
		handle:         handle,
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
				recover()
				_ = <-r.routineChan
			}()

			if r.routineMaxTime > 0 {
				ctx, _ := context.WithTimeout(context.Background(), r.routineMaxTime)
				for {
					v, ok := <-r.paramsChan
					if ok == false {
						return
					}

					select {
					case <-ctx.Done():
						saLog.Err("saGo routine time out...")
						return
					default:
						r.handle(v)
					}
				}
			} else {
				for {
					v, ok := <-r.paramsChan
					if ok == false {
						return
					}
					r.handle(v)
				}
			}
		}()
	}
}
