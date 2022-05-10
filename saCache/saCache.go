package saCache

import (
	"errors"
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
		t, _ := saData.ToInt64(mode[:len(mode)-1])
		if isMinute == false && t < 10 {
			//最小10秒
		} else if t > 0 && t <= 60 {
			if isMinute {
				retentionSecond += 60 * 60
			} else {
				retentionSecond += t
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
	}
	cacheKind.TotalCnt++
	cacheKind.LastTime = now

	var c *cacheItem = nil
	var cIdx = -1
	for i, v := range cacheKind.Ary {
		if v.Id == id {
			c = &v
			cIdx = i
			break
		}
	}

	if cIdx >= 0 {
		c.Cnt++
		if c.Cnt > cacheKind.MaxCnt {
			cacheKind.MaxCnt = c.Cnt
		}
		c.LastTime = now
		cacheKind.Ary[cIdx] = *c
	} else if handle != nil {
		v, err := handle(id)
		if err != nil {
			return nil, err
		}

		c = new(cacheItem)
		c.Cnt = saHit.Int(cacheKind.MaxCnt > 1, cacheKind.MaxCnt/2, 1)
		c.V = v
		c.Id = id
		c.LastTime = now

		//每个类目最多保存的数量
		if len(cacheKind.Ary) < maxCountForOneKind {
			cacheKind.Ary = append(cacheKind.Ary, *c)
		} else {
			//取次数最小的，替换掉；idx1是超时，idx2是未超时
			var minIdx1 = -1
			var minIdx2 = -1
			for i, v := range cacheKind.Ary {
				if v.LastTime+retentionSecond < now {
					if minIdx1 == -1 || cacheKind.Ary[minIdx1].Cnt > v.Cnt {
						minIdx1 = i
					}
				} else {
					if minIdx2 == -1 || cacheKind.Ary[minIdx2].Cnt > v.Cnt {
						minIdx2 = i
					}
				}
			}

			if retentionSecond > 0 {
				//优先删除已超时的数据
				if minIdx1 >= 0 {
					cacheKind.Ary = append(cacheKind.Ary[:minIdx1], cacheKind.Ary[minIdx1+1:]...)
				} else if minIdx2 >= 0 {
					cacheKind.Ary = append(cacheKind.Ary[:minIdx2], cacheKind.Ary[minIdx2+1:]...)
				}
			} else {
				if minIdx1 >= 0 {
					if minIdx2 == -1 || cacheKind.Ary[minIdx1].Cnt < cacheKind.Ary[minIdx2].Cnt {
						cacheKind.Ary = append(cacheKind.Ary[:minIdx1], cacheKind.Ary[minIdx1+1:]...)
					}
				} else {
					cacheKind.Ary = append(cacheKind.Ary[:minIdx2], cacheKind.Ary[minIdx2+1:]...)
				}
			}
		}
	} else {
		return nil, errors.New("取值方法缺失")
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
	return c.V, nil
}

// MGet
// 只有提前调用了RegisterHandle将方法注册进来后才可以调用该接口，否则返回数据会是空的
func MGet(key string, mode string, id string) (value interface{}, err error) {
	value, err = MGetWithFunc(key, id, mode, _handle[key])
	return
}
