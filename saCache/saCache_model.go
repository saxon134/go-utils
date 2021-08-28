package saCache

type cache struct {
	maxCnt int
	ary    []cacheItem
}

type cacheItem struct {
	cnt int
	id  string
	v   interface{}
}
