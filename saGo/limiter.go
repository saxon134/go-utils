package saGo

import (
	"github.com/garyburd/redigo/redis"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saHit"
	"strings"
	"sync"
	"time"
)

type limiter struct {
	lastTime int64 //毫秒
	locker   sync.Mutex
}

var limiterDIC = map[string]*limiter{}
var limiterLocker = sync.Mutex{}

type LimiterOption string

const LimiterGlobalOption = LimiterOption("global")

// 会阻塞
// milliSecond 2次执行最小间隔（秒，可以是小数）
// maxMilliSecond  锁最大时间（秒），防止死锁
// 默认只本地锁
func LimiterLock(key string, minSecond float32, maxSecond float32, options ...any) {
	if key == "" {
		return
	}

	//防止负数
	var minMilSecond = int64(saHit.Float(minSecond >= 0, minSecond, 0) * 1000)
	var maxMilSecond = int64(saHit.Float(maxSecond >= 0, maxSecond, 0) * 1000)

	limiterLocker.Lock()
	var lm = limiterDIC[key]
	var now = time.Now().UnixMilli()
	var needLock = true
	if lm == nil {
		lm = &limiter{lastTime: now}
		limiterDIC[key] = lm
	} else
	//判断是否超出最大等待时间
	{
		if maxMilSecond > 0 && now-lm.lastTime >= maxMilSecond {
			lm.lastTime = now
			lm.locker.TryLock()
			needLock = false
		}
	}
	limiterLocker.Unlock()

	if needLock {
		lm.locker.Lock()
	}

	now = time.Now().UnixMilli()
	var diff = minMilSecond - (now - lm.lastTime)
	if diff > 0 {
		time.Sleep(time.Millisecond * time.Duration(diff))
	}

	//默认仅本地锁
	var isGlobal = false
	for _, v := range options {
		if opt, ok := v.(LimiterOption); ok && opt == LimiterGlobalOption {
			isGlobal = true
			break
		}
	}

	if _redis != nil && isGlobal {
		var redisKey = "saGo:limiter:" + key
		for {
			//最大10小时
			var expireSecond = saHit.Int64(maxMilSecond > 0, maxMilSecond/1000+1, 36000)
			var res, _ = redis.String(_redis.Do("SET", redisKey, now, "EX", expireSecond, "NX"))
			if strings.ToUpper(res) == "OK" {
				break
			}
			time.Sleep(time.Millisecond * time.Duration(saHit.OrInt64(minMilSecond, 100)))
		}
	}

	lm.lastTime = time.Now().UnixMilli()
}

// 不阻塞
func LimiterTryLock(key string, minSecond float32, options ...any) bool {
	if key == "" {
		return false
	}

	//防止负数
	var minMilliSecond = int64(saHit.Float(minSecond >= 0, minSecond, 0) * 1000)

	limiterLocker.Lock()
	var lm = limiterDIC[key]
	var now = time.Now().UnixMilli()
	if lm == nil {
		lm = &limiter{lastTime: now}
		limiterDIC[key] = lm
	}
	limiterLocker.Unlock()

	var ok = lm.locker.TryLock()
	if ok == false {
		return false
	}

	//默认仅本地锁
	var isGlobal = false
	var maxMilliSecond int64
	for _, v := range options {
		if opt, ok := v.(LimiterOption); ok && opt == LimiterGlobalOption {
			isGlobal = true
		} else {
			var maxSecond = saData.Float32(opt)
			if maxSecond > 0 {
				maxMilliSecond = int64(saHit.Float(maxSecond >= 0, maxSecond, 0) * 1000)
			}
		}
	}

	var redisKey = "saGo:limiter:" + key
	var lastTime = lm.lastTime
	if _redis != nil && isGlobal {
		var t, _ = redis.Int64(_redis.Do("GET", redisKey))
		if t > 0 && lm.lastTime < t {
			lastTime = t
		}
	}
	var diff = minMilliSecond - (now - lastTime)
	if diff > 0 {
		time.Sleep(time.Millisecond * time.Duration(diff))
	}
	lm.lastTime = lastTime

	//最大10分钟
	if _redis != nil && isGlobal {
		now = time.Now().UnixMilli()
		var expireSecond = saHit.Int64(maxMilliSecond > 0, maxMilliSecond/1000+1, 600)
		var res, _ = redis.String(_redis.Do("SET", redisKey, now, "EX", expireSecond, "NX"))
		if strings.ToUpper(res) == "OK" {
			lm.lastTime = now
			return true
		} else {
			lm.locker.Unlock()
			return false
		}
	}
	return true
}

// 解锁，不阻塞
func LimiterUnLock(key string) {
	limiterLocker.Lock()
	var lm = limiterDIC[key]
	if lm != nil {
		lm.lastTime = time.Now().UnixMilli()
		lm.locker.Unlock()
	}

	if _redis != nil {
		_, _ = _redis.Do("DEL", "saGo:limiter:"+key)
	}
	limiterLocker.Unlock()
}
