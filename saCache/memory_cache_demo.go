/**
建议每个项目，单独搞一个缓存的包
统一处理缓存未命中时的处理，避免分散在各处代码
*/
package saCache

import (
	"sync"
	"time"
)

var levels struct {
	locker    sync.RWMutex
	updatedAt time.Time

	Ary []*Level
}

type Level struct {
	Name string
	V    int
}

func GetLevels() []*Level {
	if len(levels.Ary) > 0 &&
		levels.updatedAt.Add(time.Minute*10).After(time.Now()) {
		return levels.Ary
	}

	levels.locker.Lock()
	defer levels.locker.Unlock()

	//获取数据

	levels.updatedAt = time.Now()
	return levels.Ary
}

func GetLevel(name string) *Level {
	return nil
}
