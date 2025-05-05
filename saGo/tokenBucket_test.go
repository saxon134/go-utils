package saGo

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	fmt.Println("开始", time.Now().Format(time.DateTime))
	var bucket = NewBucket(1, 0.5, func(bucket *Bucket, args interface{}) {
		time.Sleep(time.Millisecond * (10 + time.Duration(rand.Int64N(50))))
	})

	for i := 0; i < 5; i++ {
		fmt.Println(i, bucket.Desc())
		bucket.Invoke(i + 1)
	}
	bucket.Done()
	fmt.Println("完成", time.Now().Format(time.DateTime))
}
