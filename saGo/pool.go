// QPS & 并发数量控制
package saGo

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

type Pool struct {
	fn func(interface{})
	wg *sync.WaitGroup

	total       int
	qps         int
	qpsLastTime int64
	qpm         int
	qpmLastTime int64

	done   int
	isDone bool
	ch     chan interface{}

	slow       int
	expectTime int64 //预期耗时
	maxTime    int64 //最大耗时
	totalTime  int64 //总共耗时，计算平均耗时用
}

// total - 总的任务数量
// size - 并发执行数量
// qps - 秒限制，0表示无限制
// qpm - 分钟限制，0表示无限制
// fc - 执行接口
// 注意：size最节省资源的计算公式： size = qps * 每次执行的耗时
// 假如qps为20，执行耗时0.2秒，则size设置为4最节省资源
// size如果设置较大，不影响qps，只是浪费了些资源，如果执行比较费时可以通过加大size值改善执行速度
func NewPool(total int, size int, qps int, qpm int, fn func(args interface{})) *Pool {
	if total <= 0 || size <= 0 {
		return nil
	}

	var p = &Pool{
		fn: fn,
		wg: &sync.WaitGroup{},
	}
	p.wg.Add(total)

	if size > total {
		size = total
	}

	//仅有QPM限制，需设置一下QPS，尽量均匀点，否则会导致短期QPS很高
	if qps == 0 && qpm > 0 {
		qps = qpm/60 + 1
	}

	//计算任务执行目标时间，实际低于目标时间，则表明任务执行较慢
	if qps > 0 {
		p.expectTime = int64(1000/qps + 1)
	}

	p.total = total
	p.qps = qps
	p.qpm = qpm
	p.ch = make(chan interface{}, size)

	for i := 0; i < size; i++ {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Println(e)
					debug.PrintStack()
					return
				}
			}()

			for {
				if args, ok := <-p.ch; ok {
					defer p.wg.Done()

					var begin = time.Now().UnixMilli()
					p.fn(args)
					var diff = time.Now().UnixMilli() - begin
					p.totalTime += diff
					if diff > p.expectTime {
						p.slow++
						if diff > p.maxTime {
							p.maxTime = diff
						}
					}
				} else {
					return
				}
			}
		}()
	}

	return p
}

// 执行
func (p *Pool) Invoke(args interface{}) {
	if e := recover(); e != nil {
		fmt.Println(e)
		debug.PrintStack()
		return
	}

	if p.isDone {
		return
	}

	if p.done == 0 {
		var t = time.Now().UnixMilli()
		p.qpsLastTime = t
		p.qpmLastTime = t
	}

	p.done++
	p.ch <- args

	if p.qps > 0 && p.done%p.qps == 0 {
		var t = time.Now().UnixMilli()
		var diff = t - p.qpsLastTime
		if diff < 1000 {
			time.Sleep(time.Millisecond * time.Duration(1005-diff))
		}
		p.qpsLastTime = time.Now().UnixMilli()
	}

	if p.qpm > 0 && p.done%p.qpm == 0 {
		var t = time.Now().UnixMilli()
		var diff = t - p.qpmLastTime
		if diff < 60000 {
			time.Sleep(time.Millisecond * time.Duration(60100-diff))
		}
		p.qpmLastTime = time.Now().UnixMilli()
	}

	if p.done >= p.total {
		close(p.ch)
		p.isDone = true
	}
}

// 等待所有执行完，也可以不调用
func (p *Pool) Wait() {
	if p == nil || p.total == 0 || p.wg == nil {
		return
	}
	p.wg.Wait()
}

func (p *Pool) Desc() string {
	if p.done == 2000 && 20*p.slow > p.done {

	}

	var msg = ""
	if p.done <= 0 {
		msg = "待开始"
	} else {
		msg = fmt.Sprintf(
			"saGo 总数：%d，已完成：%d，慢任务:%d，平均执行时间：%d，最大执行时间：%d",
			p.total, p.done, p.slow, p.totalTime/int64(p.done), p.maxTime,
		)
	}
	return msg
}
