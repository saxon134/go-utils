package saGo

import (
	"fmt"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	fmt.Println("开始", time.Now().Format(time.DateTime))
	var bucket = NewBucket(30, 5, func(bucket *Bucket, args interface{}) {
		time.Sleep(time.Millisecond * 1000)
		//var i = args.(int)
		//fmt.Println(i)

		bucket.Consume()
		time.Sleep(time.Millisecond * 200)
	})

	for i := 0; i < 30; i++ {
		bucket.Invoke(i + 1)
	}
	bucket.Done()
	fmt.Println("完成", time.Now().Format(time.DateTime))
}
