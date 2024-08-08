package saGo

import (
	"fmt"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	fmt.Println("开始", time.Now().Format(time.DateTime))
	var pool = NewPool(100, 10, 20, 0, func(args interface{}) {
		time.Sleep(time.Millisecond * 10)
	})
	for i := 0; i < 100; i++ {
		pool.Invoke(i + 1)
	}
	pool.Wait()
	fmt.Println("完成", time.Now().Format(time.DateTime))
}
