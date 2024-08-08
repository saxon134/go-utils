// QPS & 并发数量控制
package saGo

import (
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
}

// total - 总的任务数量
// size - 并发执行数量
// qps - 秒限制，0表示无限制
// qpm - 分钟限制，0表示无限制
// fc - 执行接口
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

	if qps > 0 && size > qps {
		size = qps
	}

	p.total = total
	p.qps = qps
	p.qpm = qpm
	p.ch = make(chan interface{}, size)

	for i := 0; i < size; i++ {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					return
				}
			}()

			for {
				if args, ok := <-p.ch; ok {
					p.wg.Done()
					p.fn(args)
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
		if t-p.qpsLastTime < 1000 {
			time.Sleep(time.Millisecond * time.Duration(1010-t+p.qpsLastTime))
		}
		p.qpsLastTime = time.Now().UnixMilli()
	}

	if p.qpm > 0 && p.done%p.qpm == 0 {
		var t = time.Now().UnixMilli()
		if t-p.qpmLastTime < 60000 {
			time.Sleep(time.Millisecond * time.Duration(60100-t+p.qpmLastTime))
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
	p.wg.Wait()
}
