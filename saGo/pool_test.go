package saGo

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/saxon134/go-utils/saRedis"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

func uniqueTestKey(name string) string {
	return fmt.Sprintf("%s-%d", name, time.Now().UnixNano())
}

func TestPool(t *testing.T) {
	fmt.Println("开始", time.Now().Format(time.DateTime))
	var pool = NewPool(10, 10, func(p *Pool, args interface{}) {
		LimiterLock("test", 1, 2)
		time.Sleep(time.Second * 5)
		fmt.Println("执行中：", args)
		LimiterUnLock("test")
	})
	for i := 0; i < 20; i++ {
		pool.Invoke(i + 1)
	}
	pool.Wait()
	fmt.Println("完成", time.Now().Format(time.DateTime))
}

func TestPoolWaitReturnsWhenTaskPanics(t *testing.T) {
	var pool = NewPool(1, 1000, func(p *Pool, args interface{}) {
		panic("task panic")
	})

	pool.Invoke("panic")

	done := make(chan struct{})
	go func() {
		pool.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Pool.Wait blocked after a task panic")
	}
}

func TestPoolWaitDoesNotReturnBeforeConcurrentInvoke(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var releaseOnce sync.Once
	releaseTask := func() {
		releaseOnce.Do(func() { close(release) })
	}
	t.Cleanup(releaseTask)
	pool := NewPool(1, 1000, func(p *Pool, args interface{}) {
		close(started)
		<-release
	})

	// Keep submission in Consume so Wait cannot observe the task through the
	// WaitGroup until the submission is allowed to continue.
	pool.lock.Lock()
	invokeDone := make(chan struct{})
	go func() {
		pool.Invoke(1)
		close(invokeDone)
	}()
	time.Sleep(10 * time.Millisecond)

	waitDone := make(chan struct{})
	go func() {
		pool.Wait()
		close(waitDone)
	}()
	time.Sleep(10 * time.Millisecond)
	pool.lock.Unlock()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("concurrent Invoke was not submitted")
	}

	select {
	case <-waitDone:
		t.Fatal("Pool.Wait returned before the concurrently submitted task completed")
	case <-time.After(20 * time.Millisecond):
	}

	releaseTask()
	select {
	case <-waitDone:
	case <-time.After(time.Second):
		t.Fatal("Pool.Wait did not return after the task completed")
	}
	<-invokeDone
}

func TestPoolInvokeAfterWaitIsIgnored(t *testing.T) {
	var calls int
	pool := NewPool(1, 1000, func(p *Pool, args interface{}) {
		calls++
	})

	pool.Wait()
	pool.Invoke(1)

	if calls != 0 {
		t.Fatalf("Invoke after Wait executed callback %d times, want 0", calls)
	}
}

func TestGoWithParamsRecoversPanic(t *testing.T) {
	if os.Getenv("SAGO_GO_WITH_PARAMS_PANIC") == "1" {
		done := make(chan struct{})
		GoWithParams(nil, func(params interface{}) {
			close(done)
			panic("goroutine panic")
		})
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("GoWithParams did not run callback")
		}
		time.Sleep(50 * time.Millisecond)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestGoWithParamsRecoversPanic")
	cmd.Env = append(os.Environ(), "SAGO_GO_WITH_PARAMS_PANIC=1")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("GoWithParams did not recover panic: %v\n%s", err, out)
	}
}

func TestLimiterLockMaxSecondDoesNotBypassHeldLock(t *testing.T) {
	key := uniqueTestKey("limiter-max")
	LimiterLock(key, 0, 0)

	time.Sleep(20 * time.Millisecond)

	entered := make(chan struct{})
	go func() {
		LimiterLock(key, 0, 0.01)
		close(entered)
	}()

	select {
	case <-entered:
		t.Fatal("LimiterLock entered while the existing lock was still held")
	case <-time.After(30 * time.Millisecond):
	}

	LimiterUnLock(key)

	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("LimiterLock did not enter after the existing lock was released")
	}
	LimiterUnLock(key)
}

func TestLimiterTryLockDoesNotSleepWhenIntervalNotElapsed(t *testing.T) {
	key := uniqueTestKey("limiter-try")
	LimiterLock(key, 0, 0)
	LimiterUnLock(key)

	start := time.Now()
	ok := LimiterTryLock(key, 0.2)
	elapsed := time.Since(start)
	if ok {
		LimiterUnLock(key)
		t.Fatal("LimiterTryLock succeeded before the minimum interval elapsed")
	}
	if elapsed > 50*time.Millisecond {
		t.Fatalf("LimiterTryLock blocked for %s", elapsed)
	}

	if !LimiterTryLock(key, 0) {
		t.Fatal("LimiterTryLock did not release the local lock after interval failure")
	}
	LimiterUnLock(key)
}

