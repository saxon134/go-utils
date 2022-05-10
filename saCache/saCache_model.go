package saCache

type cache struct {
	MaxCnt   int
	TotalCnt int
	LastTime int64
	Ary      []cacheItem
}

type cacheItem struct {
	Cnt      int
	Id       string
	LastTime int64 //10位时间戳
	V        interface{}
}
