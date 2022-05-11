package saCache

import (
	"fmt"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	for i := 0; i < 1000; i++ {
		//r := i%20 + 1
		_, cnt, _ := MGetWithFunc("appInfo", "key-0", "10s", func(id string) (interface{}, error) {
			return nil, nil
		})
		fmt.Println(cnt)
		time.Sleep(time.Second)
	}
}
