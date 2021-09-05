package saCache

import (
	"github.com/saxon134/go-utils/saHit"
	"sync"
)

type CacheHandle func(id string) (interface{}, error)

var _cache = make(map[string]*cache, 20)
var _locker sync.RWMutex
var _handle map[string]CacheHandle

/**
建议都提前调用该接口注册方法，注册后其他获取缓存接口行为才会一致
否则获取缓存可能无法命中 */
func RegisterHandle(key string, handle CacheHandle) {
	if key == "" || handle == nil {
		return
	}
	if _handle == nil || len(_handle) == 0 {
		_handle = make(map[string]CacheHandle, 20)
	}

	_handle[key] = handle
}

/**
该接口的handle不会自动注册的，因为各业务内handle可能都会有一点点差异*/
func MGetWithFunc(key string, id string, handle CacheHandle) (value interface{}, err error) {
	if key == "" {
		return
	}

	cacheValue := _cache[key]
	if cacheValue == nil {
		cacheValue = new(cache)
	}

	var c *cacheItem = nil
	for _, v := range cacheValue.ary {
		if v.id == id {
			c = &v
			break
		}
	}

	if c != nil {
		_locker.Lock()
		defer _locker.Unlock()

		c.cnt++
		if c.cnt > cacheValue.maxCnt {
			cacheValue.maxCnt = c.cnt
		}
		return c.v, nil
	} else if handle != nil {
		v, err := handle(id)
		if err != nil {
			return nil, err
		}

		c = new(cacheItem)
		_locker.Lock()
		defer _locker.Unlock()

		c.cnt = saHit.Int(cacheValue.maxCnt > 1, cacheValue.maxCnt/2, 1)
		c.v = v
		c.id = id
		if len(cacheValue.ary) < 100 {
			cacheValue.ary = append(cacheValue.ary, *c)
		} else {
			//取次数最小的，替换掉
			var min *cacheItem
			var idx = 0
			for i, v := range cacheValue.ary {
				if min == nil || min.cnt == 0 || min.cnt > v.cnt {
					min = &v
					idx = i
				}
			}
			if min != nil && min.cnt >= 0 {
				cacheValue.ary = append(cacheValue.ary[:idx], cacheValue.ary[idx+1:]...)
			}
		}

		//最多保存50类数据
		if len(_cache) < 50 {
			_cache[key] = cacheValue
		}
		return v, nil
	}
	return nil, nil
}

/**
只有提前调用了RegisterHandle将方法注册进来后才可以调用该接口，否则返回数据会是空的
除非之前有调用MGetWithFunc有缓存才可能命中*/
func MGet(key string, id string) (value interface{}, err error) {
	value, err = MGetWithFunc(key, id, _handle[key])
	return
}
