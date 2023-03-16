package saCache

import (
	"time"
)

type Cache struct {
	Before     time.Time //有效期
	Data       interface{}
	UpdateAt   time.Time //更新时间1秒内也不会更新
	isUpdating bool      //是否正在更新中，避免并发导致多次异步更新处理
}
