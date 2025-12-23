package saRedis

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/saxon134/go-utils/saData"
	"strconv"
	"strings"
	"time"
)

type Redis struct {
	*redis.Pool
	Project string //默认不会在key上面拼接，如果需要区分项目，需要手动拼接
}

func Init(uri string, pass string, db int) (r *Redis, err error) {
	if uri == "" {
		return nil, errors.New("URI不能空")
	}

	r = &Redis{Pool: &redis.Pool{
		MaxIdle:         3,
		MaxActive:       20,
		IdleTimeout:     180 * time.Second,
		MaxConnLifetime: 10 * time.Minute,
		Wait:            true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", uri)
			if err != nil {
				return nil, err
			}

			if pass != "" {
				if _, err = c.Do("AUTH", pass); err != nil {
					_ = c.Close()
					return nil, err
				}
			}

			if db >= 0 {
				if _, err = c.Do("SELECT", db); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}}

	_, err = r.Pool.Dial()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r Redis) Set(k string, v interface{}, expire time.Duration) error {
	if k == "" {
		return errors.New("缺少必要参数")
	}

	c := r.Pool.Get()
	defer c.Close()

	var err error
	if v == "" {
		_, err = c.Do("DEL", k)
	} else {
		if expire > 0 {
			_, err = c.Do("SET", k, v, "EX", strconv.FormatInt(int64(expire.Seconds()), 10))
		} else {
			_, err = c.Do("SET", k, v)
		}
	}
	return err
}

func (r Redis) Get(k string) (interface{}, error) {
	c := r.Pool.Get()
	defer c.Close()

	if k != "" {
		return c.Do("GET", k)
	}
	return nil, errors.New("")
}

func (r Redis) Clear(k string) error {
	if k != "" {
		c := r.Pool.Get()
		defer c.Close()

		_, err := c.Do("DEL", k)
		return err
	}
	return nil
}

func (r Redis) SetObj(k string, objPtr interface{}, expire time.Duration) error {
	if k == "" {
		return errors.New("缺少必要参数")
	}

	c := r.Pool.Get()
	defer c.Close()

	if objPtr == "" || objPtr == nil {
		_, err := c.Do("DEL", k)
		return err
	} else {
		if bAry, err := json.Marshal(objPtr); err == nil && bAry != nil {
			if expire > 0 {
				_, err := c.Do("SET", k, bAry, "EX", strconv.FormatInt(int64(expire.Seconds()), 10))
				return err
			} else {
				_, err := c.Do("SET", k, bAry)
				return err
			}
		} else {
			return err
		}
	}
}

func (r Redis) GetObj(k string, objPtr interface{}) error {
	if k != "" {
		c := r.Pool.Get()
		defer c.Close()

		if bAry, err := redis.Bytes(c.Do("GET", k)); err == nil {
			err = saData.BytesToModel(bAry, objPtr)
			return err
		} else {
			return err
		}
	}
	return errors.New("缺少key")
}

func (r Redis) GetString(k string) (string, error) {
	if k != "" {
		c := r.Pool.Get()
		defer c.Close()

		str, err := redis.String(c.Do("GET", k))
		if err != nil {
			return "", err
		}

		return str, nil
	}

	return "", errors.New("缺少key")
}

func (r Redis) GetInt64(k string) (int64, error) {
	if k != "" {
		c := r.Pool.Get()
		defer c.Close()

		str, err := redis.String(c.Do("GET", k))
		if err != nil {
			return 0, err
		}

		return saData.ToInt64(str)
	}

	return 0, errors.New("缺少key")
}

func (r Redis) IsError(err error) bool {
	if err == nil || strings.Index(err.Error(), "nil returned") >= 0 {
		return false
	}
	return true
}
