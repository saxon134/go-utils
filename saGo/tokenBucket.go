package saGo

import (
	"fmt"
	"math/rand/v2"
	"runtime/debug"
	"sync"
	"time"
)

type Bucket struct {
	fn func(*Bucket, interface{})
	wg *sync.WaitGroup

	Qps float32
	qpm int

	lock                   sync.Mutex
	lastMicrosecond        int64
	minIntervalMicrosecond int64
	second                 int64
	count                  int

	doneCnt int
	isDone  bool
	ch      chan interface{}

	slow       int
	expectTime int64 //预期耗时
	maxTime    int64 //最大耗时
	totalTime  int64 //总共耗时，计算平均耗时用
}

// size - 并发执行数量
// qps - 秒限制，必须大于0；0-1，如0.2，表示4秒允许执行一次
// fc - 执行接口
// 注意：size最节省资源的计算公式： size = qps * 每次执行的耗时
// 假如qps为20，执行耗时0.2秒，则size设置为4最节省资源
// size如果设置较大，不影响qps，只是浪费了些资源，如果执行比较费时可以通过加大size值改善执行速度
func NewBucket(size int, qps float32, fn func(b *Bucket, args interface{})) *Bucket {
	if size <= 0 || qps <= 0 {
		return nil
	}

	var b = &Bucket{
		fn:   fn,
		wg:   &sync.WaitGroup{},
		lock: sync.Mutex{},
	}

	if qps > 0 && qps < 1 {
		size = 1
	}

	//计算任务执行目标时间，实际低于目标时间，则表明任务执行较慢
	if qps > 0 {
		b.expectTime = int64(float32(size*1000)/qps + 1)
	}

	b.Qps = qps
	if qps >= 1 {
		b.minIntervalMicrosecond = int64(500000 / qps) // 1000000÷QPS×0.5
	} else {
		b.minIntervalMicrosecond = int64(1000000 / qps) // 1000000÷QPS
	}

	b.ch = make(chan interface{}, size+1)

	//消耗
	for i := 0; i < size; i++ {
		Go(func() {
			for {
				if args, ok := <-b.ch; ok {
					if b.fn != nil {
						var begin = time.Now().UnixMilli()
						b.fn(b, args)
						var diff = time.Now().UnixMilli() - begin
						b.totalTime += diff
						if diff > b.expectTime {
							b.slow++
						}
						if diff > b.maxTime {
							b.maxTime = diff
						}
					}
					b.wg.Done()
				} else {
					break
				}
			}
		})
	}
	return b
}

// 执行
func (b *Bucket) Invoke(args interface{}) {
	if e := recover(); e != nil {
		fmt.Println(e)
		debug.PrintStack()
		return
	}

	b.Consume()
	b.ch <- args
	b.doneCnt++
	b.wg.Add(1)
}

// 消耗，阻塞，消耗成功才能执行
func (b *Bucket) Consume() {
	for {
		if b.isDone {
			return
		}

		b.lock.Lock()
		var now = time.Now().UnixMicro()
		var second = now / 1000000
		if second != b.second {
			b.second = second
			b.count = 0
		}

		//限制最小执行间隔，为QPS间隔的60%
		var t = now - b.lastMicrosecond
		if t < b.minIntervalMicrosecond {
			b.lock.Unlock()

			//稍微错开点时间，减少lock竞争
			t = b.minIntervalMicrosecond - t
			if b.minIntervalMicrosecond > 1000 {
				var r = int64(float64(b.minIntervalMicrosecond) * 0.01)
				t += rand.Int64N(r)
			}
			time.Sleep(time.Duration(t+10) * time.Microsecond)
			continue
		} else {
			if b.Qps >= 1 && b.count >= int(b.Qps) {
				b.lock.Unlock()
				time.Sleep(time.Duration(b.minIntervalMicrosecond+10) * time.Microsecond)
				continue
			} else {
				b.lastMicrosecond = now
				b.count++
				b.lock.Unlock()
				return
			}
		}
	}
}

// 需要手动结束，会等待所有执行完再返回
func (b *Bucket) Done() {
	if b == nil {
		return
	}

	b.wg.Wait()
	close(b.ch)
	b.isDone = true
}

func (b *Bucket) Desc() string {
	var msg = ""
	if b.doneCnt <= 0 {
		msg = "待开始"
	} else {
		msg = fmt.Sprintf(
			"saGo 已完成：%d，慢任务:%d，平均执行时间：%d，最大执行时间：%d",
			b.doneCnt, b.slow, b.totalTime/int64(b.doneCnt), b.maxTime,
		)
	}
	return msg
}
