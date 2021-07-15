package saGo

import (
	"fmt"
	"testing"
	"time"
)

func TestRoutine(t *testing.T) {
	r := NewRoutine(10, 3, func(params interface{}) {
		i, ok := params.(int)
		if ok {
			time.Sleep(time.Second * 1)
			fmt.Println(i)
		}
	})

	for i := 0; i < 100; i++ {
		r.Do(i)
	}
	time.Sleep(time.Second * 11)
}
