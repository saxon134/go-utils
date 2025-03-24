package saGo

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	fmt.Println("开始", time.Now().Format(time.DateTime))
	var bucket = NewBucket(150, 50, func(bucket *Bucket, args interface{}) {
		time.Sleep(time.Millisecond * (10 + time.Duration(rand.Int64N(50))))
		if time.Now().Second()%10 == 0 {
			time.Sleep(time.Second * 5)
		}
	})

	for i := 0; i < 6001; i++ {
		if i%3000 == 0 {
			fmt.Println(i, bucket.Desc())
		}
		bucket.Invoke(i + 1)
	}
	bucket.Done()
	fmt.Println("完成", time.Now().Format(time.DateTime))
}
