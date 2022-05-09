package saCache

import (
	"fmt"
	"github.com/saxon134/go-utils/saData"
	mathrand "math/rand"
	"testing"
)

func TestCache(t *testing.T) {
	for i := 0; i < 1000; i++ {
		r := mathrand.Intn(50)
		v, err := MGetWithFunc("appInfo", "key-"+saData.Itos(r), "10s", func(id string) (interface{}, error) {
			return i, nil
		})
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Println(r, v)
		fmt.Println(saData.String(_cache))
	}

}
