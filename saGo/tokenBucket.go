package saGo

import (
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type Bucket struct {
	fn func(*Bucket, interface{})
	wg *sync.WaitGroup

	Qps int
	qpm int

	lock                sync.Mutex
	lastMillisecond     int64
	intervalMillisecond int64

	tokens int32        // 当前令牌数
	ticker *time.Ticker // 定时器
	second int32

	done   int
	isDone bool
	ch     chan interface{}

	slow       int
	expectTime int64 //预期耗时
	maxTime    int64 //最大耗时
	totalTime  int64 //总共耗时，计算平均耗时用
}

// size - 并发执行数量
// qps - 秒限制，0表示无限制
// fc - 执行接口
// 注意：size最节省资源的计算公式： size = qps * 每次执行的耗时
// 假如qps为20，执行耗时0.2秒，则size设置为4最节省资源
// size如果设置较大，不影响qps，只是浪费了些资源，如果执行比较费时可以通过加大size值改善执行速度
func NewBucket(size int, qps int, fn func(bucket *Bucket, args interface{})) *Bucket {
	var p = &Bucket{
		fn: fn,
		wg: &sync.WaitGroup{},

		isDone: false,
		tokens: 0,
		ticker: time.NewTicker(time.Second / time.Duration(qps)),

		lock: sync.Mutex{},
	}

	if size <= 0 {
		return nil
	}

	//计算任务执行目标时间，实际低于目标时间，则表明任务执行较慢
	if qps > 0 {
		p.expectTime = int64(1000/qps + 1)
	}

	p.Qps = qps
	p.intervalMillisecond = 1000 / int64(qps)
	p.ch = make(chan interface{}, size)

	//生产
	go func() {
		for range p.ticker.C {
			//1秒后复位
			var second = int32(time.Now().Second())
			if second != p.second {
				atomic.StoreInt32(&p.second, second)
				atomic.StoreInt32(&p.tokens, 1)
			} else {
				atomic.AddInt32(&p.tokens, 1)
			}
		}
	}()

	//消耗
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
					if p.fn != nil {
						var begin = time.Now().UnixMilli()
						p.fn(p, args)
						var diff = time.Now().UnixMilli() - begin
						p.totalTime += diff
						if diff > p.expectTime {
							p.slow++
						}
						if diff > p.maxTime {
							p.maxTime = diff
						}
					}
					p.wg.Done()
				} else {
					break
				}
			}
		}()
	}

	return p
}

// 执行
func (p *Bucket) Invoke(args interface{}) {
	if e := recover(); e != nil {
		fmt.Println(e)
		debug.PrintStack()
		return
	}

	if p.isDone {
		return
	}

	p.done++
	p.ch <- args
	p.wg.Add(1)
}

// 消耗，阻塞，消耗成功才能执行
func (b *Bucket) Consume() {
	b.lock.Lock()
	defer b.lock.Unlock()

	var now = time.Now().UnixMilli()
	var t = int64(now - b.lastMillisecond)
	if t < b.intervalMillisecond {
		time.Sleep(time.Duration(b.intervalMillisecond-t+5) * time.Millisecond)
	}

	for {
		if b.tokens <= 0 {
			time.Sleep(time.Millisecond * 100)
			continue
		}
		b.tokens--
		b.lastMillisecond = time.Now().UnixMilli()
		return
	}
}

// 需要手动结束，会等待所有执行完再返回
func (p *Bucket) Done() {
	if p == nil || p.isDone {
		return
	}

	close(p.ch)
	p.isDone = true
	if p.wg != nil {
		p.wg.Wait()
	}

	p.ticker.Stop()
}

func (p *Bucket) Desc() string {
	var msg = ""
	if p.done <= 0 {
		msg = "待开始"
	} else {
		msg = fmt.Sprintf(
			"saGo 已完成：%d，慢任务:%d，平均执行时间：%d，最大执行时间：%d",
			p.done, p.slow, p.totalTime/int64(p.done), p.maxTime,
		)
	}
	return msg
}
