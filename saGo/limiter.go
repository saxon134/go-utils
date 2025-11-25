package saGo

import (
	"github.com/saxon134/go-utils/saData/saHit"
	"sync"
	"time"
)

type limiter struct {
	LastTime int64
	Lock     sync.Mutex
}

var limiterDIC = map[string]*limiter{}
var limiterLock = sync.Mutex{}

// 会阻塞
// milliSecond 2次执行最小间隔，毫秒
func LimiterLock(key string, milliSecond int) {
	if key == "" {
		return
	}

	//防止负数
	milliSecond = saHit.Int(milliSecond >= 0, milliSecond, 0)

	limiterLock.Lock()
	var lm = limiterDIC[key]
	if lm == nil {
		lm = &limiter{}
		limiterDIC[key] = lm
	}
	limiterLock.Unlock()

	lm.Lock.Lock()
	var now = time.Now().UnixMilli()
	var diff = int64(milliSecond) - (now - lm.LastTime)
	if diff > 0 {
		time.Sleep(time.Millisecond * time.Duration(diff))
	}
}

// 不阻塞
func LimiterTryLock(key string, milliSecond int) bool {
	if key == "" {
		return false
	}

	//防止负数
	milliSecond = saHit.Int(milliSecond >= 0, milliSecond, 0)

	limiterLock.Lock()
	var lm = limiterDIC[key]
	if lm == nil {
		lm = &limiter{}
		limiterDIC[key] = lm
	}
	limiterLock.Unlock()

	var ok = lm.Lock.TryLock()
	if ok == false {
		return false
	}

	var now = time.Now().UnixMilli()
	var diff = int64(milliSecond) - (now - lm.LastTime)
	if diff > 0 {
		time.Sleep(time.Millisecond * time.Duration(diff))
	}
	lm.LastTime = now
	return true
}

func LimiterUnLock(key string) {
	limiterLock.Lock()
	var lm = limiterDIC[key]
	if lm != nil {
		lm.LastTime = time.Now().UnixMilli()
		lm.Lock.Unlock()
		delete(limiterDIC, key)
	}
	limiterLock.Unlock()
}
