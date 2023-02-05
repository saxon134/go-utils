package saCache

import (
	"sort"
	"time"
)

var simpleCacheData map[string]simpleCache

type simpleCache struct {
	V interface{}
	T time.Time
}

// SMGet 简单版本的缓存
func SMGet(key string) (value interface{}) {
	if simpleCacheData == nil {
		simpleCacheData = map[string]simpleCache{}
	}

	data, ok := simpleCacheData[key]
	if ok == false {
		return nil
	}

	var now = time.Now()
	if data.T.IsZero() || data.T.After(now) {
		delete(simpleCacheData, key)
		return nil
	}
	return data.V
}

// SMSet 简单版本的缓存 最多存储100条
func SMSet(key string, value interface{}, duration time.Duration) {
	if simpleCacheData == nil {
		simpleCacheData = map[string]simpleCache{}
	}
	
	if len(simpleCacheData) > 100 {
		var timeAry = make([]int64, 0, len(simpleCacheData))
		for _, v := range simpleCacheData {
			timeAry = append(timeAry, v.T.Unix())
		}
		sort.Slice(timeAry, func(i, j int) bool {
			return timeAry[i] < timeAry[j]
		})

		var t = timeAry[50]
		for k, v := range simpleCacheData {
			if v.T.Unix() < t {
				delete(simpleCacheData, k)
			}
		}
	}

	if key == "" || value == nil || duration <= 0 {
		return
	}

	var data = simpleCache{
		V: value,
		T: time.Now().Add(duration),
	}
	simpleCacheData[key] = data
}