func TestLimiterUnLockDoesNotDeleteAnotherRedisOwnerLock(t *testing.T) {
	store := newLimiterRedisStore()
	oldRedis := _redis
	_redis = newLimiterTestRedis(store)
	defer func() {
		_redis = oldRedis
	}()

	key := uniqueTestKey("limiter-redis-owner")
	if !LimiterTryLock(key, 0, LimiterGlobalOption, 1) {
		t.Fatal("LimiterTryLock did not acquire the Redis lock")
	}

	redisKey := "saGo:limiter:" + key
	store.Set(redisKey, "another-owner")
	LimiterUnLock(key)

	if got := store.Get(redisKey); got != "another-owner" {
		t.Fatalf("Redis lock value = %q, want %q", got, "another-owner")
	}
}

func TestLimiterUnLockDeletesMatchingRedisOwnerLock(t *testing.T) {
	store := newLimiterRedisStore()
	oldRedis := _redis
	_redis = newLimiterTestRedis(store)
	defer func() {
		_redis = oldRedis
	}()

	key := uniqueTestKey("limiter-redis-release")
	if !LimiterTryLock(key, 0, LimiterGlobalOption, 1) {
		t.Fatal("LimiterTryLock did not acquire the Redis lock")
	}

	redisKey := "saGo:limiter:" + key
	LimiterUnLock(key)

	if got := store.Get(redisKey); got != "" {
		t.Fatalf("Redis lock value = %q, want empty", got)
	}
}

func TestLimiterLockUnLockDoesNotDeleteAnotherRedisOwnerLock(t *testing.T) {
	store := newLimiterRedisStore()
	oldRedis := _redis
	_redis = newLimiterTestRedis(store)
	defer func() {
		_redis = oldRedis
	}()

	key := uniqueTestKey("limiter-lock-redis-owner")
	LimiterLock(key, 0, 1, LimiterGlobalOption)

	redisKey := "saGo:limiter:" + key
	store.Set(redisKey, "another-owner")
	LimiterUnLock(key)

	if got := store.Get(redisKey); got != "another-owner" {
		t.Fatalf("Redis lock value = %q, want %q", got, "another-owner")
	}
}

func TestLimiterLockUsesTenMinuteDefaultRedisLease(t *testing.T) {
	store := newLimiterRedisStore()
	oldRedis := _redis
	_redis = newLimiterTestRedis(store)
	defer func() {
		_redis = oldRedis
	}()

	key := uniqueTestKey("limiter-default-lease")
	LimiterLock(key, 0, 0, LimiterGlobalOption)
	defer LimiterUnLock(key)

	redisKey := "saGo:limiter:" + key
	if got := store.Expire(redisKey); got != "600" {
		t.Fatalf("Redis lease = %s seconds, want 600", got)
	}
}

func TestLimiterUnLockLocalLockDoesNotDeleteRedisLock(t *testing.T) {
	store := newLimiterRedisStore()
	oldRedis := _redis
	_redis = newLimiterTestRedis(store)
	defer func() {
		_redis = oldRedis
	}()

	key := uniqueTestKey("limiter-local-owner")
	redisKey := "saGo:limiter:" + key
	store.Set(redisKey, "another-owner")

	LimiterLock(key, 0, 0)
	LimiterUnLock(key)

	if got := store.Get(redisKey); got != "another-owner" {
		t.Fatalf("Redis lock value = %q, want %q", got, "another-owner")
	}
}

func TestNewBucketSetsPositiveIntervalForQPSOnly(t *testing.T) {
	oldRedis := _redis
	_redis = nil
	defer func() {
		_redis = oldRedis
	}()

	bucket := NewBucket(10, 0)
	if bucket == nil {
		t.Fatal("NewBucket returned nil")
	}
	if bucket.minIntervalMillisecond <= 0 {
		t.Fatalf("minIntervalMillisecond = %d, want positive", bucket.minIntervalMillisecond)
	}
}

func TestBucketConsumeDoesNotSendInvalidRedisIncr(t *testing.T) {
	recorder := &fakeRedisRecorder{}
	r := &saRedis.Redis{Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return &fakeRedisConn{recorder: recorder}, nil
		},
	}}

	bucket := NewBucket(10, 0, r, uniqueTestKey("bucket-redis"))
	bucket.Consume()

	for _, cmd := range recorder.Commands() {
		if cmd.name == "INCR" && len(cmd.args) != 1 {
			t.Fatalf("INCR called with invalid args: %#v", cmd.args)
		}
	}
}

func TestCleanUpdatesLastCleanTime(t *testing.T) {
	cacheLock.Lock()
	caches = map[string]*CacheItem{}
	lastCleanTime = 0
	cacheLock.Unlock()

	GetCache(uniqueTestKey("cache-clean"), 0, nil)

	cacheLock.Lock()
	defer cacheLock.Unlock()
	if lastCleanTime == 0 {
		t.Fatal("clean did not update lastCleanTime")
	}
}

