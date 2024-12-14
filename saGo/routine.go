package saGo

import (
	"context"
	"fmt"
	"github.com/saxon134/go-utils/saLog"
	"runtime/debug"
	"time"
)

type Routine struct {
	routineMaxCnt  int              //goroutine最大并发数量
	routineMaxTime time.Duration    //所有任务执行最大时间
	paramsChan     chan interface{} //传输参数
	handle         func(params interface{})
}

// NewRoutine
// @Description: 通过channel，分发事务，控制事务并发数量
// @param routineMaxCnt 协程数量
// @param routineMaxTime 所有任务执行完最大总时间。单个任务协程无法被打断，时间控制没有意义。
// @param handle
// @return *Routine
func NewRoutine(routineMaxCnt int, routineMaxTime time.Duration, handle func(params interface{})) *Routine {
	if routineMaxCnt <= 0 {
		routineMaxCnt = 20
	}

	m := &Routine{
		routineMaxCnt:  routineMaxCnt,
		routineMaxTime: routineMaxTime,
		paramsChan:     make(chan interface{}, routineMaxCnt+1),
		handle:         handle,
	}

	if m.routineMaxCnt > 1 {
		for i := 0; i < m.routineMaxCnt; i++ {
			go func() {
				if e := recover(); e != nil {
					fmt.Println(e)
					debug.PrintStack()
					return
				}

				if m.routineMaxTime > 1 {
					ctx, cancelFunc := context.WithTimeout(context.Background(), m.routineMaxTime)
					for {
						v, ok := <-m.paramsChan
						if ok == false {
							cancelFunc()
							return
						}

						select {
						case <-ctx.Done():
							saLog.Err("saGo routine time out...")
							cancelFunc()
							return
						default:
							m.handle(v)
						}
					}
				} else {
					for {
						v, ok := <-m.paramsChan
						if ok == false {
							return
						}
						m.handle(v)
					}
				}
			}()
		}
	}

	return m
}

// Do
// @Description: 执行协程任务，任务数量大于最大协程数时会阻塞
// @receiver r
// @param params
func (r *Routine) Do(params interface{}) {
	if r.routineMaxCnt == 0 || r.paramsChan == nil {
		panic("请通过NewRoutine初始化")
		return
	}

	if r.routineMaxCnt > 1 {
		r.paramsChan <- params
	} else {
		r.handle(params)
	}
}
