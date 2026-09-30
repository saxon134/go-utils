package saGo

import (
	"errors"
	"math"
	"sync"
	"time"
)

type limiterRegistry struct {
	mu        sync.Mutex
	entries   map[string]*limiter
	started   bool
	retention time.Duration
	stop      chan struct{}
	done      chan struct{}
}

var limiters = &limiterRegistry{entries: make(map[string]*limiter)}

// ConfigureLimiterCleanup 启用空闲 limiter 的自动回收，默认不启用。
// 必须在首次加锁调用前配置，且只允许配置一次；retention 必须大于 0。
// 启用后，所有 minSecond 必须是有限值，按现有毫秒精度计算的间隔不得超过 retention：
// LimiterTryLock 对超限参数返回 false，LimiterLock 则 panic，避免未加锁就继续执行。
// 持锁者、等待者和正在获取锁的调用均阻止回收。
func ConfigureLimiterCleanup(retention time.Duration) error {
	return limiters.configureCleanup(retention)
}

// LimiterForget 在资源生命周期结束时删除指定 key 的本地 limiter。
// 仅在没有持锁者、等待者或正在获取锁的调用时删除成功并返回 true。
// key 不存在或仍在使用时返回 false。删除会重置该 key 的冷却历史。
// 该方法不会删除 Redis 中的锁。
func LimiterForget(key string) bool {
	limiters.mu.Lock()
	defer limiters.mu.Unlock()

	lm := limiters.entries[key]
	if lm == nil || lm.refs != 0 {
		return false
	}
	delete(limiters.entries, key)
	return true
}

func (r *limiterRegistry) acquire(key string, minSecond float32, delayFirst bool) (*limiter, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.retention > 0 {
		seconds := float64(minSecond)
		// Match the lock's millisecond precision. Comparing raw float32 seconds
		// would incorrectly reject 0.1 against a 100ms retention window.
		milliseconds := math.Trunc(float64(minSecond * 1000))
		if math.IsNaN(seconds) || math.IsInf(seconds, 0) || milliseconds > float64(r.retention.Milliseconds()) {
			return nil, false
		}
	}
	r.started = true
	lm := r.entries[key]
	if lm == nil {
		lm = &limiter{}
		if delayFirst {
			lm.lastTime = time.Now().UnixMilli()
		}
		r.entries[key] = lm
	}
	// Register the reference before exposing the pointer or waiting on locker.
	lm.refs++
	lm.idleSince = time.Time{}
	return lm, true
}

func (r *limiterRegistry) lookup(key string) *limiter {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.entries[key]
}

func (r *limiterRegistry) release(lm *limiter) {
	r.mu.Lock()
	defer r.mu.Unlock()

	lm.refs--
	if lm.refs == 0 {
		lm.idleSince = time.Now()
	}
}

func (r *limiterRegistry) configureCleanup(retention time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if retention <= 0 {
		return errors.New("saGo: limiter cleanup retention must be positive")
	}
	if r.started || r.retention > 0 {
		return errors.New("saGo: limiter cleanup must be configured once, before the first lock call")
	}
	r.retention = retention
	r.stop = make(chan struct{})
	r.done = make(chan struct{})
	go r.cleanupLoop()
	return nil
}

func (r *limiterRegistry) cleanupLoop() {
	defer close(r.done)
	interval := r.retention / 2
	if interval < time.Millisecond {
		interval = time.Millisecond
	}
	if interval > time.Minute {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-r.stop:
			return
		case <-ticker.C:
			r.cleanup(time.Now())
		}
	}
}

func (r *limiterRegistry) cleanup(now time.Time) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.retention <= 0 {
		return 0
	}
	removed := 0
	for key, lm := range r.entries {
		if lm.refs == 0 && !lm.idleSince.IsZero() && now.Sub(lm.idleSince) >= r.retention {
			delete(r.entries, key)
			removed++
		}
	}
	return removed
}
