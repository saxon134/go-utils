package saGo

import (
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saHit"
	"github.com/saxon134/go-utils/saRedis"
	"sync"
	"time"
)

type Bucket struct {
	locker sync.Mutex

	qps                    int
	qpm                    int
	minIntervalMillisecond int64
	lastTime               int64 //yyyymmddhhmmss、yyyymmddhhmm
	count                  int

	key   string
	redis *saRedis.Redis
}

func NewBucket(qps int, qpm int, args ...any) *Bucket {
	if qps <= 0 && qpm <= 0 {
		return nil
	}

	//取小的那个，只允许一个为0
	if qpm > 0 && qps > 0 {
		if qps*60 > qpm {
			qps = 0
		} else {
			qpm = 0
		}
	}

	var b = &Bucket{locker: sync.Mutex{}, qps: qps, qpm: qpm}
	var m1, m2 int
	if qps > 0 {
		m1 = 600 / qps // (1000/qps)*60%
	}
	if qpm > 0 {
		m2 = 36000 / qpm // (1000*60/qpm)*60%
	}
	b.minIntervalMillisecond = int64(saHit.Int(m1 > 0 && m1 < m2, m1, m2)) //最小间隔取小的值

	for _, v := range args {
		if redis, ok := v.(*saRedis.Redis); ok == true {
			b.redis = redis
		}

		if key, ok := v.(string); ok == true {
			b.key = key
		}
	}

	if b.redis == nil && _redis != nil {
		b.redis = _redis
	}

	if b.key == "" {
		b.key = "saGo:tokenBucket:" + saData.RandomStr()
	}
	return b
}

// 消耗，阻塞，消耗成功才能执行
func (b *Bucket) Consume() {
	for {
		b.locker.Lock()

		var now int64
		var limit = saHit.OrInt(b.qps, b.qpm)
		if b.qps > 0 {
			now = saData.Stoi64(time.Now().Format("20060102030405"))
		} else {
			now = saData.Stoi64(time.Now().Format("200601020304"))
		}

		if now == b.lastTime {
			if b.count >= limit {
				b.locker.Unlock()
				var r = int64(float64(b.minIntervalMillisecond)*0.01) + int64(b.count%5)
				time.Sleep(time.Duration(b.minIntervalMillisecond+r) * time.Millisecond)
				continue
			} else {
				b.count++
			}
		} else {
			b.lastTime = now
			b.count = 1
		}

		//本地拿到了令牌，再去读Redis令牌
		if b.redis != nil {
			var key = b.key + ":" + saData.String(now)
			var count, err = b.redis.GetInt64(key)
			if err == nil {
				if count >= int64(limit) {
					b.locker.Unlock()
					var r = int64(float64(b.minIntervalMillisecond)*0.01) + int64(b.count%5)
					time.Sleep(time.Duration(b.minIntervalMillisecond+r) * time.Millisecond)
					continue
				}
			}

			_, _ = b.redis.Do("INCR", key, 1, "EX", saHit.Int(b.qps > 0, 1, 60), "NX")
			b.locker.Unlock()
			break
		}

		b.locker.Unlock()
		break
	}
}
