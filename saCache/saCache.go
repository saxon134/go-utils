package saCache

import (
	"sort"
	"sync"
	"time"
)

type CacheHandle func() (interface{}, error)

var cacheData map[string]*Cache
var cacheLocker sync.RWMutex

func MGet(key string) (value interface{}, expired bool) {
	if cacheData == nil {
		cacheData = map[string]*Cache{}
	}

	data, ok := cacheData[key]
	if ok == false || data == nil {
		return nil, true
	}

	var now = time.Now()
	if data.Before.IsZero() || data.Before.After(now) {
		value = data.Data
		expired = true
		delete(cacheData, key)
		return value, true
	}

	return data.Data, false
}

// MSetWithFunc 最多500条，超出会删除最早的数据，一次性删除三分之一
// 防止并发更新：调用handler前必须抢到锁，并且更新时间超过1秒
func MSetWithFunc(key string, duration time.Duration, handler CacheHandle) (value interface{}, expired bool) {
	if key == "" || duration <= time.Second || handler == nil {
		return nil, false
	}

	//大锁
	cacheLocker.Lock()
	defer cacheLocker.Unlock()

	if cacheData == nil {
		cacheData = map[string]*Cache{}
	}

	//数据较多时，删除数据
	defer func() {
		if len(cacheData) >= 500 {
			var timeAry = make([]int64, 0, len(cacheData))
			for _, v := range cacheData {
				timeAry = append(timeAry, v.Before.Unix())
			}
			sort.Slice(timeAry, func(i, j int) bool {
				return timeAry[i] < timeAry[j]
			})

			var t = timeAry[300]
			for k, v := range cacheData {
				if v.Before.Unix() < t {
					delete(cacheData, k)
				}
			}
		}
	}()

	var data = cacheData[key]
	if data == nil {
		data = &Cache{
			Before:   time.Now().Add(duration),
			Data:     nil,
			UpdateAt: time.Time{},
		}
		cacheData[key] = data
	}

	//无数据情况，同步调用获取数据的处理
	if data.Data == nil {
		data.isUpdating = true
		v, e := handler()
		if e == nil {
			data.Data = v
			data.UpdateAt = time.Now()
			data.isUpdating = false
			cacheData[key] = data
			return data.Data, false
		}
		return nil, false
	}

	//有数据，并且未过期
	if data.Data != nil  {
		var now = time.Now()

		//未过期
		if data.Before.IsZero() == false && data.Before.Before(now) {
			return data.Data, false
		}

		//刚更新过，避免并发多次更新
		if data.UpdateAt.IsZero()==false && data.UpdateAt.After(now.Add(time.Second*-1)) {
			return data.Data, false
		}
	}

	//有数据，但是过期了，先返回数据，再异步调用获取数据处理
	if data.isUpdating == false {
		data.isUpdating = true
		cacheData[key] = data
		go func() {
			v, e := handler()
			if e == nil {
				data.Data = v
				data.UpdateAt = time.Now()
				cacheData[key] = data
			}
			data.isUpdating = false
		}()
	}
	return data.Data, true
}
