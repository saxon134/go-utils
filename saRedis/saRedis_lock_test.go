package saRedis

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func TestUnlockDoesNotDeleteAnotherOwnerLock(t *testing.T) {
	store := newFakeRedisStore()
	r := Redis{Pool: &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return &fakeRedisConn{store: store}, nil
		},
	}}

	lock, err := r.TryLock("resource", 10)
	if err != nil {
		t.Fatal(err)
	}
	store.Set(lock.key(), "other-owner")

	if err := lock.Unlock(); err != nil {
		t.Fatal(err)
	}

	if got := store.Get(lock.key()); got != "other-owner" {
		t.Fatalf("lock value = %q, want other-owner", got)
	}
}

type fakeRedisStore struct {
	lock   sync.Mutex
	values map[string]string
}

func newFakeRedisStore() *fakeRedisStore {
	return &fakeRedisStore{values: map[string]string{}}
}

func (s *fakeRedisStore) Set(key, value string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.values[key] = value
}

func (s *fakeRedisStore) Get(key string) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.values[key]
}

type fakeRedisConn struct {
	store *fakeRedisStore
}

func (c *fakeRedisConn) Close() error { return nil }
func (c *fakeRedisConn) Err() error   { return nil }

func (c *fakeRedisConn) Do(commandName string, args ...interface{}) (interface{}, error) {
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
		c.store.lock.Lock()
		defer c.store.lock.Unlock()
		for _, arg := range args {
			delete(c.store.values, fmt.Sprint(arg))
		}
		return int64(len(args)), nil
	case "EXPIRE":
		return int64(1), nil
	case "EVAL":
		key := fmt.Sprint(args[2])
		code := fmt.Sprint(args[3])
		c.store.lock.Lock()
		defer c.store.lock.Unlock()
		if c.store.values[key] == code {
			delete(c.store.values, key)
			return int64(1), nil
		}
		return int64(0), nil
	default:
		return nil, fmt.Errorf("unexpected command %s", commandName)
	}
}

func (c *fakeRedisConn) Send(commandName string, args ...interface{}) error { return nil }
func (c *fakeRedisConn) Flush() error                                       { return nil }
func (c *fakeRedisConn) Receive() (interface{}, error)                      { return nil, nil }
