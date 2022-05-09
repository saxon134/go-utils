package saCache

type cache struct {
	maxCnt   int
	totalCnt int
	lastTime int64 //10位时间戳
	ary      []cacheItem
}

type cacheItem struct {
	cnt      int
	id       string
	lastTime int64 //10位时间戳
	v        interface{}
}
