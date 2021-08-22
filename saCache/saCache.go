package saCache

import "github.com/saxon134/go-utils/saHit"

var _cache = make(map[string]*cache, 10)

//内存cache
func MGet(key string, id string, getHandle func() (interface{}, error)) interface{} {
	if key == "" || id == "" {
		return nil
	}

	item := _cache[key]
	if item == nil {
		item = new(cache)
	}

	var c *cacheItem = nil
	for _, v := range item.ary {
		if v.id == id {
			c = &v
			break
		}
	}

	if c != nil {
		item.locker.Lock()
		item.locker.Unlock()

		c.cnt++
		if c.cnt > item.maxCnt {
			item.maxCnt = c.cnt
		}
		return c.v
	} else if getHandle != nil {
		//做多保存10类数据
		if len(_cache) > 20 {
			return nil
		}

		v, err := getHandle()
		if err != nil {
			return nil
		}

		c = new(cacheItem)
		item.locker.Lock()
		defer item.locker.Unlock()

		c.cnt = saHit.Int(item.maxCnt > 1, item.maxCnt/2, 1)
		c.v = v
		c.id = id
		if len(item.ary) < 20 {
			item.ary = append(item.ary, *c)
		} else {
			//取次数最小的，替换掉
			var min *cacheItem
			var idx = 0
			for i, v := range item.ary {
				if min == nil || min.cnt == 0 || min.cnt > v.cnt {
					min = &v
					idx = i
				}
			}
			if min != nil && min.cnt >= 0 {
				item.ary = append(item.ary[:idx], item.ary[idx+1:]...)
			}
		}
	}
	return nil
}

//Redis cache
func RGet(key string, id string, getHandle func() (interface{}, error)) interface{} {
	//todo
	return nil
}
