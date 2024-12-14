package saGo

import (
	"fmt"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	fmt.Println("开始", time.Now().Format(time.DateTime))
	var bucket = NewBucket(30, 2, func(bucket *Bucket, args interface{}) {
		time.Sleep(time.Millisecond * 1000)

		for i := 0; i < 10; i++ {
			fmt.Println(fmt.Sprintf("重试：%d %d %d", args, i, time.Now().Unix()))
			bucket.Consume()
		}
	})

	for i := 0; i < 20; i++ {
		bucket.Invoke(i + 1)
	}
	bucket.Done()
	fmt.Println("完成", time.Now().Format(time.DateTime))
}
