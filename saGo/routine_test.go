package saGo

import (
	"fmt"
	"testing"
	"time"
)

func TestRoutine(t *testing.T) {
	r := NewRoutine(3, time.Second*2, func(params interface{}) {
		i, ok := params.(int)
		if !ok {
			return
		}

		time.Sleep(time.Second * 1)
		fmt.Println(i)
	})

	for i := 0; i < 10; i++ {
		r.Do(i)
	}
	time.Sleep(time.Second * 10)
}
