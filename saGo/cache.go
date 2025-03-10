package saGo

import (
	"sync"
	"time"
)

//简单的内存缓存，存储一些配置信息，设置过期时间

var caches = map[string]*CacheItem{}
var cacheLock sync.Mutex
var lastCleanTime int64

type CacheItem struct {
	ExpireAt time.Time
	GetAt    time.Time
	Count    int
	Value    interface{}
}

type CacheHandler func(duration time.Duration) interface{}

// 只能在GetCache内调用，处理未加锁，依赖GetCache内的锁
func clean() {
	//10秒之内不清理
	var now = time.Now()
	if now.Unix()-lastCleanTime < 10 {
		return
	}

	var keys = make([]string, 0, len(caches))
	for k, _ := range caches {
		keys = append(keys, k)
	}

	//删除过期数据
	for _, k := range keys {
		var v = caches[k]
		if v.ExpireAt.IsZero() == false {
			if v.ExpireAt.Before(now) {
				delete(caches, k)
			}
		}
	}

	//数据太长，删除访问频率低的
	if len(caches) > 1000 {
		var count = 0
		for _, v := range caches {
			count += v.Count
		}

		var avage = count / len(caches)
		var keys = make([]string, 0, 100)
		for k, v := range caches {
			if v.Count < avage {
				if v.GetAt.Before(now.Add(time.Second * 30)) {
					keys = append(keys, k)
				}
			}
		}
		for _, k := range keys {
			delete(caches, k)
		}
	}
}

func GetCache(key string, duration time.Duration, fn CacheHandler) interface{} {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	var now = time.Now()
	var value = caches[key]
	if value != nil && value.ExpireAt.IsZero() == false {
		if value.ExpireAt.Before(now) {
			delete(caches, key)
			value = nil
		}
	}

	if value == nil {
		if fn != nil {
			var v = fn(duration)
			if v != nil {
				value = &CacheItem{
					Value:    v,
					ExpireAt: now.Add(duration),
					GetAt:    now,
					Count:    1,
				}
				caches[key] = value
				return value.Value
			}
		}
	} else {
		value.GetAt = now
		value.Count++
		return value.Value
	}

	clean()
	return nil
}

func SetCache(key string, value interface{}, duration time.Duration) {
	var now = time.Now()
	if key == "" || value == nil || duration <= 0 {
		return
	}

	cacheLock.Lock()
	defer cacheLock.Unlock()

	var item = &CacheItem{
		Value: value,
		GetAt: now,
	}

	//未设置表示不失效
	if duration > 0 {
		item.ExpireAt = time.Now().Add(duration)
	}

	caches[key] = item
}
