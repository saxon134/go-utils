package saCache

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	for i := 0; i < 1000; i++ {
		r := i%20 + 1
		_, _ = MGetWithFunc("appInfo", "key-"+saData.Itos(r), "10s", func(id string) (interface{}, error) {
			return r, nil
		})
		_, _ = MGetWithFunc("appInfo", "key-0", "10s", func(id string) (interface{}, error) {
			return "230", nil
		})
		fmt.Println(saData.String(_cache))
		time.Sleep(time.Millisecond)
	}
}
