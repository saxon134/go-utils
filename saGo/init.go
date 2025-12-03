package saGo

import "github.com/saxon134/go-utils/saRedis"

var _redis *saRedis.Redis

func Init(redis *saRedis.Redis) {
	_redis = redis
}
