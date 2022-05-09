package saCache

import (
	"github.com/saxon134/go-utils/saData"
	"github.com/saxon134/go-utils/saHit"
	"strings"
	"sync"
	"time"
)

type CacheHandle func(id string) (interface{}, error)

//最多保存的类型，超出会覆盖掉最早访问的存储项
const maxCountForKind = 50

//每个项目最多保存的数量，超出会按照参数mode规则去替换
const maxCountForOneKind = 100

var _cache = make(map[string]*cache, maxCountForKind)
var _locker sync.RWMutex
var _handle map[string]CacheHandle

// RegisterHandle
// 建议都提前调用该接口注册方法，注册后其他获取缓存接口行为才会一致
// 否则获取缓存可能无法命中
func RegisterHandle(key string, handle CacheHandle) {
	if key == "" || handle == nil {
		return
	}
	if _handle == nil || len(_handle) == 0 {
		_handle = make(map[string]CacheHandle, 20)
	}

	_handle[key] = handle
}

// MGetWithFunc
// 该接口的handle不会自动注册的，因为各业务内handle可能都会有一点点差异
// mode
//    r或者空 -替换使用频次最低的存储项，
//             典型场景：缓存存在数据库的配置，减少数据库压力
//    5m      -替换最后访问时间最早的那个，如果最早的那个在5m之内，则还是会替换访问次数最少的存储项
//             5m支持范围：10s - 60m
//             典型场景：IP限流
func MGetWithFunc(key string, id string, mode string, handle CacheHandle) (value interface{}, err error) {
	if key == "" {
		return
	}

	now := time.Now().Unix()
	var retentionSecond int64
	if mode != "" {
		var isMinute = true
		if strings.HasSuffix(mode, "s") {
			isMinute = false
		}
		t, _ := saData.ToInt(mode[:len(mode)-1])
		if isMinute == false && t < 10 {
			//最小10秒
		} else if t > 0 && t <= 60 {
			if isMinute {
				retentionSecond += 60 * 60
			}
		}
	}

	cacheKind := _cache[key]
	if cacheKind == nil {
		cacheKind = new(cache)
	}
	cacheKind.totalCnt++

	var c *cacheItem = nil
	for _, v := range cacheKind.ary {
		if v.id == id {
			c = &v
			break
		}
	}

	if c != nil {
		_locker.Lock()
		defer _locker.Unlock()

		c.cnt++
		if c.cnt > cacheKind.maxCnt {
			cacheKind.maxCnt = c.cnt
		}
		c.lastTime = now
		return c.v, nil
	} else if handle != nil {
		v, err := handle(id)
		if err != nil {
			return nil, err
		}

		c = new(cacheItem)
		_locker.Lock()
		defer _locker.Unlock()

		c.cnt = saHit.Int(cacheKind.maxCnt > 1, cacheKind.maxCnt/2, 1)
		c.v = v
		c.id = id
		c.lastTime = now

		//每个类目最多保存的数量
		if len(cacheKind.ary) < maxCountForOneKind {
			cacheKind.ary = append(cacheKind.ary, *c)
		} else {
			//取次数最小的，替换掉
			var min1 *cacheItem
			var min2 *cacheItem
			var idx1 = 0
			var idx2 = 0
			for i, v := range cacheKind.ary {
				if v.lastTime+retentionSecond < now {
					if min1 == nil || min1.cnt == 0 || min1.cnt > v.cnt {
						min1 = &v
						idx1 = i
					}
				} else {
					if min2 == nil || min2.cnt == 0 || min2.cnt > v.cnt {
						min2 = &v
						idx2 = i
					}
				}
			}

			if min1 != nil && min1.cnt > 0 {
				cacheKind.ary = append(cacheKind.ary[:idx1], cacheKind.ary[idx1+1:]...)
			} else if min2 != nil && min2.cnt > 0 {
				cacheKind.ary = append(cacheKind.ary[:idx2], cacheKind.ary[idx2+1:]...)
			}
		}

		//最多保存50类数据
		if len(_cache) < maxCountForKind {
			_cache[key] = cacheKind
		} else {
			var min *cache
			var minK string
			for k, v := range _cache {
				if v.lastTime+retentionSecond < now {
					if min == nil || min.totalCnt == 0 || min.totalCnt > v.totalCnt {
						min = v
						minK = k
					}
				}
			}

			if min != nil && min.totalCnt > 0 {
				_cache[minK] = cacheKind
			}
		}
		return v, nil
	}
	return nil, nil
}

// MGet
// 只有提前调用了RegisterHandle将方法注册进来后才可以调用该接口，否则返回数据会是空的
func MGet(key string, mode string, id string) (value interface{}, err error) {
	value, err = MGetWithFunc(key, id, mode, _handle[key])
	return
}
