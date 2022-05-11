package saCache

import (
	"errors"
	"github.com/saxon134/go-utils/saData"
	"sync"
	"time"
)

type CacheHandle func(id string) (interface{}, error)

//最多保存的类型，超出会覆盖掉最早访问的存储项
const maxCountForKind = 50

//每个项目最多保存的数量，超出会按照参数mode规则去替换
const maxCountForOneKind = 50

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
//             5m支持范围：10s - 60m 默认10s
//             典型场景：IP限流
func MGetWithFunc(key string, id string, mode string, handle CacheHandle) (value interface{}, cnt int, err error) {
	if key == "" || id == "" || handle == nil {
		return nil, 0, errors.New("缺少参数或取值方法")
	}

	now := time.Now().Unix()
	var retentionSecond int64 = 10
	if mode != "" {
		m := mode[len(mode)-1:]
		t, _ := saData.ToInt64(mode[:len(mode)-1])
		if m == "m" {
			if t > 0 && t <= 60 {
				retentionSecond = t * 60
			}
		} else if m == "s" {
			if t > 10 && t <= 60 {
				retentionSecond = t
			}
		}
	}

	_locker.Lock()
	defer _locker.Unlock()

	cacheKind := _cache[key]
	if cacheKind == nil {
		cacheKind = new(cache)
		cacheKind.Ary = make([]cacheItem, 0, 100)
		cacheKind.MaxCnt = 1
		cacheKind.TotalCnt = 1
		cacheKind.LastTime = now //只有新建的时候才设置时间
	}
	cacheKind.TotalCnt++

	var item *cacheItem = nil
	var itemIdx = -1
	for i, v := range cacheKind.Ary {
		if v.Id == id {
			item = &v
			itemIdx = i
			break
		}
	}

	//不存在，或已过期
	if itemIdx < 0 || item.LastTime+retentionSecond < now {
		if itemIdx >= 0 {
			cacheKind.Ary = append(cacheKind.Ary[:itemIdx], cacheKind.Ary[itemIdx+1:]...)
		}

		v, err := handle(id)
		if err != nil {
			return nil, 0, err
		}

		item = &cacheItem{
			Cnt:      1,
			Id:       id,
			LastTime: now,
			V:        v,
		}
		//每个类目最多保存的数量
		if len(cacheKind.Ary) >= maxCountForOneKind {
			//取次数最小的，替换掉
			var minIdx = -1
			for i, v := range cacheKind.Ary {
				if minIdx == -1 || cacheKind.Ary[minIdx].Cnt > v.Cnt {
					minIdx = i
				}
			}
			cacheKind.Ary = append(cacheKind.Ary[:minIdx], cacheKind.Ary[minIdx+1:]...)
		}
		cacheKind.Ary = append(cacheKind.Ary, *item)
	} else {
		item.Cnt++
		cacheKind.Ary[itemIdx] = *item

		//更新最大次数，新增加的时候会使用到
		cacheKind.MaxCnt = 0
		for _, v := range cacheKind.Ary {
			if v.Cnt > cacheKind.MaxCnt {
				cacheKind.MaxCnt = v.Cnt
			}
		}
	}

	//最多保存50类数据
	if len(_cache) >= maxCountForKind && _cache[key] == nil {
		var min *cache
		var minK string
		for k, v := range _cache {
			if min == nil || min.TotalCnt == 0 || min.TotalCnt > v.TotalCnt {
				min = v
				minK = k
			}
		}

		if minK != "" {
			delete(_cache, minK)
		}
	}

	_cache[key] = cacheKind
	return item.V, item.Cnt, nil
}

// MGet
// 只有提前调用了RegisterHandle将方法注册进来后才可以调用该接口，否则返回数据会是空的
func MGet(key string, mode string, id string) (value interface{}, cnt int, err error) {
	value, cnt, err = MGetWithFunc(key, id, mode, _handle[key])
	return
}
