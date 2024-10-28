package saGo

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	fmt.Println("开始", time.Now().Format(time.DateTime))
	var lock sync.Mutex
	var ary = make([]int, 0, 100)
	var pool = NewPool(9000, 100, 0, 1000, func(args interface{}) {
		time.Sleep(time.Millisecond * 5000)
		var i = args.(int)
		lock.Lock()
		ary = append(ary, i+1)
		lock.Unlock()
	})
	for i := 0; i < 9000; i++ {
		pool.Invoke(i + 1)
		if i%100 == 0 {
			fmt.Println(time.Now().Format(time.DateTime), "已完成：", i+1)
		}
	}
	pool.Wait()
	fmt.Println("完成", time.Now().Format(time.DateTime))
}

func TestPool2(t *testing.T) {
	fmt.Println("开始", time.Now().Format(time.DateTime))
	var pool = NewPool(0, 5, 10, 0, func(args interface{}) {
		time.Sleep(time.Millisecond * 5000)
		var i = args.(int)
		fmt.Println(i)
	})
	for i := 0; i < 20; i++ {
		pool.Invoke(i + 1)
		if i == 8 {
			pool.Done()
			break
		}
	}
	pool.Wait()
	fmt.Println("完成", time.Now().Format(time.DateTime))
}
