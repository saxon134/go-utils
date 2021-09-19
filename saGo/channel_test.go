package saGo

import (
	"fmt"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	for i := 0; i < 20; i++ {
		go func(v int) {
			Channel("channel_test1", 1, func(params ...interface{}) {
				fmt.Println(params...)
				time.Sleep(time.Second * 1)
			}).Do(v)
		}(i)
	}

	time.Sleep(time.Second * 22)
}
