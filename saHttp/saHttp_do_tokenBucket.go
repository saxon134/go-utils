package saHttp

import (
	"sync"
	"sync/atomic"
	"time"
)

type TokenBucket struct {
	qps    int
	isDone bool
	wg     *sync.WaitGroup

	tokens int32        // 当前令牌数
	ticker *time.Ticker // 定时器
	second int32

	ch chan *BucketChannel
}

type BucketChannel struct {
	Params
	ResPtr interface{}
}

func NewTokenBucket(size int, qps int) *TokenBucket {
	if qps > 1000 {
		qps = 1000
	}
	if qps < 1 {
		qps = 1
	}

	if size <= 0 {
		return nil
	}

	var tb = &TokenBucket{
		wg:     &sync.WaitGroup{},
		isDone: false,
		tokens: 0,
		ticker: time.NewTicker(time.Second / time.Duration(qps)),
		ch:     make(chan *BucketChannel, qps),
	}
	tb.wg.Add(1)

	//生产
	go func() {
		for range tb.ticker.C {
			//1秒后复位
			var second = int32(time.Now().Second())
			if second != tb.second {
				atomic.StoreInt32(&tb.second, second)
				atomic.StoreInt32(&tb.tokens, 0)
			} else {
				atomic.AddInt32(&tb.tokens, 1)
			}

			if tb.isDone {
				time.Sleep(time.Second)
				tb.ticker.Stop()
				break
			}
		}
	}()

	//并发执行任务
	for i := 0; i < size; i++ {
		go func() {
			for {
				var args, unclosed = <-tb.ch
				if unclosed == false && len(tb.ch) == 0 && args == nil {
					tb.isDone = true
					tb.wg.Done()
					break
				}

				//最多重试100次
				for i := 0; i <= 100; i++ {
					//消耗
					tb.Consume(1)

					//请求
					var needRetry = true
					var err = _do(args.Params, args.ResPtr)
					if err != nil {
						needRetry = false
					} else {
						if args.Params.Retry == nil || args.Params.Retry(i+1, args.ResPtr, err) == false {
							needRetry = false
						} else
						//达到最大重试次数
						if i+1 >= 100 {
							needRetry = false
						}
					}

					//无需重试
					if needRetry == false {
						break
					}
				}
			}
		}()
	}
	return tb
}

func (tb *TokenBucket) Consume(num int32) {
	for {
		tokens := atomic.LoadInt32(&tb.tokens)
		if tokens < num {
			time.Sleep(time.Millisecond * 10)
			continue
		}

		atomic.StoreInt32(&tb.tokens, tokens-num)
		return
	}
}

// 执行完成，必须调用，否则会一直等待
func (tb *TokenBucket) Done() {
	if tb == nil || tb.isDone || tb.wg == nil {
		return
	}

	close(tb.ch)
	tb.wg.Wait()
}
