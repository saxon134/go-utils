package saGo

import (
	"fmt"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	fmt.Println("开始", time.Now().Format(time.DateTime))
	var pool = NewPool(10, 10, func(p *Pool, args interface{}) {
		LimiterLock("test", 1, 2)
		time.Sleep(time.Second * 5)
		fmt.Println("执行中：", args)
		LimiterUnLock("test")
	})
	for i := 0; i < 20; i++ {
		pool.Invoke(i + 1)
	}
	pool.Wait()
	fmt.Println("完成", time.Now().Format(time.DateTime))
}
