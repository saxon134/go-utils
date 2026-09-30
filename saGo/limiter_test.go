package saGo

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func isolateLimiters(t *testing.T) *limiterRegistry {
	t.Helper()
	oldLimiters, oldRedis := limiters, _redis
	r := &limiterRegistry{entries: make(map[string]*limiter)}
	limiters, _redis = r, nil
	t.Cleanup(func() {
		if r.stop != nil {
			close(r.stop)
			<-r.done
		}
		limiters, _redis = oldLimiters, oldRedis
	})
	return r
}

func TestLimiterTryLockRedisLeaseOptions(t *testing.T) {
	tests := []struct {
		name    string
		options []any
		want    string
	}{
		{"default", nil, "600"},
		{"int", []any{3}, "3"},
		{"int64", []any{int64(3)}, "3"},
		{"float32", []any{float32(3)}, "3"},
		{"float64", []any{float64(3)}, "3"},
		{"string", []any{"3"}, "3"},
		{"fraction", []any{3.2}, "4"},
		{"half-second", []any{0.5}, "1"},
		{"sub-millisecond", []any{0.0001}, "1"},
		{"long", []any{1800}, "1800"},
		{"last-valid", []any{3, 7, "invalid", -1}, "7"},
		{"unknown-option", []any{LimiterOption("3")}, "600"},
		{"zero", []any{0}, "600"},
		{"negative", []any{-1}, "600"},
		{"nan", []any{math.NaN()}, "600"},
		{"infinite", []any{math.Inf(1)}, "600"},
		{"overflow", []any{1e30}, "600"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolateLimiters(t)
			store := newLimiterRedisStore()
			_redis = newLimiterTestRedis(store)
			t.Cleanup(func() { _redis.Pool.Close() })
			options := append(append([]any{}, tt.options...), LimiterGlobalOption)
			if !LimiterTryLock("lease", 0, options...) {
				t.Fatal("failed to acquire lock")
			}
			defer LimiterUnLock("lease")
			if got := store.Expire("saGo:limiter:lease"); got != tt.want {
				t.Fatalf("Redis EX = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLimiterLockRedisLeaseRounding(t *testing.T) {
	tests := []struct {
		seconds float32
		want    string
	}{
		{0, "600"},
		{3, "3"},
		{3.2, "4"},
		{0.5, "1"},
		{0.0001, "1"},
		{1800, "1800"},
		{float32(math.Inf(1)), "600"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.seconds), func(t *testing.T) {
			isolateLimiters(t)
			store := newLimiterRedisStore()
			_redis = newLimiterTestRedis(store)
			t.Cleanup(func() { _redis.Pool.Close() })
			LimiterLock("lease", 0, tt.seconds, LimiterGlobalOption)
			defer LimiterUnLock("lease")
			if got := store.Expire("saGo:limiter:lease"); got != tt.want {
				t.Fatalf("Redis EX = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLimiterForgetOnlyIdleEntriesAndResetsCooldown(t *testing.T) {
	isolateLimiters(t)
	if LimiterForget("missing") {
		t.Fatal("forgot a missing key")
	}
	if !LimiterTryLock("resource", 60) {
		t.Fatal("failed to acquire first lock")
	}
	if LimiterForget("resource") {
		t.Fatal("forgot a held lock")
	}
	if LimiterTryLock("resource", 0) {
		t.Fatal("acquired a held lock")
	}
	LimiterUnLock("resource")
	if LimiterTryLock("resource", 60) {
		LimiterUnLock("resource")
		t.Fatal("ignored cooldown")
	}
	if !LimiterForget("resource") {
		t.Fatal("failed acquisitions leaked references")
	}
	if !LimiterTryLock("resource", 60) {
		t.Fatal("forget did not reset cooldown")
	}
	LimiterUnLock("resource")
}

func TestLimiterForgetAfterRedisContention(t *testing.T) {
	isolateLimiters(t)
	store := newLimiterRedisStore()
	_redis = newLimiterTestRedis(store)
	t.Cleanup(func() { _redis.Pool.Close() })
	store.Set("saGo:limiter:resource", "other-owner")
	if LimiterTryLock("resource", 0, LimiterGlobalOption) {
		t.Fatal("acquired another owner's lock")
	}
	if !LimiterForget("resource") {
		t.Fatal("Redis failure leaked a reference")
	}
	if got := store.Get("saGo:limiter:resource"); got != "other-owner" {
		t.Fatal("forget modified the Redis lock")
	}
}

func TestLimiterForgetAfterRedisConnectionError(t *testing.T) {
	isolateLimiters(t)
	_redis = newLimiterTestRedis(newLimiterRedisStore())
	_redis.Pool.Close()
	if LimiterTryLock("resource", 0, LimiterGlobalOption) {
		t.Fatal("acquired a lock through a closed Redis pool")
	}
	if !LimiterForget("resource") {
		t.Fatal("Redis connection error leaked a reference")
	}
}

func waitForLimiterCondition(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for limiter state")
		}
		time.Sleep(time.Millisecond)
	}
}

func TestLimiterCleanupAndForgetKeepWaiters(t *testing.T) {
	r := isolateLimiters(t)
	if err := ConfigureLimiterCleanup(time.Hour); err != nil {
		t.Fatal(err)
	}
	LimiterLock("resource", 0, 0)
	entered, release, done := make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		defer close(done)
		LimiterLock("resource", 0, 0)
		close(entered)
		<-release
		LimiterUnLock("resource")
	}()
	ownerReleased := false
	t.Cleanup(func() {
		if !ownerReleased {
			LimiterUnLock("resource")
		}
		close(release)
		<-done
	})

	waitForLimiterCondition(t, func() bool {
		r.mu.Lock()
		defer r.mu.Unlock()
		return r.entries["resource"].refs == 2
	})
	if LimiterForget("resource") || r.cleanup(time.Now().Add(2*time.Hour)) != 0 {
		t.Fatal("removed a lock with a registered waiter")
	}
	LimiterUnLock("resource")
	ownerReleased = true
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("waiter failed to acquire lock")
	}
	if LimiterForget("resource") || LimiterTryLock("resource", 0) {
		t.Fatal("waiter lost exclusive ownership")
	}
}

func TestLimiterCleanupKeepsCooldownAndReclaimsDynamicKeys(t *testing.T) {
	r := isolateLimiters(t)
	if err := ConfigureLimiterCleanup(time.Hour); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("resource-%d", i)
		if !LimiterTryLock(key, 3600) {
			t.Fatal("failed to acquire first lock")
		}
		LimiterUnLock(key)
	}
	if got := r.cleanup(time.Now().Add(30 * time.Minute)); got != 0 {
		t.Fatalf("removed %d entries inside their retention window", got)
	}
	if LimiterTryLock("resource-0", 3600) {
		LimiterUnLock("resource-0")
		t.Fatal("cleanup bypassed cooldown")
	}
	if got := r.cleanup(time.Now().Add(2 * time.Hour)); got != 1000 {
		t.Fatalf("removed %d expired entries, want 1000", got)
	}
}

func TestLimiterCleanupRunsInBackground(t *testing.T) {
	r := isolateLimiters(t)
	if err := ConfigureLimiterCleanup(10 * time.Millisecond); err != nil {
		t.Fatal(err)
	}
	LimiterLock("resource", 0, 0)
	LimiterUnLock("resource")
	waitForLimiterCondition(t, func() bool { return r.lookup("resource") == nil })
}

func TestLimiterCleanupConfigurationAndIntervalBound(t *testing.T) {
	r := isolateLimiters(t)
	for _, retention := range []time.Duration{0, -time.Second} {
		if ConfigureLimiterCleanup(retention) == nil {
			t.Fatal("accepted nonpositive retention")
		}
	}
	if err := ConfigureLimiterCleanup(time.Second); err != nil {
		t.Fatal(err)
	}
	if ConfigureLimiterCleanup(2*time.Second) == nil {
		t.Fatal("allowed changing the interval bound after cleanup was configured")
	}
	for _, seconds := range []float32{2, float32(math.NaN()), float32(math.Inf(1))} {
		if LimiterTryLock("invalid", seconds) {
			LimiterUnLock("invalid")
			t.Fatalf("accepted minSecond=%v", seconds)
		}
	}
	func() {
		defer func() {
			if recover() == nil {
				t.Error("blocking Lock must not return without enforcing its interval bound")
			}
		}()
		LimiterLock("invalid", 2, 0)
		LimiterUnLock("invalid")
	}()
	if r.lookup("invalid") != nil {
		t.Fatal("invalid intervals created a registry entry")
	}
	if !LimiterTryLock("boundary", 1) {
		t.Fatal("rejected interval equal to retention")
	}
	LimiterUnLock("boundary")
}

func TestLimiterCleanupCannotEnableAfterUse(t *testing.T) {
	r := isolateLimiters(t)
	if !LimiterTryLock("resource", 86400) {
		t.Fatal("default mode must allow intervals without a retention bound")
	}
	LimiterUnLock("resource")
	if got := r.cleanup(time.Now().Add(48 * time.Hour)); got != 0 {
		t.Fatal("cleanup ran without configuration")
	}
	if !LimiterForget("resource") {
		t.Fatal("failed to forget idle resource")
	}
	if ConfigureLimiterCleanup(time.Second) == nil {
		t.Fatal("allowed enabling cleanup after use")
	}
}

func TestLimiterCleanupFractionalIntervalBound(t *testing.T) {
	tests := []struct {
		retention time.Duration
		seconds   float32
		want      bool
	}{
		{100 * time.Millisecond, 0.1, true},
		{150 * time.Millisecond, 0.15, true},
		{100 * time.Millisecond, 0.101, false},
		{time.Second, math.MaxFloat32, false},
		{time.Second, -math.MaxFloat32, true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s/%v", tt.retention, tt.seconds), func(t *testing.T) {
			isolateLimiters(t)
			if err := ConfigureLimiterCleanup(tt.retention); err != nil {
				t.Fatal(err)
			}
			got := LimiterTryLock("resource", tt.seconds)
			if got {
				LimiterUnLock("resource")
			}
			if got != tt.want {
				t.Fatalf("TryLock(%v) = %v, want %v", tt.seconds, got, tt.want)
			}
		})
	}
}

func TestLimiterConcurrentCleanupPreservesMutualExclusion(t *testing.T) {
	r := isolateLimiters(t)
	if err := ConfigureLimiterCleanup(time.Nanosecond); err != nil {
		t.Fatal(err)
	}
	const keyCount = 4
	var active [keyCount]atomic.Int32
	var violations atomic.Int32
	stop, done := make(chan struct{}), make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-stop:
				return
			default:
				r.cleanup(time.Now())
				for i := 0; i < keyCount; i++ {
					LimiterForget(fmt.Sprint(i))
				}
				runtime.Gosched()
			}
		}
	}()

	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				index := (worker + i) % keyCount
				key := fmt.Sprint(index)
				if i%2 == 0 {
					LimiterLock(key, 0, 0)
				} else if !LimiterTryLock(key, 0) {
					continue
				}
				if active[index].Add(1) != 1 {
					violations.Add(1)
				}
				runtime.Gosched()
				active[index].Add(-1)
				LimiterUnLock(key)
			}
		}(worker)
	}
	wg.Wait()
	close(stop)
	<-done
	if got := violations.Load(); got != 0 {
		t.Fatalf("same-key critical sections overlapped %d times", got)
	}
	r.cleanup(time.Now().Add(time.Second))
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.entries) != 0 {
		t.Fatalf("%d entries could not be reclaimed after callers finished", len(r.entries))
	}
}