func TestCleanKeepsRecentlyAccessedLowFrequencyItems(t *testing.T) {
	cacheLock.Lock()
	caches = map[string]*CacheItem{}
	lastCleanTime = 0
	now := time.Now()
	recentKey := uniqueTestKey("cache-recent")
	caches[recentKey] = &CacheItem{Value: "recent", GetAt: now, Count: 0}
	for i := 0; i < 1001; i++ {
		caches[fmt.Sprintf("old-%d", i)] = &CacheItem{
			Value: "old",
			GetAt: now.Add(-time.Minute),
			Count: 2,
		}
	}
	clean()
	_, ok := caches[recentKey]
	caches = map[string]*CacheItem{}
	cacheLock.Unlock()

	if !ok {
		t.Fatal("clean removed a recently accessed cache item")
	}
}

type redisCommand struct {
	name string
	args []interface{}
}

type fakeRedisRecorder struct {
	lock     sync.Mutex
	commands []redisCommand
}

func (r *fakeRedisRecorder) Record(name string, args []interface{}) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.commands = append(r.commands, redisCommand{
		name: strings.ToUpper(name),
		args: append([]interface{}(nil), args...),
	})
}

func (r *fakeRedisRecorder) Commands() []redisCommand {
	r.lock.Lock()
	defer r.lock.Unlock()
	return append([]redisCommand(nil), r.commands...)
}

type fakeRedisConn struct {
	recorder *fakeRedisRecorder
}

func (c *fakeRedisConn) Close() error {
	return nil
}

func (c *fakeRedisConn) Err() error {
	return nil
}

func (c *fakeRedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	c.recorder.Record(commandName, args)
	switch strings.ToUpper(commandName) {
	case "GET":
		return nil, redis.ErrNil
	case "INCR":
		if len(args) != 1 {
			return nil, fmt.Errorf("wrong number of arguments for INCR")
		}
		return int64(1), nil
	case "EXPIRE":
		return int64(1), nil
	case "EVAL":
		return int64(1), nil
	default:
		return "OK", nil
	}
}

func (c *fakeRedisConn) Send(commandName string, args ...interface{}) error {
	c.recorder.Record(commandName, args)
	return nil
}

func (c *fakeRedisConn) Flush() error {
	return nil
}

func (c *fakeRedisConn) Receive() (interface{}, error) {
	return nil, nil
}

type limiterRedisStore struct {
	lock    sync.Mutex
	values  map[string]string
	expires map[string]string
}

func newLimiterRedisStore() *limiterRedisStore {
	return &limiterRedisStore{
		values:  map[string]string{},
		expires: map[string]string{},
	}
}

func (s *limiterRedisStore) Set(key, value string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values[key] = value
}

func (s *limiterRedisStore) Get(key string) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.values[key]
}

func (s *limiterRedisStore) Expire(key string) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.expires[key]
}

func newLimiterTestRedis(store *limiterRedisStore) *saRedis.Redis {
	return &saRedis.Redis{Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return &limiterRedisConn{store: store}, nil
		},
	}}
}

type limiterRedisConn struct {
	store *limiterRedisStore
}

func (c *limiterRedisConn) Close() error { return nil }
func (c *limiterRedisConn) Err() error   { return nil }

func (c *limiterRedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	switch strings.ToUpper(commandName) {
	case "SET":
		key := fmt.Sprint(args[0])
		value := fmt.Sprint(args[1])
		c.store.lock.Lock()
		defer c.store.lock.Unlock()
		if _, ok := c.store.values[key]; ok {
			return nil, redis.ErrNil
		}
		c.store.values[key] = value
		c.store.expires[key] = fmt.Sprint(args[3])
		return "OK", nil
	case "GET":
		key := fmt.Sprint(args[0])
		c.store.lock.Lock()
		defer c.store.lock.Unlock()
		value, ok := c.store.values[key]
		if !ok {
			return nil, redis.ErrNil
		}
		return value, nil
	case "DEL":
		key := fmt.Sprint(args[0])
		c.store.lock.Lock()
		defer c.store.lock.Unlock()
		delete(c.store.values, key)
		return int64(1), nil
	case "EVAL":
		key := fmt.Sprint(args[2])
		token := fmt.Sprint(args[3])
		c.store.lock.Lock()
		defer c.store.lock.Unlock()
		if c.store.values[key] == token {
			delete(c.store.values, key)
			return int64(1), nil
		}
		return int64(0), nil
	default:
		return nil, fmt.Errorf("unexpected command %s", commandName)
	}
}

func (c *limiterRedisConn) Send(commandName string, args ...interface{}) error { return nil }
func (c *limiterRedisConn) Flush() error                                       { return nil }
func (c *limiterRedisConn) Receive() (interface{}, error)                      { return nil, nil }
