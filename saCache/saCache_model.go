package saCache

import (
	"sync"
)

type cache struct {
	locker sync.RWMutex
	maxCnt int
	ary    []cacheItem
}

type cacheItem struct {
	id  string
	cnt int
	v   interface{}
}
