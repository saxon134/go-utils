package saRedis

import (
	"github.com/gomodule/redigo/redis"
	"strings"
)

func (r Redis) Do(command string, args ...interface{}) (res interface{}, err error) {
	c := r.Pool.Get()
	defer c.Close()

	res, err = c.Do(command, args...)
	return
}

func (r Redis) WriteExpireAndNotExist(key, value string, expireTime int64) (ok bool) {
	c := r.Pool.Get()
	defer c.Close()

	res, _ := redis.String(c.Do("SET", key, value, "EX", expireTime, "NX"))
	return strings.ToUpper(res) == "OK"
}
