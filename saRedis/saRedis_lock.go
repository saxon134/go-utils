package saRedis

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"strings"
	"time"
)

type Lock struct {
	resource        string
	second          int    // 时间
	startTimeSecond int64  //非零，表示有最小锁定时间限制
	code            string // 锁的值，避免刷新其他协程的时间
	conn            *Redis
	c               chan int // 主动解锁的时候，发送消息退出协程
}

func (r Redis) TryLock(resource string, timeoutSecond int) (lock *Lock, err error) {
	lock = &Lock{
		resource:        resource,
		second:          timeoutSecond,
		startTimeSecond: 0,
		code:            strings.Replace(uuid.NewV4().String(), "-", "", -1),
		conn:            &r,
	}

	err = lock.tryLock()
	if err != nil {
		return nil, err
	}

	return
}

func (r Redis) TryLockWithMinTime(resource string, minTimeSecond int) (lock *Lock, err error) {
	lock = &Lock{
		resource:        resource,
		second:          minTimeSecond,
		startTimeSecond: time.Now().Unix(),
		code:            strings.Replace(uuid.NewV4().String(), "-", "", -1),
		conn:            &r,
	}

	err = lock.tryLock()
	if err != nil {
		return nil, err
	}

	return
}

func (lock *Lock) Unlock() (err error) {
	if lock.startTimeSecond <= 0 {
		_, err = lock.conn.Do("EVAL", redisUnlockScript, 1, lock.key(), lock.code)
	}
	if lock.c != nil {
		select {
		case lock.c <- 1: //退出刷新时间协程
		default:
		}
	}
	return
}

func (lock *Lock) tryLock() (err error) {
	if lock == nil || lock.conn == nil || lock.resource == "" {
		return errors.New("error, connect is empty")
	}
	if lock.second <= 0 {
		return errors.New("lock time is empty")
	}

	var res string
	res, err = redis.String(lock.conn.Do("SET", lock.key(), lock.code, "NX", "EX", lock.second))
	if err != nil {
		return err
	}

	if strings.ToUpper(res) == "OK" {
		lock.c = make(chan int, 1)
		go lock.refreshTimeout()

		return nil
	}

	return errors.New("加锁失败")
}

// 每2/3过期时间，刷新一下锁时间，避免处理过程超过了超时时间，导致锁被释放
func (lock *Lock) refreshTimeout() {
	for {
		second := 1
		if lock.second >= 3 {
			second = int(2 * lock.second / 3)
		}

		t := time.NewTimer(time.Second * (time.Duration(second)))
		select {
		case <-lock.c:
			if lock.startTimeSecond <= 0 {
				t.Stop()
				return
			} else {
				if lock.startTimeSecond+int64(lock.second) <= time.Now().Unix() {
					t.Stop()
					return
				} else {
					lock.resource = ""
				}
			}
		case <-t.C:
			if lock == nil || lock.resource == "" {
				t.Stop()
				return
			}

			_, err := lock.conn.Do("EVAL", redisRefreshScript, 1, lock.key(), lock.code, lock.second)
			if err != nil {
				t.Stop()
				return
			}
		}
	}
}

func (lock *Lock) key() string {
	return fmt.Sprintf("lock:%s", lock.resource)
}

const redisUnlockScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`

const redisRefreshScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("EXPIRE", KEYS[1], ARGV[2])
end
return 0
`
