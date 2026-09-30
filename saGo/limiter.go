package saGo

import (
	"math"
	"strings"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saData/saHit"
	"github.com/saxon134/go-utils/saRedis"
)

type limiter struct {
	lastTime   int64 //毫秒
	locker     sync.Mutex
	redis      *saRedis.Redis
	redisValue string

	// Protected by limiterRegistry.mu, including while locker is held.
	refs      int
	idleSince time.Time
}

type LimiterOption string

const LimiterGlobalOption = LimiterOption("global")

// 会阻塞
// milliSecond 2次执行最小间隔（秒，可以是小数）
// maxMilliSecond  锁最大时间（秒），防止死锁
// 默认只本地锁
// 启用自动回收后，minSecond 超过保留窗口或不是有限值时 panic。
func LimiterLock(key string, minSecond float32, maxSecond float32, options ...any) {
	if key == "" {
		return
	}

	//防止负数
	var minMilSecond = int64(saHit.Float(minSecond >= 0, minSecond, 0) * 1000)

	lm, ok := limiters.acquire(key, minSecond, true)
	if !ok {
		panic("saGo: minSecond must be finite and must not exceed the limiter cleanup retention")
	}

	lm.locker.Lock()

	now := time.Now().UnixMilli()
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
		var redisValue = strings.Replace(uuid.NewV4().String(), "-", "", -1)
		expireSecond := int64(600)
		if seconds, ok := limiterLeaseSeconds(maxSecond); ok {
			expireSecond = seconds
		}
		for {
			//默认10分钟
			var res, _ = redis.String(_redis.Do("SET", redisKey, redisValue, "EX", expireSecond, "NX"))
			if strings.ToUpper(res) == "OK" {
				lm.redis = _redis
				lm.redisValue = redisValue
				break
			}
			time.Sleep(time.Millisecond * time.Duration(saHit.OrInt64(minMilSecond, 100)))
		}
	}

	lm.lastTime = time.Now().UnixMilli()
}

// 不阻塞
// 启用自动回收后，minSecond 超过保留窗口或不是有限值时返回 false。
func LimiterTryLock(key string, minSecond float32, options ...any) bool {
	if key == "" {
		return false
	}

	//防止负数
	var minMilliSecond = int64(saHit.Float(minSecond >= 0, minSecond, 0) * 1000)

	//默认仅本地锁
	var isGlobal = false
	expireSecond := int64(600)
	for _, v := range options {
		if opt, ok := v.(LimiterOption); ok {
			isGlobal = isGlobal || opt == LimiterGlobalOption
		} else if seconds, ok := limiterLeaseSeconds(v); ok {
			expireSecond = seconds
		}
	}

	lm, ok := limiters.acquire(key, minSecond, false)
	if !ok {
		return false
	}
	if !lm.locker.TryLock() {
		limiters.release(lm)
		return false
	}
	acquired := false
	defer func() {
		if !acquired {
			lm.locker.Unlock()
			limiters.release(lm)
		}
	}()

	now := time.Now().UnixMilli()
	var diff = minMilliSecond - (now - lm.lastTime)
	if diff > 0 {
		return false
	}

	//默认10分钟
	if _redis != nil && isGlobal {
		now = time.Now().UnixMilli()
		var redisKey = "saGo:limiter:" + key
		var redisValue = strings.Replace(uuid.NewV4().String(), "-", "", -1)
		var res, _ = redis.String(_redis.Do("SET", redisKey, redisValue, "EX", expireSecond, "NX"))
		if strings.ToUpper(res) == "OK" {
			lm.lastTime = now
			lm.redis = _redis
			lm.redisValue = redisValue
			acquired = true
			return true
		} else {
			return false
		}
	}
	acquired = true
	return true
}

// 解锁，不阻塞
func LimiterUnLock(key string) {
	lm := limiters.lookup(key)
	if lm == nil {
		return
	}

	var redisConn = lm.redis
	var redisValue = lm.redisValue
	lm.redis = nil
	lm.redisValue = ""

	if redisConn != nil && redisValue != "" {
		_, _ = redisConn.Do("EVAL", limiterRedisUnlockScript, 1, "saGo:limiter:"+key, redisValue)
	}

	lm.lastTime = time.Now().UnixMilli()
	lm.locker.Unlock()
	limiters.release(lm)
}

// Return whole seconds for Redis EX without losing sub-millisecond leases.
// Invalid options are ignored. Bound leases to the range of time.Duration.
func limiterLeaseSeconds(value any) (int64, bool) {
	seconds, err := saData.ToFloat64(value)
	if err != nil || math.IsNaN(seconds) || seconds <= 0 ||
		seconds > float64((1<<63-1)/int64(time.Second)) {
		return 0, false
	}
	return int64(math.Ceil(seconds)), true
}

const limiterRedisUnlockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`
